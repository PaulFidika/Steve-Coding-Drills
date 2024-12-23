from diffusers.pipelines.stable_diffusion_xl.pipeline_stable_diffusion_xl import StableDiffusionXLPipeline
import torch
import time

# Load the SDXL pipeline with FP16 precision for reduced memory usage
pipeline = StableDiffusionXLPipeline.from_pretrained(
    'stabilityai/stable-diffusion-xl-base-1.0', 
    torch_dtype=torch.float16
)

# Move the pipeline to GPU
start_time = time.time()
pipeline = pipeline.to('cuda')
print(f"Moving pipeline to GPU took {time.time() - start_time:.2f} seconds")

# Enable attention slicing to optimize memory (optional)
pipeline.enable_attention_slicing()
# pipeline.enable_xformers_memory_efficient_attention()
pipeline.enable_vae_tiling()
# pipe.enable_model_cpu_offload()

# Compile the UNet model
start_time = time.time()
pipeline.unet = torch.compile(pipeline.unet)
print(f"Compiling UNet model took {time.time() - start_time:.2f} seconds")

# Compile the text encoders
start_time = time.time()
pipeline.text_encoder = torch.compile(pipeline.text_encoder)
pipeline.text_encoder_2 = torch.compile(pipeline.text_encoder_2)
print(f"Compiling text encoders took {time.time() - start_time:.2f} seconds")

# Compile the VAE decoder
start_time = time.time()
pipeline.vae.decoder = torch.compile(pipeline.vae.decoder)
print(f"Compiling VAE decoder took {time.time() - start_time:.2f} seconds")

# Run inference with the compiled pipeline
start_time = time.time()
prompt = "An astronaut riding a horse on the moon."
images = pipeline(prompt=prompt, num_inference_steps=25, height=1024, width=1024).images
print(f"Running inference took {time.time() - start_time:.2f} seconds")

# Save the output image
start_time = time.time()
images[0].save("output.png")
print(f"Saving output image took {time.time() - start_time:.2f} seconds")
