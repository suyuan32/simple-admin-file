# Custom configuration | 独立配置
# Service name | 项目名称
SERVICE=Fms
# Service name in specific style | 项目经过style格式化的名称
SERVICE_STYLE=fms
# Service name in lowercase | 项目名称全小写格式
SERVICE_LOWER=fms
# Service name in snake format | 项目名称下划线格式
SERVICE_SNAKE=fms
# Service name in snake format | 项目名称短杠格式
SERVICE_DASH=fms

# The project version, if you don't use git, you should set it manually | 项目版本，如果不使用git请手动设置
VERSION=v1.3.10

# The project file name style | 项目文件命名风格
PROJECT_STYLE=go_zero

# Whether to use i18n | 是否启用 i18n
PROJECT_I18N=true

# The suffix after build or compile | 构建后缀
PROJECT_BUILD_SUFFIX=api

# Swagger type, support yml,json | Swagger 文件类型，支持yml,json
SWAGGER_TYPE=json

# Ent enabled features | Ent 启用的官方特性
ENT_FEATURE=sql/execquery,intercept

# Auto generate API data for initialization | 自动生成 API 初始化数据
AUTO_API_INIT_DATA=true

# The arch of the build | 构建的架构
GOARCH=amd64

# ---- You may not need to modify the codes below | 下面的代码大概率不需要更改 ----

GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
LDFLAGS := -s -w

.PHONY: test
test: # Run test for the project | 运行项目测试
	go test -v --cover ./internal/..

.PHONY: fmt
fmt: # Format the codes | 格式化代码
	$(GOFMT) -w $(GOFILES)

.PHONY: lint
lint: # Run go linter | 运行代码错误分析
	golangci-lint run -D staticcheck

.PHONY: tools
tools: # Install the necessary tools | 安装必要的工具
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;
	$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest


.PHONY: docker
docker: # Build the docker image | 构建 docker 镜像
	docker build -f Dockerfile -t ${DOCKER_USERNAME}/$(SERVICE_DASH)-$(PROJECT_BUILD_SUFFIX)-docker:${VERSION} .
	@echo "Build docker successfully"

.PHONY: publish-docker
publish-docker: # Publish docker image | 发布 docker 镜像
	echo "${DOCKER_PASSWORD}" | docker login --username ${DOCKER_USERNAME} --password-stdin https://${REPO}
	docker push ${DOCKER_USERNAME}/$(SERVICE_DASH)-$(PROJECT_BUILD_SUFFIX):${VERSION}
	@echo "Publish docker successfully"

.PHONY: gen-swagger
gen-swagger: # Generate swagger file | 生成 swagger 文件
	swagger generate spec --output=./$(SERVICE_STYLE).$(SWAGGER_TYPE) --scan-models
	@echo "Generate swagger successfully"

.PHONY: serve-swagger
serve-swagger: # Run the swagger server | 运行 swagger 服务
	lsof -i:36666 | awk 'NR!=1 {print $2}' | xargs killall -9 || true
	swagger serve -F=swagger --port 36666 $(SERVICE_STYLE).$(SWAGGER_TYPE)
	@echo "Serve swagger-ui successfully"

.PHONY: gen-api
gen-api: # Generate API files | 生成 API 的代码
	goctls api go --api ./desc/all.api --dir ./ --trans_err=true --style=$(PROJECT_STYLE)
	swagger generate spec --output=./$(SERVICE_STYLE).$(SWAGGER_TYPE) --scan-models
	@echo "Generate API codes successfully"

.PHONY: gen-ent
gen-ent: # Generate Ent codes | 生成 Ent 的代码
	go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./ent/template/*.tmpl" ./ent/schema --feature $(ENT_FEATURE)
	@echo "Generate Ent codes successfully"

.PHONY: gen-api-ent-logic
gen-api-ent-logic: # Generate CRUD logic from Ent, need to set model and group | 根据 Ent 生成 CRUD 代码，需要设置 model 和 group
	goctls api ent --schema=./ent/schema --style=$(PROJECT_STYLE) --api_service_name=$(SERVICE) --output=./ --model=$(model) --group=$(group) --i18n=$(PROJECT_I18N) --overwrite=true --api_data=$(AUTO_API_INIT_DATA)
	@echo "Generate CRUD codes from Ent successfully"

.PHONY: build-win
build-win: # Build project for Windows | 构建Windows下的可执行文件
	env CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_$(PROJECT_BUILD_SUFFIX).exe $(SERVICE_STYLE).go
	@echo "Build project for Windows successfully"

.PHONY: build-mac
build-mac: # Build project for MacOS | 构建MacOS下的可执行文件
	env CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_$(PROJECT_BUILD_SUFFIX) $(SERVICE_STYLE).go
	@echo "Build project for MacOS successfully"

.PHONY: build-linux
build-linux: # Build project for Linux | 构建Linux下的可执行文件
	env CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_$(PROJECT_BUILD_SUFFIX) $(SERVICE_STYLE).go
	@echo "Build project for Linux successfully"

.PHONY: help
help: # Show help | 显示帮助
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done
