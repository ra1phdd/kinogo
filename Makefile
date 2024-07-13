PHONY: test cover build build-backend build-frontend generate \
	generate-backend generate-backend-movies generate-backend-interactions generate-backend-comments generate-backend-metrics \
	generate-frontend generate-frontend-movies generate-frontend-interactions generate-frontend-comments generate-frontend-metrics

# Переменные
BACKEND_DIR := backend
FRONTEND_DIR := frontend
PROTO_DIR := protos
BUILD_DIR := $(BACKEND_DIR)/build
PROTOC_GEN_TS := $(FRONTEND_DIR)/node_modules/.bin/protoc-gen-ts

# Юнит тесты и покрытие кода на backend
test:
	cd $(BACKEND_DIR) && go test -race -count 1 ./...

cover:
	cd $(BACKEND_DIR) && go test -short -race -count 1 -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out
	rm $(BACKEND_DIR)/coverage.out

# Сборка
build: build-backend build-frontend

build-backend:
	mkdir -p $(BUILD_DIR)
	rm -rf $(BUILD_DIR)/*
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/kinogo_amd64 ./cmd/main/main.go
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/kinogo_arm64 ./cmd/main/main.go

build-frontend:
	rm -rf $(FRONTEND_DIR)/dist
	cd $(FRONTEND_DIR) && npm run build

# Генерация proto-файлов
generate: generate-backend generate-frontend

generate-backend: generate-backend-movies generate-backend-interactions generate-backend-comments generate-backend-metrics

generate-frontend: generate-frontend-movies generate-frontend-interactions generate-frontend-comments generate-frontend-metrics

# Шаблон для генерации proto-файлов на backend
define generate_backend_template
	@echo "Генерация файлов $(1).pb.go и $(1)_grpc.pb.go"
	@rm -rf $(BACKEND_DIR)/pkg/$(1)
	@mkdir -p $(BACKEND_DIR)/pkg/$(1)
	@protoc --go_out=$(BACKEND_DIR)/pkg/$(1) --go-grpc_out=$(BACKEND_DIR)/pkg/$(1) $(PROTO_DIR)/$(1)/$(1).proto
	@mv $(BACKEND_DIR)/pkg/$(1)/$(1)/* $(BACKEND_DIR)/pkg/$(1)
	@rm -rf $(BACKEND_DIR)/pkg/$(1)/$(1)/
endef

generate-backend-movies:
	$(call generate_backend_template,movies_v1)

generate-backend-interactions:
	$(call generate_backend_template,interactions_v1)

generate-backend-comments:
	$(call generate_backend_template,comments_v1)

generate-backend-metrics:
	$(call generate_backend_template,metrics_v1)

# Шаблон для генерации proto-файлов на frontend
define generate_frontend_template
	@echo "Генерация файла $(1).ts"
	@protoc -I=. --plugin="protoc-gen-ts=$(PROTOC_GEN_TS)" --ts_opt=target=web --ts_out=$(FRONTEND_DIR)/src $(PROTO_DIR)/$(1)/$(1).proto
	@echo "Отключение правил ESLint и проверки типов TypeScript в файле $(1).ts..."
	@sed -i '1i/* eslint-disable */\n// @ts-nocheck' $(FRONTEND_DIR)/src/protos/$(1)/$(1).ts
endef

generate-frontend-movies:
	$(call generate_frontend_template,movies_v1)

generate-frontend-interactions:
	$(call generate_frontend_template,interactions_v1)

generate-frontend-comments:
	$(call generate_frontend_template,comments_v1)
	@sed -i '1i/* eslint-disable */\n// @ts-nocheck' $(FRONTEND_DIR)/src/google/protobuf/timestamp.ts

generate-frontend-metrics:
	$(call generate_frontend_template,metrics_v1)
