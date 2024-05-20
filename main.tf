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

# data "devops-bootcamp_devops" "devops_test" {}

# output "devops_test" {
#   value = data.devops-bootcamp_devops.devops_test
# }

resource "devops-bootcamp_engineer-resource" "grant" {
  name  = "grant"
  email = "grant@google.com"
}

resource "devops-bootcamp_engineer-resource" "jocko" {
  name  = "jocko"
  email = "jocko@google.com"
}

resource "devops-bootcamp_engineer-resource" "ben" {
  name  = "ben"
  email = "ben@google.com"
}

resource "devops-bootcamp_engineer-resource" "myles" {
  name = "myles"
  email = "myles@google.com"
}
resource "devops-bootcamp_engineer-resource" "wick" {
  name = "wick"
  email = "wick@google.com"
}


resource "devops-bootcamp_ops-resource" "example" {
  name = "example-ops"
  engineers = [
      {
      id = devops-bootcamp_engineer-resource.grant.id
    },
    {
      id = devops-bootcamp_engineer-resource.jocko.id
    },
    {
     id = devops-bootcamp_engineer-resource.myles.id
    },
    {
      id = devops-bootcamp_engineer-resource.jocko.id
    },
    # {
    #   id = devops-bootcamp_engineer-resource.wick.id
    # }
  ]

}

# output "example_engineer" {
#   value = devops-bootcamp_engineer-resource.example
# }

# output "example_engineer-2" {
#   value = devops-bootcamp_engineer-resource.example-2
# }