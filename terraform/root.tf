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