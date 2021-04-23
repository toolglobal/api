ifeq ($(mode),debug)
	LDFLAGS="-X 'main.buildTime=`date`' -X 'main.goVersion=`go version`' -X main.gitHash=`git rev-parse HEAD`"
else
	LDFLAGS="-s -w -X 'main.buildTime=`date`' -X 'main.goVersion=`go version`' -X main.gitHash=`git rev-parse HEAD`"
endif

.PHONY: build
build:
	export GOPROXY="https://goproxy.io,direct"
	rm -rf ./build && mkdir -p ./build/config
	go build -ldflags ${LDFLAGS} -o build/api cmd/main.go
	cp -r config/config.toml ./build/config
	cp -r static ./build
clean:
	rm -rf ./build

.PHONY: docs
docs:
	#swag init --parseDependency  -g cmd/api/main.go
	swag init -g cmd/main.go

