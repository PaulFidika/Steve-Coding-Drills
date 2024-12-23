from diffusers import FluxPipeline
import torch
import time
import random
import os

auth_token = "hf_y"

# List all available CUDA devices
cuda_devices = torch.cuda.device_count()
print(f"Number of available CUDA devices: {cuda_devices}")

for i in range(cuda_devices):
    device = torch.cuda.get_device_properties(i)
    print(f"CUDA Device {i}:")
    print(f"  Name: {device.name}")
    print(f"  Compute Capability: {device.major}.{device.minor}")
    print(f"  Total Memory: {device.total_memory / 1024**3:.2f} GB")

# Empty PyTorch cache
torch.cuda.empty_cache()

pipeline = FluxPipeline.from_pretrained(
    "black-forest-labs/FLUX.1-schnell",
    torch_dtype=torch.bfloat16,
    token=auth_token
).to("cuda")

# black-forest-labs/FLUX.1-dev

# pipeline.unet.to(memory_format=torch.channels_last)
# pipeline.unet = torch.compile(pipeline.unet, mode="reduce-overhead", fullgraph=True)

# , num_images_per_prompt=4
for i in range(4):
    # torch.cuda.empty_cache()
    start_time = time.time()
    prompts = [
        "Astronaut in a jungle, cold color palette, muted colors, detailed, 8k",
        "Futuristic cityscape at night, neon lights, cyberpunk style, highly detailed",
        "Underwater scene with bioluminescent creatures, deep sea, ethereal atmosphere",
        "Ancient ruins overgrown with lush vegetation, mystical ambiance, golden hour lighting",
        "Steampunk airship flying through a storm, dramatic clouds, intricate machinery",
        "Alien landscape with bizarre flora, twin moons in the sky, otherworldly colors",
        "Medieval fantasy battle scene, dragons and knights, epic scale, hyper-realistic",
        "Surreal dreamscape with floating islands, impossible architecture, vivid imagination",
        "Post-apocalyptic wasteland, abandoned cities reclaimed by nature, atmospheric haze",
        "Microscopic view of exotic particles, abstract patterns, scientific visualization"
    ]
    prompt = random.choice(prompts)
    
    image = pipeline(
        prompt=prompt,
        num_images_per_prompt=1,
        num_inference_steps=4,
        width=1024,
        height=1024
    ).images[0]

    random_number = random.randint(10000, 99999)
    
    # Save the image in the 'images' folder
    os.makedirs('images', exist_ok=True)
    image.save(f"images/flux_{random_number}.png")
    
    print(f"Time: {time.time() - start_time}")

