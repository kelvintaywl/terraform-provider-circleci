resource "circleci_runner_resource_class" "machine_linux" {
  resource_class = "kelvintaywl-tf/machine-linux"
  description    = "Amazon Linux 2"

}

resource "circleci_runner_token" "from_tf" {
  resource_class = circleci_runner_resource_class.machine_linux.resource_class
  nickname       = "main"
}
