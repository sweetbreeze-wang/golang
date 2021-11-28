export tag=v2.0
root:
	export ROOT=github.com/cncamp/golang

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t httpserver:${tag} .
tag:
	echo "make new tag for docker push"
	docker tag  sweetbreeze0207/httpserver:v2.0 httpserver:v2.0

push: release
	echo "pushing cncamp/httpserver"
	docker push sweetbreeze0207/httpserver 
run:
	echo "run httpservevr"
	docker run  -p 192.168.0.105:80:80 -d httpserver:v2.0
