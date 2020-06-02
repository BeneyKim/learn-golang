locals {
  cloudsim = "cloudsim"
}

module "cloudsim" {
  source = "./ecs-service"

  project_name = var.project_name
  service_name = "nt-${local.cloudsim}"

  cluster_name              = aws_ecs_cluster.cloudsim.name
  container_definition_file = "${path.root}/container-definitions/${local.cloudsim}.json"
  container_name            = local.cloudsim
  volumes = []
  az_count       = 1
  app_count      = 1
  fargate_cpu    = 256
  fargate_memory = 512

  vpc_id            = aws_vpc.main.id
  security_group_id = aws_security_group.main.id
  subnet_id         = aws_subnet.public[0].id

  //  nlb_arn = aws_lb.lb.arn
  //  lb_tcp_ports = [
  //    80
  //  ]
  //  lb_tls_ports = [
  //    443
  //  ]
  //
  //  zone_name   = data.aws_route53_zone.ntel.name
  //  cert_domain = aws_route53_record.pbx.name
}


resource "aws_ecs_cluster" "cloudsim" {
  name               = "${var.project_name}-cluster"
  capacity_providers = ["FARGATE"]

  setting {
    name  = "containerInsights"
    value = "disabled"
  }

  tags = local.default_tag
}