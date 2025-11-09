.PHONY: build openapi upload

# 根据当前操作系统自动选择编译方式
build:
	ifeq ($(OS),Windows_NT)
		$env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="1"; go build -o bin/narcissus-blog-linux main.go
	else ifeq ($(shell uname),Darwin)
		# $env:GOOS="darwin"; $env:GOARCH="amd64"; $env:CGO_ENABLED="1"; go build -o bin/narcissus-blog-darwin main.go
	else
		GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/narcissus-blog-linux main.go
	endif

upload:
	scp -r bin/narcissus-blog-linux narcissus:/home/narcissus/workspace/deployment
#	scp -r conf narcissus:/home/narcissus/workspace/deployment

openapi:
  swag init -g ./cmd/blog/main.go -o docs

db-model-generate:
  # TODO
  # "root:root@tcp(172.26.21.6:3306)/blog_narcissus?charset=utf8mb4&parseTime=true&loc=Local"