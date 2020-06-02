//data aws_lb "nlb" {
//  count = local.lb_exist
//  arn  = var.nlb_arn
//}
//
//resource "aws_lb_target_group" "tcp" {
//  for_each = local.tcp_ports
//
//  name     = "lbtg-tcp-${each.key}"
//  port     = each.value
//  protocol = "TCP"
//  target_type = "ip"
//  vpc_id   = data.aws_vpc.main.id
//
//  health_check {
//    protocol            = "TCP"
//    port                = each.value
//    interval            = 10
//    healthy_threshold   = 3
//    unhealthy_threshold = 3
//  }
//
//  stickiness {
//    enabled = false
//    type = "lb_cookie"
//  }
//}
//
//resource "aws_lb_listener" "tcp_port" {
//  for_each = local.tcp_ports
//
//  load_balancer_arn = data.aws_lb.nlb[0].arn
//  port              = each.value
//  protocol          = "TCP"
//
//  default_action {
//    type             = "forward"
//    target_group_arn = aws_lb_target_group.tcp[each.key].arn
//  }
//}
//
//resource "aws_lb_target_group" "tls" {
//  for_each = local.tls_ports
//
//  name     = "lbtg-tls-${each.key}"
//  port     = each.value
//  protocol = "TLS"
//  target_type = "ip"
//  vpc_id   = data.aws_vpc.main.id
//
//  health_check {
//    protocol            = "TCP"
//    interval            = 30
//    healthy_threshold   = 5
//    unhealthy_threshold = 5
//  }
//}
//
//data "aws_route53_zone" "zone" {
//  count = var.zone_name != null ? 1 : 0
//  name = var.zone_name
//}
//
//resource "aws_acm_certificate" "cert" {
//  count = var.cert_domain != null ? 1 : 0
//
//  domain_name       = var.cert_domain
//  validation_method = "DNS"
//}
//
//resource "aws_route53_record" "record" {
//  count = var.cert_domain != null ? 1 : 0
//
//  name    = aws_acm_certificate.cert[0].domain_validation_options[0].resource_record_name
//  type    = aws_acm_certificate.cert[0].domain_validation_options[0].resource_record_type
//  zone_id = data.aws_route53_zone.zone[0].zone_id
//  records = [aws_acm_certificate.cert[0].domain_validation_options[0].resource_record_value]
//  ttl     = 60
//}
//
//resource "aws_acm_certificate_validation" "cert" {
//  count = var.cert_domain != null ? 1 : 0
//
//  certificate_arn         = aws_acm_certificate.cert[0].arn
//  validation_record_fqdns = [
//    aws_route53_record.record[0].fqdn
//  ]
//}
//
//resource "aws_lb_listener" "tls_port" {
//  for_each = local.tls_ports
//
//  load_balancer_arn = data.aws_lb.nlb[0].arn
//  port              = each.value
//  protocol          = "TLS"
//  ssl_policy        = "ELBSecurityPolicy-2016-08"
//  certificate_arn   = aws_acm_certificate_validation.cert[0].certificate_arn
//
//  default_action {
//    type             = "forward"
//    target_group_arn = aws_lb_target_group.tls[each.key].arn
//  }
//}
