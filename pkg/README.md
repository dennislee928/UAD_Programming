# pkg/

This directory contains **public, reusable packages** that can be imported by external projects.

## Purpose

- Public API for UAD language tools
- Libraries that other Go projects can depend on
- Stable, documented interfaces

## Guidelines

- All code in `pkg/` should be well-documented
- Maintain backward compatibility
- Follow semantic versioning principles
- Do not include internal implementation details

## Planned Packages

- `pkg/uad/`: Core UAD language client library
- `pkg/ast/`: Public AST manipulation utilities
- `pkg/runtime/`: Runtime configuration and interfaces

## Status

🚧 This directory is currently empty. Packages will be added as the API stabilizes.


