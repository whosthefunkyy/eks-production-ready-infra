# Production-like AWS EKS Platform

## Overview

This project demonstrates a production-like DevOps platform deployed on AWS using Terraform, Kubernetes and GitHub Actions.

The goal of the project was to build a complete infrastructure and deployment pipeline from scratch following modern DevOps practices, including Infrastructure as Code, secure authentication, automated CI/CD, monitoring, and private database connectivity.

The platform provisions AWS infrastructure, deploys a Go application to Amazon EKS using Helm, exposes it through an AWS Application Load Balancer, connects it to a private PostgreSQL database, and provides monitoring with Prometheus and Grafana.

---

## Project Structure

```text
bootstrap/
├── GitHub OIDC
├── IAM Roles
└── Remote Terraform Backend

infrastructure/
├── VPC
├── EKS
├── RDS
├── ECR
├── AWS Load Balancer Controller
├── kube-prometheus-stack
└── Terraform configuration

application/
├── Go API
├── Dockerfile
├── Helm Chart
└── Kubernetes manifests
```

---

## Infrastructure

Terraform provisions the following AWS resources:

* VPC
* Public and Private Subnets
* Internet Gateway
* NAT Gateway
* Amazon EKS Cluster
* Managed Node Groups
* Amazon ECR
* Amazon RDS PostgreSQL
* IAM Roles
* OIDC Provider
* Security Groups

The PostgreSQL database is deployed into private subnets and is accessible only from Kubernetes worker nodes through a dedicated Security Group.

---

## Bootstrap

Infrastructure authentication is configured separately using the bootstrap configuration.

Bootstrap creates:

* GitHub OIDC Provider
* IAM Role for GitHub Actions
* IAM Policies
* S3 Remote State
* DynamoDB State Locking

This allows GitHub Actions to authenticate to AWS without storing long-lived AWS Access Keys.

---

## CI/CD

The project contains four independent GitHub Actions workflows.

### Infrastructure

* Create infrastructure
* Destroy infrastructure

Terraform provisions and removes AWS resources automatically.

### Application

* Build and Deploy
* Destroy Application

The deployment workflow:

1. Builds the Go application Docker image
2. Pushes the image to Amazon ECR
3. Updates the Helm release
4. Performs a rolling update inside Kubernetes

Manual approval gates are used before infrastructure changes.

---

## Kubernetes

The application is deployed using Helm.

The platform includes:

* Deployment
* Service
* Ingress
* ConfigMap
* Secret
* ServiceAccount
* Horizontal Pod Autoscaler
* ServiceMonitor

Deployment configuration includes:

* Rolling Updates
* Readiness Probes
* Liveness Probes
* CPU and Memory Requests/Limits

---

## Networking

Traffic flow:

Internet

↓

AWS Application Load Balancer

↓

Ingress

↓

Kubernetes Service

↓

Go API Pods

↓

Amazon RDS PostgreSQL

AWS Load Balancer Controller automatically creates and manages the ALB based on Kubernetes Ingress resources.

Multiple applications can share a single ALB using IngressGroup.

---

## Security

The project follows modern AWS security practices.

Implemented:

* GitHub OIDC authentication
* IRSA (IAM Roles for Service Accounts)
* No static AWS credentials
* Private RDS deployment
* Least-Privilege IAM permissions
* Dedicated Security Groups

---

## Monitoring

Monitoring is implemented using kube-prometheus-stack.

The platform includes:

* Prometheus
* Grafana
* ServiceMonitor

The Go application exposes custom Prometheus metrics including:

* http_requests_total
* http_request_duration_seconds

Grafana is published through the existing AWS ALB using a dedicated `/grafana` path.

---

## Troubleshooting Experience

During development the following production-like issues were identified and resolved:

* Terraform dependency ordering
* GitHub OIDC authentication
* IRSA configuration
* Helm deployment conflicts
* AWS Load Balancer Controller permissions
* ALB Ingress routing priority
* Shared ALB configuration using IngressGroup
* Kubernetes rollout behavior
* Prometheus ServiceMonitor configuration

---
> **Note:** Replace `262778473495` with your own AWS account ID in 
> `bootstrap/iam/github-oidc-trust-policy.json` before running bootstrap.

## Quick Start
1. Run bootstrap once
2. Trigger Infrastructure UP workflow
3. Trigger Application Deploy workflow
4. Access API: http://<alb-dns>/health
5. Access Grafana: http://<alb-dns>/grafana

## Technologies

* Terraform
* AWS (EKS, VPC, IAM, ALB, RDS, ECR, S3)
* Docker
* Kubernetes
* Helm
* GitHub Actions
* GitHub OIDC
* Go
* PostgreSQL
* Prometheus
* Grafana
* Linux
