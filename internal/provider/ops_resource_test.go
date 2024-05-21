package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `


				resource "devops-bootcamp_engineer-resource" "ben" {
					name  = "ben"
					email = "ben@google.com"
				  }
				  resource "devops-bootcamp_engineer-resource" "wick" {
					name = "wick"
					email = "wick@google.com"
				  }
resource "devops-bootcamp_ops-resource" "test" {
	name  = "ops_example_1"
	engineers = [
		{
			id = devops-bootcamp_engineer-resource.ben.id
		},
	]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(

					// Verify ops name and amount of engineers
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "name", "ops_example_1"),
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "engineers.#", "1"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devops-bootcamp_ops-resource.test", "id"),

					// Verify the first & second engineer to ensure all attributes are set
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "engineers.0.name", "ben"),
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "engineers.0.email", "ben@google.com"),

					resource.TestCheckResourceAttrSet("devops-bootcamp_ops-resource.test", "engineers.0.id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_ops-resource.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `

				  resource "devops-bootcamp_engineer-resource" "wick" {
					name = "wick"
					email = "wick@google.com"
				  }
				  resource "devops-bootcamp_engineer-resource" "ben" {
					name  = "ben"
					email = "ben@google.com"
				  }
resource "devops-bootcamp_ops-resource" "test" {
	name  = "ops_example_2"
	engineers = [

		{
			id = devops-bootcamp_engineer-resource.wick.id
		}

	]
}

`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify ops name and amount of engineers
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "name", "ops_example_2"),
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "engineers.#", "1"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devops-bootcamp_ops-resource.test", "id"),

					// Verify the first engineer to ensure all attributes are set
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "engineers.0.name", "wick"),
					resource.TestCheckResourceAttr("devops-bootcamp_ops-resource.test", "engineers.0.email", "wick@google.com"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_ops-resource.test", "engineers.0.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
