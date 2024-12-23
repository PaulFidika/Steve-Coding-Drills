import torch
from diffusers import FluxPipeline
from PIL import Image
import numpy as np
import time
import argparse
from torchvision.transforms import ToPILImage
from datetime import datetime
from torchao.quantization import float8_dynamic_activation_float8_weight, quantize_, int8_dynamic_activation_int8_weight
from torchao.quantization.quant_api import PerRow
import random

print("Configuring torch inductor settings")
# CUDA setup
torch.set_float32_matmul_precision("high")

print("Configuring torch inductor settings")
torch._inductor.config.conv_1x1_as_mm = True
torch._inductor.config.coordinate_descent_tuning = True
torch._inductor.config.epilogue_fusion = False
torch._inductor.config.coordinate_descent_check_all_directions = True

# Model dictionary
MODEL_DICT = {
    "dev": "black-forest-labs/FLUX.1-dev",
    "schnell": "black-forest-labs/FLUX.1-schnell"
}

def setup_model(model_name, compile):
    if model_name not in MODEL_DICT:
        raise ValueError(f"Unknown model: {model_name}. Available models: {', '.join(MODEL_DICT.keys())}")
    
    auth_token = "hf_y"
    
    repo_id = MODEL_DICT[model_name]
    start_time = time.time()
    pipe = FluxPipeline.from_pretrained(
        repo_id,
        torch_dtype=torch.bfloat16,
        local_files_only=False,
        token=auth_token
    ).to("cuda")
    loading_time = time.time() - start_time
    print(f"Model {model_name} loaded successfully in {loading_time:.2f} seconds")
    
    pipe.set_progress_bar_config(disable=True)
    
    # Quantization step
    transformer_start_time = time.time()

    quantize_(pipe.transformer, int8_dynamic_activation_int8_weight())
    # quantize_(pipe.transformer, float8_dynamic_activation_float8_weight(granularity=PerRow()))
    transformer_quantization_time = time.time() - transformer_start_time
    print(f"Transformer quantization completed in {transformer_quantization_time:.2f} seconds")

    vae_start_time = time.time()
    quantize_(pipe.vae, int8_dynamic_activation_int8_weight())
    # quantize_(pipe.vae, float8_dynamic_activation_float8_weight(granularity=PerRow()))
    vae_quantization_time = time.time() - vae_start_time
    print(f"VAE quantization completed in {vae_quantization_time:.2f} seconds")

    # Torch compilation step
    if compile:
        transformer_start_time = time.time()
        pipe.transformer = torch.compile(pipe.transformer, mode="max-autotune", fullgraph=True)
        transformer_compilation_time = time.time() - transformer_start_time
        
        vae_start_time = time.time()
        pipe.vae = torch.compile(pipe.vae, mode="max-autotune", fullgraph=True)
        vae_compilation_time = time.time() - vae_start_time
        
        pipe.transformer.to(memory_format=torch.channels_last)
        pipe.vae.to(memory_format=torch.channels_last)
        
        print(f"Transformer compilation completed in {transformer_compilation_time:.2f} seconds")
        print(f"VAE compilation completed in {vae_compilation_time:.2f} seconds")
    else:
        print("Model not compiled")

    
    # Perform GPU -> CPU -> GPU round trips
    # gpu_to_cpu_times = []
    # cpu_to_gpu_times = []
    
    # for i in range(3):
    #     # GPU to CPU
    #     start_time = time.time()
    #     pipe.to("cpu")
    #     gpu_to_cpu_time = time.time() - start_time
    #     gpu_to_cpu_times.append(gpu_to_cpu_time)
        
    #     # CPU to GPU
    #     start_time = time.time()
    #     pipe.to("cuda")
    #     cpu_to_gpu_time = time.time() - start_time
    #     cpu_to_gpu_times.append(cpu_to_gpu_time)
    
    # # Calculate and print average times
    # avg_gpu_to_cpu = sum(gpu_to_cpu_times) / len(gpu_to_cpu_times)
    # avg_cpu_to_gpu = sum(cpu_to_gpu_times) / len(cpu_to_gpu_times)
    # print(f"\nAverage GPU -> CPU time: {avg_gpu_to_cpu:.4f}s")
    # print(f"Average CPU -> GPU time: {avg_cpu_to_gpu:.4f}s")
    
    return pipe

