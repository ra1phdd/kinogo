.PHONY: generate-frontend generate-frontend-movies generate-frontend-interactions

generate-frontend: generate-frontend-movies generate-frontend-interactions

generate-frontend-movies:
	protoc --plugin="protoc-gen-grpc-web=frontend/node_modules/.bin/protoc-gen-grpc-web" --grpc-web_out=import_style=typescript,mode=grpcwebtext:frontend/src \
		--plugin="protoc-gen-js=frontend/node_modules/.bin/protoc-gen-js" --js_out=import_style=commonjs:frontend/src \
		protos/movies_v1/movies_v1.proto
	protoc -I=. --plugin="protoc-gen-js=frontend/node_modules/.bin/protoc-gen-js" --js_out=frontend/src \
		protos/movies_v1/movies_v1.proto

generate-frontend-interactions:
	protoc --plugin="protoc-gen-grpc-web=frontend/node_modules/.bin/protoc-gen-grpc-web" --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:frontend/src \
		   protos/interactions_v1/interactions_v1.proto