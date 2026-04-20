# AP2 Contracts Proto Repository

This repository contains only protobuf contracts used by Order and Payment services.

## Repository Link

- https://github.com/cureeeeee/advancedprogramming2/tree/main/contracts-proto

## Generated Repository Link

- https://github.com/cureeeeee/advancedprogramming2/tree/main/contracts-generated

## Local Generation
Use Buf:

```bash
buf generate
```

Generated code should be pushed to a separate repository (`ap2-contracts-generated`) by GitHub Actions.

## CI Remote Generation

Workflow file:

- `.github/workflows/remote-generate.yml`

Required secret:

- `GENERATED_REPO_PAT` (token with push access to generated repository)
