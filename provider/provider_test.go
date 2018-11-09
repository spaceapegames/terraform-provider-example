package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"example": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SERVICE_ADDRESS"); v == "" {
		t.Fatal("SERVICE_ADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICE_PORT"); v == "" {
		t.Fatal("SERVICE_PORT must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICE_TOKEN"); v == "" {
		t.Fatal("SERVICE_TOKEN must be set for acceptance tests")
	}
}
