#!/bin/bash
go build -o cloudpilot-build cmd/web/*.go && ./cloudpilot-build -dbname=cloudpilot -dbuser=postgres -cache=false -production=false