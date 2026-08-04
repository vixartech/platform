# AGENTS.md

This repository contains shared local platform tooling for Vixartech projects.

## Scope

- Keep Docker Compose deployments, Dockerfiles, linters, and Make targets here.
- Do not add application business logic here.
- Do not add project-specific Go packages here. Reusable Go code belongs in `github.com/vixartech/gear`.
- Keep commands stable because service repositories include this repository as the `platform` submodule.

## Conventions

- `platform.mk` is the public Make contract consumed by projects.
- `deployments/` contains infrastructure services only.
- `tools/make/` contains reusable Make targets.
- `linters/` contains shared linter config.
- `bin/` is local generated state and must stay ignored.
