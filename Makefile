.PHONY: generate generate-backend generate-backend-movies generate-backend-interactions generate-backend-comments generate-frontend generate-frontend-movies generate-frontend-interactions generate-frontend-comments
generate: generate-backend generate-frontend

generate-backend: generate-backend-movies generate-backend-interactions generate-backend-comments

generate-backend-movies:
	rm -rf backend/pkg/movies_v1
	mkdir -p backend/pkg/movies_v1
	protoc --go_out=backend/pkg/movies_v1 --go-grpc_out=backend/pkg/movies_v1 \
    	protos/movies_v1/movies_v1.proto
	mv backend/pkg/movies_v1/movies_v1/* backend/pkg/movies_v1
	rm -rf backend/pkg/movies_v1/movies_v1/

generate-backend-interactions:
	rm -rf backend/pkg/interactions_v1
	mkdir -p backend/pkg/interactions_v1
	protoc --go_out=backend/pkg/interactions_v1 --go-grpc_out=backend/pkg/interactions_v1 \
		protos/interactions_v1/interactions_v1.proto
	mv backend/pkg/interactions_v1/interactions_v1/* backend/pkg/interactions_v1
	rm -rf backend/pkg/interactions_v1/interactions_v1/

generate-backend-comments:
	rm -rf backend/pkg/comments_v1
	mkdir -p backend/pkg/comments_v1
	protoc --go_out=backend/pkg/comments_v1 --go-grpc_out=backend/pkg/comments_v1 \
    	protos/comments_v1/comments_v1.proto
	mv backend/pkg/comments_v1/comments_v1/* backend/pkg/comments_v1
	rm -rf backend/pkg/comments_v1/comments_v1/

generate-frontend: generate-frontend-movies generate-frontend-interactions generate-frontend-comments

generate-frontend-movies:
	protoc -I=. --plugin="protoc-gen-ts=frontend/node_modules/.bin/protoc-gen-ts" --ts_out=frontend/src \
		protos/movies_v1/movies_v1.proto

generate-frontend-interactions:
	protoc -I=. --plugin="protoc-gen-ts=frontend/node_modules/.bin/protoc-gen-ts" --ts_out=frontend/src \
    		protos/interactions_v1/interactions_v1.proto

generate-frontend-comments:
	protoc -I=. --plugin="protoc-gen-ts=frontend/node_modules/.bin/protoc-gen-ts" --ts_out=frontend/src \
    		protos/comments_v1/comments_v1.proto