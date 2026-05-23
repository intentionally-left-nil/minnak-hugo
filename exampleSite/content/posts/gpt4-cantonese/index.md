---
title: "What GPT-4 Gets Wrong About Cantonese"
date: 2026-03-01T14:00:00Z
category: ["AI"]
tag: ["llm", "language", "evaluation"]
summary: "Large language models perform surprisingly well on Mandarin but struggle with Cantonese in ways that reveal the limits of training data distribution."
cover:
  src: "feature.jpg"
  alt: "Photo by Kevin Canlas on Unsplash"
---

Large language models perform surprisingly well on Mandarin but struggle with Cantonese in ways that reveal the limits of training data distribution.

Cantonese is spoken by roughly 85 million people, primarily in Guangdong province and Hong Kong. Despite this, it is dramatically underrepresented in LLM training corpora compared to Mandarin. The consequences are predictable: models that can write fluent Mandarin prose will produce Cantonese text that reads as unnatural to native speakers.

## The written/spoken divide

Cantonese has a written form — used in informal digital contexts like WhatsApp and social media — that differs significantly from standard written Chinese. Characters like 係 (hai6, "to be"), 唔 (m4, "not"), and 嘅 (ge3, possessive particle) are distinctly Cantonese.

GPT-4 tends to avoid these characters even when explicitly prompted to write in Cantonese, defaulting instead to Mandarin grammar patterns written in traditional characters. The output looks like Cantonese to a non-speaker but reads as unnatural to a native.

## Why this happens

Training data for Chinese-language models skews heavily toward simplified Mandarin text scraped from mainland Chinese sources. Traditional character Cantonese text from Hong Kong is a small fraction of that. Models learn the statistical patterns of the dominant dialect and apply them to all requests for "Chinese."

This is a data distribution problem, not a capability problem. Models fine-tuned on Cantonese-specific corpora perform noticeably better.
