package provider

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlert_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExampleItemImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleItemExists("example_item.test_import"),
				),
			},
			{
				ResourceName:      "example_item.test_import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckExampleItemImporter_basic() string {
	return fmt.Sprintf(`
resource "example_item" "test_import" {
	name = "test_import"
	description = "testing importing of a resource"
	tags = [
		"import1",
		"import2",
		"import3"
	]
}
`)
}
