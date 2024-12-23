Results on A100 80GB:

Flux dev, 28 steps

- Disc -> GPU: 8 - 12 seconds
- GPU -> CPU: 10 seconds
- CPU -> GPU: 4 seconds

- Quantization time is negligable
- First run compilation time: ~90 seconds

- Inference without compilation: 270 seconds (36 GBs VRAM)
- Inference with compilation: 8.9 seconds (25 GBs VRAM)

Flux Schnell, 4 steps

- First run without compilation: 41 seconds
- First run with compilation: ~413 seconds

- Inference without compilation: 43 seconds (36 GBs VRAM)
- Inference without compilation (BUT xformers!): 35 seconds (36 GBs VRAM)
- Inference (corrected) without compilation: 2.6 seconds (36 GBs VRAM)
- Inference with compilation: 1.8 seconds

Notes:
- You cannot move compiled models to CPU
- Compiled models can be trained for multiple image-sizes
- For some reason, the non-compiled models have god-awful performance. I don't know what's happening.


### MiG Test:
Run on H100, Schnell with 4 steps

- 2x MiG Mode: 2.9 - 4 secs (depends on power level of the instance)
- Without MiG: 1.85 - 2.4 seconds

That is, with MiG we get 1 image every 1.7 seconds
Without MiG, we get 1 image every 2 seconds

We get about a 17% increase in throughput, which isn't super significant. This is probably because Flux needs to use all of the compute-units available to it. Perhaps weaker models require fewer compute units?

Run on H100, SDXL:
- Without MiG: 2.5 seconds (24 images per min)
- 2x MiG: 3.15 seconds (38 images per min) (3.8 and 4.5 = 29 images per min)
- 4x MiG: 6.4 seconds (medium) and 10.2 seconds (small) (34 images per min)
- 7x MiG: 21 seconds (small) (20 images per min)

The 10 GB memory requirement really makes life hard for the SXL model and slows it down; if I was handling memory better it might work better.