def tensor_to_pil(tensor: torch.Tensor) -> list[Image.Image]:
    tensor = tensor.to(dtype=torch.float16)
    transform = ToPILImage()
    return [transform(t) for t in tensor]

def run_benchmark(pipe, batch_size, num_iterations, model_name):
    prompt = "A beautiful woman short hair"
    print(f"Starting benchmark for {model_name} batch size {batch_size}")
    total_time = 0
    results = []
    
    steps = 28 if model_name == "dev" else 4
    
    image_sizes = [(1024, 1024), (1280, 768), (768, 1280)]
    for i in range(num_iterations):
        start_time = time.time()
        
        # Rotate through the three image sizes
        height, width = image_sizes[i % len(image_sizes)]
        
        images = pipe(
            prompt=prompt,
            guidance_scale=3.5,
            height=height,
            width=width,
            num_inference_steps=steps,
            num_images_per_prompt=batch_size
        ).images
        
        # Convert to PIL and save images
        for j, img in enumerate(images):
            img.save(f"image_{i}_{j}.png")
        
        inference_time = time.time() - start_time
        total_time += inference_time
        
        print(f"Batch size {batch_size}, Iteration {i+1}/{num_iterations}: Total time (inference + saving) {inference_time:.4f}s")
        results.append(f"Iteration {i+1}: {inference_time:.4f}s")
        
        # Test moving the pipe around
        # start_time = time.time()
        # pipe.to("cpu")
        # gpu_to_cpu_time = time.time() - start_time
        # print(f"GPU to CPU time: {gpu_to_cpu_time:.4f}s")

        # start_time = time.time()
        # pipe.to("cuda")
        # cpu_to_gpu_time = time.time() - start_time
        # print(f"CPU to GPU time: {cpu_to_gpu_time:.4f}s")

        # Add these times to the results
        # results.append(f"GPU to CPU time: {gpu_to_cpu_time:.4f}s")
        # results.append(f"CPU to GPU time: {cpu_to_gpu_time:.4f}s")

    avg_time = total_time / num_iterations
    images_per_second = (batch_size * num_iterations) / total_time
    
    print(f"\nBenchmark results for batch size {batch_size}:")
    print(f"  Total iterations: {num_iterations}")
    print(f"  Average time (inference + saving): {avg_time:.4f}s")
    print(f"  Images per second: {images_per_second:.2f}")

    results.append(f"Average time: {avg_time:.4f}s")
    results.append(f"Images per second: {images_per_second:.2f}")

    return avg_time, images_per_second, results

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Run FLUX model inference benchmark")
    parser.add_argument("--model", type=str, choices=MODEL_DICT.keys(), help="Model to use for inference")
    parser.add_argument("--compile", action='store_true', help="Compile the model")
    parser.add_argument("--batch_size", type=int, default=1, help="Batch size for inference")
    parser.add_argument("--num_iterations", type=int, default=9, help="Number of iterations to run")
    args = parser.parse_args()

    print(f"Starting benchmark with model: {args.model}, batch size: {args.batch_size}, iterations: {args.num_iterations}")

    pipe = setup_model(args.model, args.compile)
    avg_time, imgs_per_sec, results = run_benchmark(pipe, args.batch_size, args.num_iterations, args.model)

    print(f"\nFinal Benchmark Results for model {args.model}, batch size {args.batch_size}:")
    print(f"  Average time (inference + saving): {avg_time:.4f}s")
    print(f"  Images per second: {imgs_per_sec:.2f}")

    # Save results to a file
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    result_filename = f'flux_benchmark_results_{timestamp}.txt'
    with open(result_filename, 'w') as f:
        f.write(f"Benchmark Results for model {args.model}, batch size {args.batch_size}, iterations {args.num_iterations}\n")
        f.write("\n".join(results))
    
    print(f"Detailed results saved to: {result_filename}")
    print("Benchmark completed.")
