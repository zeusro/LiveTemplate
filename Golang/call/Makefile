.PHONY: build run proto clean

# 生成 proto 文件
proto:
	@echo "Generating proto files..."
	@mkdir -p pb
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		asr.proto
	@echo "Proto files generated in pb/ directory"

# 构建 HTTP 服务器
build-http:
	go build -tags '!rpc' -o bin/asr-server asrserver.go

# 构建 gRPC 服务器
build-rpc:
	go build -tags rpc -o bin/asr-rpc-server asrrpcserver.go

# 构建所有
build: proto build-http build-rpc

# 运行 HTTP 服务器
run-http: build-http
	./bin/asr-server -f etc/asr.yaml

# 运行 gRPC 服务器
run-rpc: build-rpc
	./bin/asr-rpc-server -f etc/asr.yaml

# 清理
clean:
	rm -rf bin/ pb/
