#!/bin/bash

VERSION=$(cat version)
echo "building terraform-provider-example_${VERSION}"
go build -o terraform-provider-example_${VERSION}