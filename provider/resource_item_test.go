package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spaceapegames/terraform-provider-example/api/client"
	"regexp"
	"testing"
)

func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleItemExists("example_item.test_item"),
					resource.TestCheckResourceAttr(
						"example_item.test_item", "name", "test"),
					resource.TestCheckResourceAttr(
						"example_item.test_item", "description", "hello"),
					resource.TestCheckResourceAttr(
						"example_item.test_item", "tags.#", "2"),
					resource.TestCheckResourceAttr("example_item.test_item", "tags.1931743815", "tag1"),
					resource.TestCheckResourceAttr("example_item.test_item", "tags.1477001604", "tag2"),
				),
			},
		},
	})
}

func TestAccItem_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleItemExists("example_item.test_update"),
					resource.TestCheckResourceAttr(
						"example_item.test_update", "name", "test_update"),
					resource.TestCheckResourceAttr(
						"example_item.test_update", "description", "hello"),
					resource.TestCheckResourceAttr(
						"example_item.test_update", "tags.#", "2"),
					resource.TestCheckResourceAttr("example_item.test_update", "tags.1931743815", "tag1"),
					resource.TestCheckResourceAttr("example_item.test_update", "tags.1477001604", "tag2"),
				),
			},
			{
				Config: testAccCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleItemExists("example_item.test_update"),
					resource.TestCheckResourceAttr(
						"example_item.test_update", "name", "test_update"),
					resource.TestCheckResourceAttr(
						"example_item.test_update", "description", "updated description"),
					resource.TestCheckResourceAttr(
						"example_item.test_update", "tags.#", "2"),
					resource.TestCheckResourceAttr("example_item.test_update", "tags.1931743815", "tag1"),
					resource.TestCheckResourceAttr("example_item.test_update", "tags.1477001604", "tag2"),
				),
			},
		},
	})
}

func TestAccItem_Multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemMultiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleItemExists("example_item.test_item"),
					testAccCheckExampleItemExists("example_item.another_item"),
				),
			},
		},
	})
}

var whiteSpaceRegex = regexp.MustCompile("name cannot contain whitespace")

func TestAccItem_WhitespaceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckItemWhitespace(),
				ExpectError: whiteSpaceRegex,
			},
		},
	})
}

func testAccCheckItemDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "example_item" {
			continue
		}

		_, err := apiClient.GetItem(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testAccCheckExampleItemExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		apiClient := testAccProvider.Meta().(*client.Client)
		_, err := apiClient.GetItem(name)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckItemBasic() string {
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

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
resource "example_item" "test_update" {
  name        = "test_update"
  description = "hello"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
resource "example_item" "test_update" {
  name        = "test_update"
  description = "updated description"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
`)
}

func testAccCheckItemMultiple() string {
	return fmt.Sprintf(`
resource "example_item" "test_item" {
  name        = "test"
  description = "hello"
  
  tags = [
	"tag1",
    "tag2",
  ]
}

resource "example_item" "another_item" {
	name        = "another_test"
	description = "hello"
	
	tags = [
	  "tag1",
	  "tag2",
	]
  }
`)
}

func testAccCheckItemWhitespace() string {
	return fmt.Sprintf(`
resource "example_item" "test_item" {
	name        = "test with whitespace"
	description = "hello"

	tags = [
		"tag1",
		"tag2",
	]
}
`)
}
