package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer-resource" "test" {
	name  = "grant"
	email = "grant@google.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(

					// Verify first order item
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "name", "grant"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "email", "grant@google.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer-resource.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_engineer-resource.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer-resource" "test" {
	name  = "joe"
	email = "joe@google.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(

					// Verify first order item
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "name", "joe"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "email", "joe@google.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer-resource.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
