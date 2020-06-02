########################################################################################################################
### Network

locals {
  ## Need to increase az_count to 2 at least to improve availabilty.
  az_count = "1" # Number of AZs, but only one region is used currently for DEV.
}

//data "aws_eip" "nlb_eip_id_1" {
//  id = var.nlb_eip_id_1
//}
//
//data "aws_eip" "nlb_eip_id_2" {
//  id = var.nlb_eip_id_2
//}

# Fetch AZs in the current region
data "aws_availability_zones" "available" {}

resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = local.default_tag
}

# Create var.az_count public subnets, each in a different AZ
resource "aws_subnet" "public" {
  count                   = min(length(data.aws_availability_zones.available), local.az_count)
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 8, local.az_count + count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  vpc_id                  = aws_vpc.main.id
  map_public_ip_on_launch = true

  tags = local.default_tag
}

# IGW for the public subnet
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id

  tags = local.default_tag
}

# Route the public subnet traffic through the IGW
resource "aws_route" "internet_access" {
  route_table_id         = aws_vpc.main.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.gw.id
}
