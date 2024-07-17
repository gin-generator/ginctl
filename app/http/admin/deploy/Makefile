export tag=prod-20240706-1
export GOPROXY=https://goproxy.cn,direct
build:
	echo "building admin binary"
	go build -o admin .

build_image:
	echo "building admin container"
	docker build -t your-registry:${tag} .