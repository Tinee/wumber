.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./hello-world/hello-world

local: 
	sam local start-api --env-vars env.development.json
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/create-workspace ./cmd/workspace/create 
	GOOS=linux GOARCH=amd64 go build -o ./bin/jwt-auth ./cmd/auth/jwt
	GOOS=linux GOARCH=amd64 go build -o ./bin/user-register ./cmd/user/register
	GOOS=linux GOARCH=amd64 go build -o ./bin/mixpanel-user-register ./cmd/mixpanel/user-created