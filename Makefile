.PHONY: help test test-verbose test-products test-categories test-users test-carts clean deps check test-parallel

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  make deps        - 安装项目依赖"
	@echo "  make test        - 运行所有测试"
	@echo "  make test-verbose - 运行所有测试（详细输出）"
	@echo "  make test-parallel - 并行运行所有测试"
	@echo "  make test-products - 只运行商品相关测试"
	@echo "  make test-categories - 只运行分类相关测试"
	@echo "  make test-users  - 只运行用户相关测试"
	@echo "  make test-carts  - 只运行购物车相关测试"
	@echo "  make clean       - 清理测试结果"
	@echo "  make check       - 检查环境配置"

# 安装项目依赖
deps:
	@echo "正在安装Go模块依赖..."
	go mod tidy
	go mod download

# 运行所有测试
test: clean
	@echo "正在运行所有API测试..."
	go test -v ./tests/... -timeout 30m

# 运行所有测试（详细输出）
test-verbose: clean
	@echo "正在运行所有API测试（详细模式）..."
	go test -v -count=1 ./tests/... -timeout 30m

# 运行商品相关测试
test-products: clean
	@echo "正在运行商品相关测试..."
	go test -v ./tests/ -run "Test.*Product" -timeout 15m

# 运行分类相关测试
test-categories: clean
	@echo "正在运行分类相关测试..."
	go test -v ./tests/ -run "Test.*Categor" -timeout 15m

# 运行用户相关测试
test-users: clean
	@echo "正在运行用户相关测试..."
	go test -v ./tests/ -run "Test.*User" -timeout 15m

# 运行购物车相关测试
test-carts: clean
	@echo "正在运行购物车相关测试..."
	go test -v ./tests/ -run "Test.*Cart" -timeout 15m

# 生成Allure报告
report:
	@echo "生成Allure报告..."
	allure generate tests/allure-results --clean -o allure-report
	@echo "Allure报告已生成到 allure-report 目录"

# 启动Allure服务器
serve:
	@echo "启动Allure服务器..."
	allure serve tests/allure-results

# 安装Allure命令行工具
install:
	@echo "安装Allure命令行工具..."
	npm install -g allure-commandline
	@echo "Allure已安装完成"

# 运行测试并生成报告
test-and-report: test report
	@echo "测试完成并已生成报告"

# 完整流程：测试 -> 报告 -> 服务
full: test report serve

# 清理测试结果
clean:
	@echo "正在清理测试结果..."
	rm -rf tests/allure-results allure-report

# 检查环境
check:
	@echo "正在检查环境..."
	@echo "Go版本:"
	@go version
	@echo "\n项目依赖状态:"
	@go mod verify

# 并行运行测试（更快）
test-parallel: clean
	@echo "正在并行运行API测试..."
	go test -v -parallel 4 ./tests/... -timeout 30m