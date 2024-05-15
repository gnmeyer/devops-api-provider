data "devops-bootcamp_engineer" "test" {}

output "test_engineer" {
  value = data.devops-bootcamp_engineer.test
}