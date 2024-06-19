.PHONY: generate-backend generate-backend-movies generate-backend-interactions

generate-backend: generate-backend-movies generate-backend-interactions

generate-backend-movies:
	rm -rf pkg/movies_v1
	mkdir -p pkg/movies_v1
	protoc --go_out=pkg/movies_v1 --go-grpc_out=pkg/movies_v1 \
    	protos/movies_v1/movies_v1.proto
	mv pkg/movies_v1/movies_v1/* pkg/movies_v1
	rm -rf pkg/movies_v1/movies_v1/

generate-backend-interactions:
	rm -rf pkg/interactions_v1
	mkdir -p pkg/interactions_v1
	protoc --go_out=pkg/interactions_v1 --go-grpc_out=pkg/interactions_v1 \
		protos/interactions_v1/interactions_v1.proto
	mv pkg/interactions_v1/interactions_v1/* pkg/interactions_v1
	rm -rf pkg/interactions_v1/interactions_v1/


generate-frontend: generate-frontend-movies generate-frontend-interactions

generate-frontend-movies:
	rm -rf web/kinogo/src/generated/movies_v1
	mkdir -p web/kinogo/src/generated/movies_v1
	protoc --ts_out=web/kinogo/src/generated/movies_v1 \
    	protos/movies_v1/movies_v1.proto
	mv web/kinogo/src/generated/movies_v1/movies_v1/* web/kinogo/src/generated/movies_v1
	rm -rf web/kinogo/src/generated/movies_v1/movies_v1/

generate-frontend-interactions:
	rm -rf web/kinogo/src/generated/interactions_v1
	mkdir -p web/kinogo/src/generated/interactions_v1
	protoc --ts_out=web/kinogo/src/generated/interactions_v1 \
		protos/interactions_v1/interactions_v1.proto
	mv web/kinogo/src/generated/interactions_v1/interactions_v1/* web/kinogo/src/generated/interactions_v1
	rm -rf web/kinogo/src/generated/interactions_v1/interactions_v1/