package client

import (
	"fmt"
	"github.com/phayes/freeport"
	"github.com/spaceapegames/terraform-provider-example/api/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_GetAll(t *testing.T) {
	testCases := []struct {
		testName  string
		seedData  map[string]server.Item
		expectErr bool
	}{
		{
			testName:  "no items",
			seedData:  nil,
			expectErr: false,
		},
		{
			testName: "one item",
			seedData: map[string]server.Item{
				"one": {
					Name:        "one",
					Description: "I'm an item",
					Tags:        nil,
				},
			},
			expectErr: false,
		},
		{
			testName: "more items",
			seedData: map[string]server.Item{
				"one": {
					Name:        "one",
					Description: "I'm an item",
					Tags:        nil,
				},
				"two": {
					Name:        "two",
					Description: "I'm an item",
					Tags:        nil,
				},
				"three": {
					Name:        "three",
					Description: "I'm an item",
					Tags:        nil,
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			port, testServer, err := setupServer(tc.seedData)
			assert.NoError(t, err, "error setting up server for test")
			go func() {
				err = testServer.ListenAndServe()
				assert.NoError(t, err, "error starting server for test")
			}()
			client := NewClient("http://localhost", port, "supersecrettoken")

			items, err := client.GetAll()
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tc.seedData == nil {
				tc.seedData = map[string]server.Item{}
			}
			assert.Equal(t, tc.seedData, *items)
		})
	}
}

func TestClient_GetItem(t *testing.T) {
	testCases := []struct {
		testName     string
		itemName     string
		seedData     map[string]server.Item
		expectErr    bool
		expectedResp *server.Item
	}{
		{
			testName: "item exists",
			itemName: "item1",
			seedData: map[string]server.Item{
				"item1": {
					Name:        "item1",
					Description: "an item",
					Tags:        nil,
				},
			},
			expectErr: false,
			expectedResp: &server.Item{
				Name:        "item1",
				Description: "an item",
				Tags:        nil,
			},
		},
		{
			testName: "item with tags",
			itemName: "item1",
			seedData: map[string]server.Item{
				"item1": {
					Name:        "item1",
					Description: "an item",
					Tags: []string{
						"tag1",
						"tag2",
					},
				},
			},
			expectErr: false,
			expectedResp: &server.Item{
				Name:        "item1",
				Description: "an item",
				Tags: []string{
					"tag1",
					"tag2",
				},
			},
		},
		{
			testName:     "item does not exist",
			itemName:     "item1",
			seedData:     nil,
			expectErr:    true,
			expectedResp: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			port, testServer, err := setupServer(tc.seedData)
			assert.NoError(t, err, "error setting up server for test")
			go func() {
				err = testServer.ListenAndServe()
				assert.NoError(t, err, "error starting server for test")
			}()
			client := NewClient("http://localhost", port, "supersecrettoken")

			item, err := client.GetItem(tc.itemName)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, item)
		})
	}
}

func TestClient_NewItem(t *testing.T) {
	testCases := []struct {
		testName  string
		newItem   *server.Item
		seedData  map[string]server.Item
		expectErr bool
	}{
		{
			testName: "success",
			newItem: &server.Item{
				Name:        "item1",
				Description: "describe me",
				Tags:        nil,
			},
			seedData:  nil,
			expectErr: false,
		},
		{
			testName: "item already exists",
			newItem: &server.Item{
				Name:        "item1",
				Description: "describe me",
				Tags:        nil,
			},
			seedData: map[string]server.Item{
				"item1": {
					Name:        "item1",
					Description: "describe me",
					Tags:        nil,
				},
			},
			expectErr: true,
		},
		{
			testName: "whitespace in item name",
			newItem: &server.Item{
				Name:        "item 1",
				Description: "describe me 1",
				Tags:        nil,
			},
			seedData:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			port, testServer, err := setupServer(tc.seedData)
			assert.NoError(t, err, "error setting up server for test")
			go func() {
				err = testServer.ListenAndServe()
				assert.NoError(t, err, "error starting server for test")
			}()
			client := NewClient("http://localhost", port, "supersecrettoken")

			err = client.NewItem(tc.newItem)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			item, err := client.GetItem(tc.newItem.Name)
			assert.NoError(t, err)
			assert.Equal(t, tc.newItem, item)
		})
	}
}

func TestClient_UpdateItem(t *testing.T) {
	testCases := []struct {
		testName    string
		updatedItem *server.Item
		seedData    map[string]server.Item
		expectErr   bool
	}{
		{
			testName: "item exists",
			updatedItem: &server.Item{
				Name:        "item1",
				Description: "describe me 1",
				Tags:        nil,
			},
			seedData: map[string]server.Item{
				"item1": {
					Name:        "item1",
					Description: "describe me",
					Tags:        nil,
				},
			},
			expectErr: false,
		},
		{
			testName: "item does not exist",
			updatedItem: &server.Item{
				Name:        "item1",
				Description: "describe me 1",
				Tags:        nil,
			},
			seedData:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			port, testServer, err := setupServer(tc.seedData)
			assert.NoError(t, err, "error setting up server for test")
			go func() {
				err = testServer.ListenAndServe()
				assert.NoError(t, err, "error starting server for test")
			}()
			client := NewClient("http://localhost", port, "supersecrettoken")

			err = client.UpdateItem(tc.updatedItem)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			item, err := client.GetItem(tc.updatedItem.Name)
			assert.NoError(t, err)
			assert.Equal(t, tc.updatedItem, item)
		})
	}
}

func TestClient_DeleteItem(t *testing.T) {
	testCases := []struct {
		testName  string
		itemName  string
		seedData  map[string]server.Item
		expectErr bool
	}{
		{
			testName: "item exists",
			itemName: "item1",
			seedData: map[string]server.Item{
				"item1": {
					Name:        "item1",
					Description: "describe me",
					Tags:        nil,
				},
			},
			expectErr: false,
		},
		{
			testName:  "item does not exist",
			itemName:  "item1",
			seedData:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			port, testServer, err := setupServer(tc.seedData)
			assert.NoError(t, err, "error setting up server for test")
			go func() {
				err = testServer.ListenAndServe()
				assert.NoError(t, err, "error starting server for test")
			}()
			client := NewClient("http://localhost", port, "supersecrettoken")

			err = client.DeleteItem(tc.itemName)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			_, err = client.GetItem(tc.itemName)
			assert.Error(t, err)
		})
	}
}

// setupServer sets up a new api server, on localhost, on a random port, return the port, the service and any errors
func setupServer(seedData map[string]server.Item) (int, *server.Service, error) {
	if seedData == nil {
		seedData = map[string]server.Item{}
	}
	port, err := freeport.GetFreePort()
	if err != nil {
		return 0, nil, err
	}
	return port, server.NewService(fmt.Sprintf("localhost:%v", port), seedData), err
}
