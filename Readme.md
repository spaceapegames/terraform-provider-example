# Building a Terraform Provider

Code to accompany the [Building a Terraform Provider](https://medium.com/spaceapetech/creating-a-terraform-provider-part-1-ed12884e06d7) blog.

Consists of several components

*  A main.go which serves as the entry point to the provider
*  A provider package which implments the provider and is consumed by main.go
*  An api package which contains of a main.go which is the entry point to the server. This would not usually live within the same repository as the provider code, it's just here so that all the code for this example lives with in a single repository
    *  The api consists of two packages:
        *  server, which is the implementation of the webserver
        *  client, which is a client that can be used to programatically interact with the server.

## Requirements

* go => 1.11

This project used Go Modules, so you will need to enable them using `export GO111MODULE=on`, otherwise your go commands (run, build and test) will fail.

## API

The API is pretty simple, it just stores items which have a name, description and some tags, tags are a slice of strings. Name serves as the id for the Item. 

``` go
type Item struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
```

### Routes

All Items are stored in memeory in a `map[string]Item`, where the key is the name of the Item.

The server has five routes:

*  POST /item  - Create an item
*  GET /item - Retrive all of the items
*  GET /item/{name} - Retrieve a single item by name
*  PUT /item/{name} - Update a single item by name
*  DELETE /item/{name} - Delete a single item by name

### Starting the Server

You can start the server by running `go run api/main.go` or `make startapi` from the root of the repository. This will start the server on `localhost:3001`

You can optionally provide a file containing json to seed the server by providing a seed flag; `go run api/main.go -seed seed.json`

### Authentication

An non-empty `Authorization` header must be provided with all requests. The server will reject any requests without this.

## Client

The client can be used to programatically interact with the Server and is what the provider will use.

There is a `NewClient` function that will return a `*Client`. The function takes a hostname, port and token (The token can be anything that is not an empty string).

This will create a client for server with the default, hard-coded settings:

``` go
client := NewClient("http://localhost", 3001, "supersecrettoken")
```

There are then 5 methods, GetAll, GetItem, NewItem, UpdateItem and DeleteItem, which map to the api endpoints of the server.