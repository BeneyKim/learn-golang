data "external" "read_ecs_public_ip" {

  program = ["python", "${path.module}/get-ecs-public-ip.py"]

  query = {
    cluster_config_name = data.aws_ecs_cluster.main.cluster_name
    container_name = var.container_name
    task_definition_name = "${aws_ecs_task_definition.task.family}:${aws_ecs_task_definition.task.revision}"
  }
}

output public_ip {
  value = data.external.read_ecs_public_ip.result["public_ip"]
}