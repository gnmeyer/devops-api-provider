# data "devops-bootcamp_engineer" "test" {}

# output "test_engineer" {
#  value = data.devops-bootcamp_engineer.test
# }

resource "devops-bootcamp_engineer-resource" "example" {
  name  = "engineer_name_222"
  email = "myles@dudes.com"
}

resource "devops-bootcamp_engineer-resource" "example-2" {
  name  = "engineer_name_2223"
  email = "myleas@dudes.com"
}

output "example_engineer" {
  value = devops-bootcamp_engineer-resource.example
}

output "example_engineer-2" {
  value = devops-bootcamp_engineer-resource.example-2
}