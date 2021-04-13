.PHONY: proto
proto:
	protoc \
		--proto_path=. \
		--proto_path=example/helloworld/proto \
		--go_out=${GOPATH}/src \
		--go-grpc_out=${GOPATH}/src \
		--go-complexity_out=${GOPATH}/src \
		example/helloworld/proto/helloworld.proto

.PHONY: build
build:
	cd cmd/protoc-gen-go-complexity && go build -o ${GOBIN}/protoc-gen-go-complexity .
