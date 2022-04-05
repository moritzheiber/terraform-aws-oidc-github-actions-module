terraform {
  required_version = ">= 1"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 3.2.0"
    }
  }
}
