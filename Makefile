PHONY: generate
generate:
	rm -rf pkg/movies_v1
	mkdir -p pkg/movies_v1
	protoc --go_out=pkg/movies_v1 --go-grpc_out=pkg/movies_v1 \
		   protos/movies_v1/service.proto
	mv pkg/movies_v1/movies_v1/* pkg/movies_v1
	rm -rf pkg/movies_v1/movies_v1