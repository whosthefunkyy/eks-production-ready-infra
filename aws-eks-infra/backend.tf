terraform {
 
  required_version = ">= 1.10.0" 

  backend "s3" {
    bucket       = "terraform-whosthefunky-bucket-2026"
    key          = "terraform/infra/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    
    
    use_lockfile = true 
  }
}