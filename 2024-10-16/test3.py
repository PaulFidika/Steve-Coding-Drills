import torch
from diffusers import FluxPipeline
from PIL import Image
import numpy as np
import time
import argparse
from torchvision.transforms import ToPILImage
from datetime import datetime
from torchao.quantization import int8_dynamic_activation_int8_weight, quantize_

# CUDA setup
torch.set_float32_matmul_precision("high")

# Model dictionary
MODEL_DICT = {
    "dev": "black-forest-labs/FLUX.1-dev",
    "schnell": "black-forest-labs/FLUX.1-schnell"
}

def setup_model(model_name):
    if model_name not in MODEL_DICT:
        raise ValueError(f"Unknown model: {model_name}. Available models: {', '.join(MODEL_DICT.keys())}")
    
    repo_id = MODEL_DICT[model_name]
    pipe = FluxPipeline.from_pretrained(repo_id, torch_dtype=torch.bfloat16, local_files_only=True).to("cuda")
    print(f"Model {model_name} loaded successfully from {repo_id}")
    
    print("Configuring torch inductor settings")
    torch._inductor.config.conv_1x1_as_mm = True
    torch._inductor.config.coordinate_descent_tuning = True
    torch._inductor.config.epilogue_fusion = False
    torch._inductor.config.coordinate_descent_check_all_directions = True
    
    pipe.set_progress_bar_config(disable=True)
    
    # Quantization step
    print("Starting quantization of transformer model to int8_dynamic_activation_int8_weight")
    quantize_(pipe.transformer, int8_dynamic_activation_int8_weight())
    print("Model transformer quantized to int8_dynamic_activation_int8_weight successfully")

    print("Starting quantization of vae model to int8_dynamic_activation_int8_weight")
    quantize_(pipe.vae, int8_dynamic_activation_int8_weight())
    print("Model transformer vae to int8_dynamic_activation_int8_weight successfully")
    
    print("Converting transformer to channels_last memory format")
    pipe.transformer.to(memory_format=torch.channels_last)

    print("Converting vae to channels_last memory format")
    pipe.vae.to(memory_format=torch.channels_last)
    
    print("Compiling transformer with mode='max-autotune' and fullgraph=True")
    pipe.transformer = torch.compile(pipe.transformer, mode="max-autotune", fullgraph=True)
    print("Transformer compilation completed")

    print("Compiling vae with mode='max-autotune' and fullgraph=True")
    pipe.vae = torch.compile(pipe.vae, mode="max-autotune", fullgraph=True)
    print("Vae compilation completed")
    
    return pipe

def tensor_to_pil(tensor: torch.Tensor) -> list[Image.Image]:
    tensor = tensor.to(dtype=torch.float16)
    transform = ToPILImage()
    return [transform(t) for t in tensor]

def run_benchmark(pipe, batch_size, num_iterations):
    prompt = "A dinosaur in modern suit"
    print(f"Starting benchmark for batch size {batch_size}")
    total_time = 0
    results = []
    
    for i in range(num_iterations):
        start_time = time.time()
        
        images = pipe(
            prompt=prompt,
            guidance_scale=3.5,
            height=1024,
            width=1024,
            num_inference_steps=28,
            num_images_per_prompt=batch_size
        ).images
        
        # Convert to PIL and save images
        for j, img in enumerate(images):
            img.save(f"image_{i}_{j}.png")
        
        inference_time = time.time() - start_time
        total_time += inference_time
        
        print(f"Batch size {batch_size}, Iteration {i+1}/{num_iterations}: Total time (inference + saving) {inference_time:.4f}s")
        results.append(f"Iteration {i+1}: {inference_time:.4f}s")

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
    parser.add_argument("model", type=str, choices=MODEL_DICT.keys(), help="Model to use for inference")
    parser.add_argument("batch_size", type=int, help="Batch size for inference")
    parser.add_argument("num_iterations", type=int, help="Number of iterations to run")
    args = parser.parse_args()

    print(f"Starting benchmark with model: {args.model}, batch size: {args.batch_size}, iterations: {args.num_iterations}")

    pipe = setup_model(args.model)
    avg_time, imgs_per_sec, results = run_benchmark(pipe, args.batch_size, args.num_iterations)

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