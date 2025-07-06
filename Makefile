.PHONY: build

# 根据当前操作系统自动选择编译方式
build:
	ifeq ($(OS),Windows_NT)
		$env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="1"; go build -o bin/narcissus-blog-linux main.go
	else ifeq ($(shell uname),Darwin)
		# $env:GOOS="darwin"; $env:GOARCH="amd64"; $env:CGO_ENABLED="1"; go build -o bin/narcissus-blog-darwin main.go
	else
		GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/narcissus-blog-linux main.go
	endif
