# Production-like EKS Infrastructure

This repository is a DevOps practice project that provisions an AWS EKS environment with Terraform and deploys a small Go API with Docker, ECR, and Helm through GitHub Actions OIDC.

## What is included

- One-time AWS bootstrap for GitHub OIDC and the GitHub Actions IAM role.
- Terraform-managed EKS, VPC, node group, ECR repository, and AWS Load Balancer Controller.
- Manual GitHub Actions workflows to create and destroy the infrastructure.
- Automated application workflow that builds the Go API image, pushes it to ECR, and deploys the Helm chart to EKS.

## Repository layout

```text
aws-bootstrap-infra/  # One-time IAM/OIDC bootstrap. Do not run during destroy.
aws-eks-infra/        # Main Terraform infrastructure managed by up/down workflows.
go-api-on-eks/        # Go application, Dockerfile, and Helm chart.
.github/workflows/    # Infrastructure and application automation.
```

## Workflow order

1. Run `aws-bootstrap-infra` once from a trusted local/admin environment.
2. Run `Infrastructure UP` from GitHub Actions to create EKS and supporting resources.
3. Run `Application Deploy` to build, push, and deploy the Go API.
4. Run `Infrastructure DOWN` when the environment is no longer needed. Type `destroy` in the confirmation input.

The destroy workflow only targets `aws-eks-infra`, so the bootstrap IAM role and GitHub OIDC provider stay in place for future runs.

## Required GitHub/AWS assumptions

- AWS account: `262778473495`
- Region: `us-east-1`
- GitHub Actions role: `arn:aws:iam::262778473495:role/GitHubActionsEKSRole`
- EKS cluster name: `terraform-eks-v3`
- ECR repository: `go-api-image`

If the repository owner/name changes, update the GitHub OIDC trust policy in `aws-bootstrap-infra/iam/github-oidc-trust-policy.json` and re-apply the bootstrap stack once.
