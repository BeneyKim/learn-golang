locals {
  // pbx_record_name = "pbx"
  cloudsim_record_name = "cloudsim"
}

data "aws_route53_zone" "ntel" {
  name = "ntel.ml."
}

resource "aws_route53_record" "cloudsim" {
  zone_id = data.aws_route53_zone.ntel.zone_id
  name    = "${local.cloudsim_record_name}.${data.aws_route53_zone.ntel.name}"
  type    = "A"
  ttl     = 300
  records = [
    module.cloudsim.public_ip
  ]
}
