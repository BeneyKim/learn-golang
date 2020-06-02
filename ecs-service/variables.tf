variable "project_name" {
  description = "Project Name"
}

variable "service_name" {
  description = "service name"
}

variable "iac_tool" {
  description = "Infrastructure as Code"
  default     = "Terraform"
}



########################################################################################################################
### ECS Task Definitions

variable "cluster_name" {
  description = "Cluster Name"
}


variable "container_definition_file" {
  description = "Container Definition File"
}

variable "container_name" {
  description = "Container Name"
}

variable "volumes" {
  type = list(string)
}

variable "az_count" {
  description = "Number of AZs to cover in a given AWS region"
  default     = "2"
}

variable "app_count" {
  description = "Number of docker containers to run"
  default     = 1
}



########################################################################################################################
### Compute

variable "fargate_cpu" {
  description = "Fargate instance CPU units to provision (1 vCPU = 1024 CPU units)"
  default     = "512"
}

variable "fargate_memory" {
  description = "Fargate instance memory to provision (in MiB)"
  default     = "256"
}



########################################################################################################################
### Network & Security

variable vpc_id {
  description = "VPC ID"
}

variable "security_group_id" {
  description = "Security Group ID"
}

variable "subnet_id" {
  description = "Security Group ID"
}



########################################################################################################################
### Load Balancer
//
//variable nlb_arn {
//  description = "Network Load Balancer ARN"
//  default = null
//}
//
//variable "lb_tcp_ports" {
//  description = "Network Load Balancer Target Group for TCP"
//  type = list(number)
//  default = []
//}
//
//variable "lb_tls_ports" {
//  description = "Network Load Balancer Target Group for TLS"
//  type = list(number)
//  default = []
//}
//
//variable zone_name {
//  description = "Domain Zone Name"
//  default = null
//}
//
//variable cert_domain {
//  description = "Domain Name for Certificate"
//  default = null
//}


########################################################################################################################
### Local Variables
locals {
  volumes = { for v in var.volumes : tostring(v) => v }

  // lb_exist = var.nlb_arn != null ? 1 : 0
  // tcp_ports = var.nlb_arn != null ? { for v in var.lb_tcp_ports : tostring(v) => v } : {}
  // tls_ports = var.nlb_arn != null ? { for v in var.lb_tls_ports : tostring(v) => v } : {}
  # zone_name = var.nlb_arn != null ? zone_name : null
  # cert_domain = var.nlb_arn != null ? cert_domain : null

  default_tag = {
    project_name = var.project_name
    iac_tool     = var.iac_tool
  }
}