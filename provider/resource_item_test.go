package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spaceapegames/terraform-provider-blog/api/client"
	"testing"
)

func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItem_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("example_item.test_item"),
					resource.TestCheckResourceAttr(
						"example_item.test_item", "name", "test"),
					resource.TestCheckResourceAttr(
						"example_item.test_item", "description", "hello"),
					resource.TestCheckResourceAttr(
						"example_item.test_item", "tags.#", "2"),
					// TODO check tag values are correct
				),
			},
		},
	})
}

func testAccCheckItemDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckWavefrontAlertExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		fmt.Println(rs.Primary)

		name := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*client.Client)
		_, err := apiClient.GetItem(name)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckItem_basic() string {
	return fmt.Sprintf(`
resource "example_item" "test_item" {
  name        = "test"
  description = "hello"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
`)
}
