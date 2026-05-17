---
title: "Setting Up a Home Lab on a Decommissioned ThinkPad"
date: 2026-02-20T09:00:00Z
category: ["Technology"]
tag: ["homelab", "linux", "networking"]
summary: "An old ThinkPad X230, a fresh NixOS install, and a few evenings of configuration later — here's what I learned running a personal server at home."
feature_image_alt: "Photo by André Furtado on Unsplash"
---

An old ThinkPad X230, a fresh NixOS install, and a few evenings of configuration later — here's what I learned running a personal server at home.

## Why bother?

Cloud services are convenient until they're not. S3 costs money. GitHub rate-limits the API. Your CI runner spins up in a distant region. A home lab gives you a local alternative for development infrastructure — fast, cheap, under your control.

## Hardware

The X230 cost $40 on eBay. It has a dual-core i5, 8 GB RAM, and a 256 GB SSD I had lying around. That's enough to run:

- A private container registry
- A local DNS resolver (Pi-hole)
- A Gitea instance for private repos
- A Wireguard VPN endpoint

## The NixOS choice

I chose NixOS because I wanted the configuration to be reproducible and version-controlled. The entire system state lives in `configuration.nix`. Roll back a bad change with one command. Rebuild from scratch on any machine with one command.

The learning curve is steep — the Nix language is genuinely odd — but the payoff is a system you can reason about completely.

## What I'd do differently

Start with Ubuntu if this is your first home lab. NixOS is worth it eventually, but fighting the package model while also learning networking fundamentals is too much context-switching.
