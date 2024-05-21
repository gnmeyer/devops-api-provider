package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "devops-bootcamp_engineer" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Engineers returned
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.#", "3"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.0.name", "Ryan"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.0.id", "H3ZTR"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.0.email", "ryan@ferrets.com"),

					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.1.name", "zach"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.1.id", "M3IGD"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.1.email", "zach@bengal.com"),

					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.2.name", "bob"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.2.id", "CTDSM"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineers.2.email", "bob@bob.com"),

					// Verify placeholder id attribute
					// resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "id", "placeholder"),
				),
			},
		},
	})
}
