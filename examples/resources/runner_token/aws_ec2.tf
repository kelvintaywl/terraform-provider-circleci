locals {
  circleci_namespace           = "acmeorg"
  circleci_org_vcs             = "github"
  circleci_org_name            = "acmeorg"
  circleci_runner_machine_name = "aws_ec2_linux2023"

  aws_ec2_ami_id        = "ami-xxxx"
  aws_ec2_instance_type = "t3.medium"
  # ASSUMPTION: existing VPC with subnets and security group already provisioned
  aws_vpc_security_group_id = "sg-xxxx"
  aws_vpc_subnet_id         = "subnet"
}

# Creates the org namespace via the CircleCI CLI locally
resource "null_resource" "circleci_namespace" {
  provisioner "local-exec" {
    command = "circleci --host https://${var.circleci_hostname} --token ${var.circleci_api_token}  namespace create ${local.circleci_namespace} ${local.circleci_org_vcs} ${local.circleci_org_name} --no-prompt"
  }
}

resource "circleci_runner_resource_class" "machine_linux" {
  resource_class = "${local.circleci_namespace}/${local.circleci_runner_machine_name}"
  description    = "Amazon Linux 2023"
  depends_on     = [null_resource.circleci_namespace]
}

resource "circleci_runner_token" "admin" {
  resource_class = circleci_runner_resource_class.machine_linux.resource_class
  nickname       = "admin"
}

resource "aws_key_pair" "aws_key_pair" {
  key_name   = "work laptop ED25519 SSH key"
  public_key = file("/local/path/to/specific/id_ed25519.pub")
}

resource "aws_instance" "aws_ec2_linux_2023" {
  # example: creating 2 instances
  count = 2

  ami           = local.aws_ec2_ami_id
  instance_type = local.aws_ec2_instance_type
  user_data = templatefile(
    "${path.module}/userdata.yaml",
    {
      token         = circleci_runner_token.admin.token
      name          = "aws_ec2_linux_2023_${count.index}"
      hostname      = var.circleci_hostname
      agent_version = "1.0.48283-583799b"
      platform      = "linux/amd64"
    }
  )
  key_name                    = aws_key_pair.aws_key_pair.key_name
  security_groups             = [local.aws_vpc_security_group_id]
  subnet_id                   = local.aws_vpc_subnet_id
  user_data_replace_on_change = true

  tags = {
    # name of EC2 instance
    Name = "${local.circleci_runner_machine_name}-${count.index}"
  }
}
