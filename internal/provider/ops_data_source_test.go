package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "devops-bootcamp_ops" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Engineers returned
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.#", "2"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.0.name", "ops_ferrets"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.0.id", "MIGFP"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.0.engineers.#", "0"),

					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.1.name", "ops_bengal"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.1.id", "YBTQO"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_ops.test", "ops.1.engineers.#", "0"),

					// Verify placeholder id attribute
					// resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "id", "placeholder"),
				),
			},
		},
	})
}
