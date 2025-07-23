.PHONY: help test test-verbose test-products test-categories test-users test-carts clean install report serve deps

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  make deps        - 安装项目依赖"
	@echo "  make test        - 运行所有测试"
	@echo "  make test-verbose - 运行所有测试（详细输出）"
	@echo "  make test-products - 只运行商品相关测试"
	@echo "  make test-categories - 只运行分类相关测试"
	@echo "  make test-users  - 只运行用户相关测试"
	@echo "  make test-carts  - 只运行购物车相关测试"
	@echo "  make report      - 生成Allure测试报告"
	@echo "  make serve       - 启动Allure报告服务器"
	@echo "  make clean       - 清理测试结果和报告"
	@echo "  make install     - 安装Allure命令行工具"

# 安装项目依赖
deps:
	@echo "正在安装Go模块依赖..."
	go mod tidy
	go mod download

# 运行所有测试
test: clean
	@echo "正在运行所有API测试..."
	mkdir -p allure-results
	go test -v ./tests/... -timeout 30m

# 运行所有测试（详细输出）
test-verbose: clean
	@echo "正在运行所有API测试（详细模式）..."
	mkdir -p allure-results
	go test -v -count=1 ./tests/... -timeout 30m

# 运行商品相关测试
test-products: clean
	@echo "正在运行商品相关测试..."
	mkdir -p allure-results
	go test -v ./tests/ -run TestProductsTestSuite -timeout 15m

# 运行分类相关测试
test-categories: clean
	@echo "正在运行分类相关测试..."
	mkdir -p allure-results
	go test -v ./tests/ -run TestCategoriesTestSuite -timeout 15m

# 运行用户相关测试
test-users: clean
	@echo "正在运行用户相关测试..."
	mkdir -p allure-results
	go test -v ./tests/ -run TestUsersTestSuite -timeout 15m

# 运行购物车相关测试
test-carts: clean
	@echo "正在运行购物车相关测试..."
	mkdir -p allure-results
	go test -v ./tests/ -run TestCartsTestSuite -timeout 15m

# 生成Allure测试报告
report:
	@echo "正在生成Allure测试报告..."
	@if command -v allure >/dev/null 2>&1; then \
		allure generate allure-results --clean -o allure-report; \
		echo "报告已生成到 allure-report 目录"; \
	else \
		echo "错误: 未找到allure命令，请先运行 'make install' 安装Allure"; \
		exit 1; \
	fi

# 启动Allure报告服务器
serve:
	@echo "正在启动Allure报告服务器..."
	@if command -v allure >/dev/null 2>&1; then \
		if [ -d "allure-report" ]; then \
			allure open allure-report; \
		else \
			echo "错误: 报告目录不存在，请先运行 'make report' 生成报告"; \
			exit 1; \
		fi; \
	else \
		echo "错误: 未找到allure命令，请先运行 'make install' 安装Allure"; \
		exit 1; \
	fi

# 清理测试结果和报告
clean:
	@echo "正在清理测试结果和报告..."
	rm -rf allure-results allure-report

# 安装Allure命令行工具
install:
	@echo "正在检查并安装Allure命令行工具..."
	@if command -v brew >/dev/null 2>&1; then \
		echo "使用Homebrew安装Allure..."; \
		brew install allure; \
	elif command -v npm >/dev/null 2>&1; then \
		echo "使用npm安装Allure..."; \
		npm install -g allure-commandline; \
	else \
		echo "错误: 未找到brew或npm，请手动安装Allure"; \
		echo "访问 https://docs.qameta.io/allure/#_installing_a_commandline 获取安装说明"; \
		exit 1; \
	fi

# 运行测试并生成报告的组合命令
test-and-report: test report
	@echo "测试完成，报告已生成"

# 完整的测试流程（测试 + 报告 + 服务器）
full: test report serve

# 检查环境
check:
	@echo "正在检查环境..."
	@echo "Go版本:"
	@go version
	@echo "\n项目依赖状态:"
	@go mod verify
	@echo "\nAllure状态:"
	@if command -v allure >/dev/null 2>&1; then \
		allure --version; \
	else \
		echo "Allure未安装，运行 'make install' 进行安装"; \
	fi

# 并行运行测试（更快）
test-parallel: clean
	@echo "正在并行运行API测试..."
	mkdir -p allure-results
	go test -v -parallel 4 ./tests/... -timeout 30m