terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~>5.0"
    }
  }
  backend "s3" {
    bucket = "terraform-etaigh6ughahneiy"
    key    = "r-ray/state"
    region = "eu-central-1"
  }
}

provider "aws" {
  region = "eu-central-1"
  default_tags {
    tags = {
      "project": "r-ray"
    }
  }
}