ARCH=amd64
OS=linux
CGO_FLAG=0
OUTPUT=webserver

all:
	CGO_ENABLE=${CGO_FLAG} GOOS=${OS} GOARCH=${ARCH} go build -o ${OUTPUT}

t := "/tmp/go-cover.$(shell /bin/bash -c "date +%Y%m%d%H%M%S").tmp"

test:
	podman run -d -p 27017:27017 --name mongo-example docker.io/library/mongo:latest
	GIN_MODE=test bash -c 'go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t'
	podman stop mongo-example
	podman rm mongo-example

clean:
	rm ${OUTPUT}
