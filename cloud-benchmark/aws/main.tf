/*
   Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

provider "aws" {
  region = "${var.region}"
}

data "aws_vpc" "default" {
  default = true
}

resource "aws_key_pair" "qed-benchmark" {
  key_name   = "qed-benchmark"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}

data "aws_subnet_ids" "all" {
  vpc_id = "${data.aws_vpc.default.id}"
}

data "aws_ami" "amazon_linux" {
  most_recent = true

  filter {
    name = "name"

    values = [
      "amzn-ami-hvm-*-x86_64-gp2",
    ]
  }

  filter {
    name = "owner-alias"

    values = [
      "amazon",
    ]
  }
}

module "security_group" {
  source = "terraform-aws-modules/security-group/aws"

  name        = "qed-benchmark"
  description = "Security group for QED benchmark usage"
  vpc_id      = "${data.aws_vpc.default.id}"

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["http-8080-tcp", "all-icmp", "ssh-tcp" ]
  egress_rules        = ["all-all"]
}

resource "aws_security_group_rule" "allow_profiling" {
  type            = "ingress"
  from_port       = 6060
  to_port         = 6060
  protocol        = "tcp"
  cidr_blocks     = ["0.0.0.0/0"]

  security_group_id = "${module.security_group.this_security_group_id}"
}


resource "aws_eip" "qed-benchmark" {
  vpc      = true
  instance = "${module.ec2.id[0]}"
}

module "ec2" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name                        = "qed-benchmark"
  ami                         = "${data.aws_ami.amazon_linux.id}"
  instance_type               = "${var.flavour}"
  subnet_id                   = "${element(data.aws_subnet_ids.all.ids, 0)}"
  vpc_security_group_ids      = ["${module.security_group.this_security_group_id}"]
  associate_public_ip_address = true
  key_name                    = "${aws_key_pair.qed-benchmark.key_name}"

  root_block_device = [{
    volume_type = "gp2"
    volume_size = "${var.volume_size}"
    delete_on_termination = true
  }]
}