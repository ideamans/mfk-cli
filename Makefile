.PHONY: build install uninstall clean test help

# デフォルトターゲット
.DEFAULT_GOAL := help

# バイナリ名
BINARY_NAME := mfk
INSTALL_PATH := /usr/local/bin

# ビルド
build: ## バイナリをビルド
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/main.go
	@echo "Build complete: $(BINARY_NAME)"

# インストール
install: build ## バイナリを /usr/local/bin にインストール（要sudo）
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo chmod 755 $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete: $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo ""
	@echo "You can now run '$(BINARY_NAME)' from anywhere."

# アンインストール
uninstall: ## インストールされたバイナリを削除（要sudo）
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_PATH)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstallation complete."

# クリーンアップ
clean: ## ビルド成果物を削除
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f tmp/*.pdf
	@echo "Clean complete."

# テスト実行
test: ## Goテストを実行
	@echo "Running tests..."
	@go test -v ./...

# 依存関係の整理
tidy: ## go mod tidyを実行
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies tidied."

# バージョン確認
version: ## インストールされているmfkのバージョンを確認
	@if [ -f "$(INSTALL_PATH)/$(BINARY_NAME)" ]; then \
		echo "Installed: $(INSTALL_PATH)/$(BINARY_NAME)"; \
		ls -lh $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		echo "$(BINARY_NAME) is not installed."; \
		echo "Run 'make install' to install."; \
	fi

# ローカルテスト実行
test-local: build ## ビルドしてローカルでテスト実行
	@echo "Running local test..."
	@./$(BINARY_NAME) --help
	@echo ""
	@echo "Test complete. Run './$(BINARY_NAME) download-invoice --help' for more info."

# ヘルプ
help: ## このヘルプメッセージを表示
	@echo "使用可能なターゲット:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "使用例:"
	@echo "  make build      # バイナリをビルド"
	@echo "  make install    # /usr/local/bin にインストール"
	@echo "  make uninstall  # インストールを削除"
	@echo "  make clean      # ビルド成果物を削除"
