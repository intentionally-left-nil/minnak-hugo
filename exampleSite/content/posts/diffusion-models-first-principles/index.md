---
title: "Diffusion Models From First Principles"
date: 2026-03-10T12:00:00Z
categories: ["AI"]
tags: ["diffusion", "ml", "generative-models"]
summary: "Denoising diffusion probabilistic models are often explained through the lens of score matching or variational inference. Here's a more intuitive path."
feature_image_alt: "Whimsical drawings on brown paper attached to a wall — Photo by Brother Yoon on Unsplash"
---

Denoising diffusion probabilistic models are often explained through the lens of score matching or variational inference. Here's a more intuitive path.

## The core idea

Start with a data distribution — say, the distribution of natural images. Repeatedly add a small amount of Gaussian noise until the distribution is indistinguishable from pure noise. This is the *forward process*.

Now train a neural network to reverse this process: given a slightly noisy image, predict what it looked like with a tiny bit less noise. Chain these denoising steps together and you have a generative model.

## Why this works

The key insight is that denoising at any given noise level is a well-posed problem. The network doesn't need to hallucinate a complete image from nothing — it just needs to remove a small, predictable amount of noise from an already-nearly-correct image.

This decomposition of generation into many small, manageable steps is what makes diffusion models both stable to train and high-quality in output.

## The connection to score matching

The denoising network is implicitly learning the *score* — the gradient of the log probability density — of the data distribution at each noise level. This is why diffusion models and score-based generative models converge to the same thing.

Understanding this connection explains why classifier-free guidance works: you're interpolating between the conditional and unconditional score functions, amplifying the signal that distinguishes a specific class from the overall distribution.
