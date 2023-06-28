---
page_title: "circleci_runner_token Resource - terraform-provider-circleci"
subcategory: ""
description: |-
  Manages a Runner token
---

# circleci_runner_token (Resource)

Manages a Runner token

## Example Usage

```terraform
resource "circleci_runner_resource_class" "machine_linux" {
  resource_class = "kelvintaywl-tf/machine-linux"
  description    = "Amazon Linux 2"

}

resource "circleci_runner_token" "from_tf" {
  resource_class = circleci_runner_resource_class.machine_linux.resource_class
  nickname       = "main"
}
```

### Machine Runner setup with AWS EC2 instance

This is a sample for your reference.

Please modify accordingly.

```terraform
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
```

Here is a sample of the userdata.yaml file:

```yaml
#cloud-config
packages:
  - docker
  - git
  - rpm-build
  - policycoreutils-devel
# NOTE: this is for Amazon Linux 2023 (amd64)
runcmd:
  - 'export platform="${platform}"'
  - 'export agent_version="${agent_version}"'
  - 'sudo curl https://raw.githubusercontent.com/CircleCI-Public/runner-installation-files/main/download-launch-agent.sh --output ./download-launch-agent.sh'
  - 'sh ./download-launch-agent.sh'
  - 'id -u circleci &>/dev/null || sudo adduser -c GECOS circleci'
  - 'sudo mkdir -p /var/opt/circleci'
  - 'sudo chmod 0750 /var/opt/circleci'
  - 'sudo chown -R circleci /var/opt/circleci /opt/circleci'
  - 'echo "circleci ALL=(ALL) NOPASSWD:ALL" | sudo tee -a /etc/sudoers'
  - 'sudo mkdir -p /etc/opt/circleci && sudo touch /etc/opt/circleci/launch-agent-config.yaml'
  - 'sudo wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq'
  - 'chmod +x /usr/bin/yq'
  - 'yq e ".api.auth_token = \"${token}\"" -i /etc/opt/circleci/launch-agent-config.yaml'
  - 'yq e ".api.url = \"https://${hostname}\"" -i /etc/opt/circleci/launch-agent-config.yaml'
  - 'yq e ".runner.name = \"${name}\"" -i /etc/opt/circleci/launch-agent-config.yaml'
  - 'yq e ".runner.working_directory = \"/var/opt/circleci/workdir\"" -i /etc/opt/circleci/launch-agent-config.yaml'
  - 'yq e ".runner.cleanup_working_directory = true" -i /etc/opt/circleci/launch-agent-config.yaml'
  - 'sudo chown -R circleci: /etc/opt/circleci'
  - 'sudo chmod 600 /etc/opt/circleci/launch-agent-config.yaml'
  - 'sudo mkdir -p /etc/opt/circleci/policy'
  - 'sudo sepolicy generate --path /etc/opt/circleci/policy --init /opt/circleci/circleci-launch-agent'
  - 'sudo curl https://raw.githubusercontent.com/CircleCI-Public/runner-installation-files/main/rhel8-install/circleci_launch_agent.te --output /etc/opt/circleci/policy/circleci_launch_agent.te'
  - 'sudo /etc/opt/circleci/policy/circleci_launch_agent.sh'
  - 'sudo /opt/circleci/circleci-launch-agent --config /etc/opt/circleci/launch-agent-config.yaml'
```

### Print token raw in output

**Note:** This is not recommended!

However, if you need to, you can print the token (sensitive) out in the output.

```terraform
resource "circleci_runner_token" "my_new_token" {
  resource_class = "kelvintaywl-tf/test"
  nickname       = "my-new-token"
}

output "token" {
  description = "kelvintaywl-tf/test runner token (my-new-token)"
  # to print token raw in the output
  value = nonsensitive(circleci_runner_token.my_new_token.token)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `nickname` (String) The Runner token alias.
- `resource_class` (String) The name of the Runner resource-class (should include namespace)

### Read-Only

- `created_at` (String) Date and time the token was created
- `id` (String) The unique ID of the Runner token.
- `token` (String, Sensitive) The Runner token value.
