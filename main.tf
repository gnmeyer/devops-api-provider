# data "devops-bootcamp_engineer" "test" {}




# output "test_engineer" {
#  value = data.devops-bootcamp_engineer.test
# }

# data "devops-bootcamp_ops" "ops_test" {}

# output "ops_test" {
#   value = data.devops-bootcamp_ops.ops_test
# }

# data "devops-bootcamp_dev" "dev_test" {}

# output "dev_test" {
#   value = data.devops-bootcamp_dev.dev_test
# }

data "devops-bootcamp_devops" "devops_test" {}

output "devops_test" {
  value = data.devops-bootcamp_devops.devops_test
}

resource "devops-bootcamp_engineer-resource" "example" {
  name  = "joeski"
  email = "myles@dudes.com"
}

resource "devops-bootcamp_engineer-resource" "example-2" {
  name  = "grant"
  email = "myleas@dudes.com"
}

output "example_engineer" {
  value = devops-bootcamp_engineer-resource.example
}

output "example_engineer-2" {
  value = devops-bootcamp_engineer-resource.example-2
}