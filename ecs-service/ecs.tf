data aws_caller_identity "current" {}
data aws_region "current" {}
data aws_vpc "main" { id = var.vpc_id }
data aws_security_group "main" { id = var.security_group_id }
data aws_subnet "main" { id = var.subnet_id }


resource "aws_ecs_task_definition" "task" {
  family       = "${var.service_name}-task"
  network_mode = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory

  container_definitions = file(var.container_definition_file)

  dynamic "volume" {
    for_each = local.volumes
    content {
      name = volume.key
    }
  }

  task_role_arn      = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/ecsTaskExecutionRole"
  execution_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/ecsTaskExecutionRole"

  tags = local.default_tag
}

resource "aws_ecs_service" "service" {
  name            = "${var.service_name}-service"
  cluster         = data.aws_ecs_cluster.main.arn
  task_definition = aws_ecs_task_definition.task.arn
  desired_count   = var.app_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups = [data.aws_security_group.main.id]
    subnets = [
      data.aws_subnet.main.id
    ]
    assign_public_ip = true
  }

//  dynamic "load_balancer" {
//    for_each = local.tcp_ports
//
//    content {
//      target_group_arn = aws_lb_target_group.tcp[load_balancer.key].arn
//      container_name = var.container_name
//      container_port = load_balancer.key
//    }
//  }
//
//  dynamic "load_balancer" {
//    for_each = local.tls_ports
//
//    content {
//      target_group_arn = aws_lb_target_group.tls[load_balancer.key].arn
//      container_name = var.container_name
//      container_port = load_balancer.key
//    }
//  }

  provisioner "local-exec" {
    command = "ecs-cli configure --cluster ${data.aws_ecs_cluster.main.cluster_name} --default-launch-type ${self.launch_type} --region ${data.aws_region.current.name} --config-name ${data.aws_ecs_cluster.main.cluster_name}"
  }
}

data aws_ecs_cluster "main" {
  cluster_name = var.cluster_name
}

//resource "aws_ecs_cluster" "pbx" {
//  name               = "${var.service_name}-cluster"
//  capacity_providers = ["FARGATE"]
//
//  setting {
//    name  = "containerInsights"
//    value = "enabled"
//  }
//
//  tags = local.default_tag
//}
