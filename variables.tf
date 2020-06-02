variable "project_name" {
  description = "Project Name"
  default     = "cloudsim"
}

variable "iac_tool" {
  description = "Infrastructure as Code"
  default     = "Terraform"
}

locals {
  default_tag = {
    project_name = var.project_name
    iac_tool     = var.iac_tool
  }
}