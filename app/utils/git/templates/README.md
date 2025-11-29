# IAC Configuration Repository

This repository contains infrastructure as code configuration managed by the iac tool.

## Directory Structure

- `stacks/` - Docker Compose stack definitions
- `env/` - Environment variables and secrets
- `backups/` - Configuration backups

## Getting Started

1. Configure your stacks in the `stacks/` directory
2. Set up environment variables using `iac env add`
3. Deploy stacks with `iac stack up`

For more information, run `iac --help`
