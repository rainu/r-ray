module "api-gateway" {
  source = "./api-gateway"
}

variable "LAMBDA_DOCKER_IMAGE" {
  description = "The complete path to the aws lambda docker image"
}

module "lambda" {
  source = "./lambda"

  docker_image = var.LAMBDA_DOCKER_IMAGE
}

variable "ECR_REPOSITORY_NAME" {
  description = "The name of the ecr repository"
}

module "ecr" {
  source = "./ecr"
  repository_name = var.ECR_REPOSITORY_NAME
}