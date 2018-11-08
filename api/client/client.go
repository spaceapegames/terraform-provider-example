package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spaceapegames/terraform-provider-example/api/server"
	"io"
	"net/http"
)

// Client holds all of the information required to connect to a server
type Client struct {
	hostname   string
	port       int
	authToken  string
	httpClient *http.Client
}

// NewClient returns a new client configured to communicate on a server with the
// given hostname and port and to send an Authorization Header with the value of
// token
func NewClient(hostname string, port int, token string) *Client {
	return &Client{
		hostname:   hostname,
		port:       port,
		authToken:  token,
		httpClient: &http.Client{},
	}
}

// GetAll Retrieves all of the Items from the server
func (c *Client) GetAll() (*map[string]server.Item, error) {
	body, err := c.httpRequest("item", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	items := map[string]server.Item{}
	err = json.NewDecoder(body).Decode(&items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

// GetItem gets an item with a specific name from the server
func (c *Client) GetItem(name string) (*server.Item, error) {
	body, err := c.httpRequest(fmt.Sprintf("item/%v", name), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	item := &server.Item{}
	err = json.NewDecoder(body).Decode(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// NewItem creates a new Item
func (c *Client) NewItem(item *server.Item) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return err
	}
	_, err = c.httpRequest("item", "POST", buf)
	if err != nil {
		return err
	}
	return nil
}

// UpdateItem updates the values of an item
func (c *Client) UpdateItem(item *server.Item) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("item/%s", item.Name), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem removes an item from the server
func (c *Client) DeleteItem(itemName string) error {
	_, err := c.httpRequest(fmt.Sprintf("item/%s", itemName), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authToken)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%s:%v/%s", c.hostname, c.port, path)
}
