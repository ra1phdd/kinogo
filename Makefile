.PHONY: test cover

test:
	cd backend && go test -race -count 1 ./...

cover:
	cd backend && go test -short -race -count 1 -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out
	rm backend/coverage.out

.PHONY: build build-backend build-frontend

build: build-backend build-frontend

build-backend:
	mkdir -p backend/build
	rm -rf build/*
	cd backend && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/kinogo_amd64 ./cmd/main/main.go
	cd backend && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o build/kinogo_arm64 ./cmd/main/main.go
build-frontend:
	rm -rf frontend/dist
	cd frontend && npm run build

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
	@echo "Generate file movies_v1.ts"
	@protoc -I=. --plugin="protoc-gen-ts=frontend/node_modules/.bin/protoc-gen-ts" --ts_opt=target=web --ts_out=frontend/src \
		protos/movies_v1/movies_v1.proto
	@echo "Adding headers to the generated file..."
	@sed -i '1i/* eslint-disable */\n// @ts-nocheck' frontend/src/protos/movies_v1/movies_v1.ts

generate-frontend-interactions:
	@echo "Generate file interactions_v1.ts"
	@protoc -I=. --plugin="protoc-gen-ts=frontend/node_modules/.bin/protoc-gen-ts" --ts_opt=target=web --ts_out=frontend/src \
        protos/interactions_v1/interactions_v1.proto
	@echo "Adding headers to the generated file..."
	@sed -i '1i/* eslint-disable */\n// @ts-nocheck' frontend/src/protos/interactions_v1/interactions_v1.ts

generate-frontend-comments:
	@echo "Generate file comments_v1.ts"
	@protoc -I=. --plugin="protoc-gen-ts=frontend/node_modules/.bin/protoc-gen-ts" --ts_opt=target=web --ts_out=frontend/src \
        protos/comments_v1/comments_v1.proto
	@echo "Adding headers to the generated file..."
	@sed -i '1i/* eslint-disable */\n// @ts-nocheck' frontend/src/protos/comments_v1/comments_v1.ts
	@sed -i '1i/* eslint-disable */\n// @ts-nocheck' frontend/src/google/protobuf/timestamp.ts