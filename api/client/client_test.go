package client

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestClient_GetAll(t *testing.T) {
	client := NewClient("http://localhost", "3001", "")
	resp, err := client.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(spew.Sdump(resp))
}
