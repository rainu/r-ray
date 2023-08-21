variable "repository_name" {
  description = "The name of the ecr repository"
}

resource "aws_ecr_repository" "_" {
  name = var.repository_name
}
