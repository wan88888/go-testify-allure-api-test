# Go + Testify + Allure API 自动化测试框架

基于 Golang、Testify 和 Allure 构建的接口自动化测试框架，使用 Platzi Fake Store API 作为测试数据源。

## 🚀 项目特性

- **现代化技术栈**: 使用 Go 1.21+ 、Testify 测试框架和 Allure 报告
- **完整的 API 覆盖**: 涵盖商品、分类、用户、购物车等所有 API 端点
- **丰富的测试场景**: 包括正常流程、异常处理、数据验证、性能测试等
- **详细的测试报告**: 集成 Allure 生成美观的测试报告
- **灵活的配置管理**: 支持 YAML 配置文件，易于环境切换
- **模块化设计**: 清晰的项目结构，易于维护和扩展
- **并发测试支持**: 支持并行执行测试用例，提高执行效率

## 📁 项目结构

```
go-testify-allure-api-test/
├── client/                 # API 客户端
│   └── api_client.go      # HTTP 客户端封装
├── config/                # 配置管理
│   └── config.go          # 配置文件解析
├── models/                # 数据模型
│   └── models.go          # API 响应结构体
├── tests/                 # 测试用例
│   ├── products_test.go   # 商品相关测试
│   ├── categories_test.go # 分类相关测试
│   ├── users_test.go      # 用户相关测试
│   └── carts_test.go      # 购物车相关测试
├── utils/                 # 工具函数
│   └── test_utils.go      # 测试辅助工具
├── config.yaml            # 配置文件
├── go.mod                 # Go 模块文件
├── main_test.go           # 主测试入口
├── Makefile              # 构建脚本
└── README.md             # 项目文档
```

## 🛠️ 环境要求

- **Go**: 1.21 或更高版本
- **Allure**: 2.13+ (用于生成测试报告)
- **网络**: 能够访问 https://fakestoreapi.com

## 📦 安装和设置

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-testify-allure-api-test
```

### 2. 安装依赖

```bash
# 安装 Go 依赖
make deps

# 或者手动安装
go mod tidy
```

### 3. 安装 Allure (可选)

```bash
# 使用 Makefile 自动安装
make install

# 或者手动安装
# macOS (使用 Homebrew)
brew install allure

# 使用 npm
npm install -g allure-commandline
```

## 🚀 快速开始

### 运行所有测试

```bash
# 运行所有测试
make test

# 运行测试并生成报告
make test-and-report

# 完整流程：测试 + 报告 + 启动服务器
make full
```

### 运行特定测试套件

```bash
# 只运行商品相关测试
make test-products

# 只运行分类相关测试
make test-categories

# 只运行用户相关测试
make test-users

# 只运行购物车相关测试
make test-carts
```

### 生成和查看报告

```bash
# 生成 Allure 报告
make report

# 启动报告服务器
make serve
```

## 📊 测试覆盖范围

### 商品 API (Products)
- ✅ 获取所有商品列表
- ✅ 根据 ID 获取单个商品
- ✅ 限制数量获取商品
- ✅ 排序获取商品
- ✅ 创建新商品
- ✅ 更新商品信息
- ✅ 删除商品
- ✅ 异常情况处理

### 分类 API (Categories)
- ✅ 获取所有商品分类
- ✅ 根据分类获取商品
- ✅ 分类数据一致性验证
- ✅ 无效分类处理

### 用户 API (Users)
- ✅ 获取所有用户列表
- ✅ 根据 ID 获取用户信息
- ✅ 用户登录认证
- ✅ 用户数据格式验证
- ✅ 无效凭据处理

### 购物车 API (Carts)
- ✅ 获取所有购物车
- ✅ 根据 ID 获取购物车
- ✅ 购物车数据一致性验证
- ✅ 性能测试

## 🔧 配置说明

项目使用 `config.yaml` 文件进行配置：

```yaml
api:
  base_url: "https://fakestoreapi.com"  # API 基础 URL
  timeout: 30                           # 请求超时时间（秒）
  retry_count: 3                        # 重试次数

allure:
  results_dir: "allure-results"         # Allure 结果目录
  report_dir: "allure-report"           # Allure 报告目录

test:
  parallel: true                        # 是否并行执行测试
  verbose: true                         # 是否显示详细输出
  cleanup: true                         # 是否自动清理

logging:
  level: "info"                         # 日志级别
  format: "json"                        # 日志格式
  output: "console"                     # 日志输出
```

## 📈 测试报告

框架集成了 Allure 测试报告，提供以下功能：

- **测试概览**: 测试执行统计、成功率、耗时等
- **测试详情**: 每个测试用例的详细执行步骤
- **请求/响应**: 完整的 HTTP 请求和响应信息
- **错误分析**: 失败用例的详细错误信息
- **趋势分析**: 历史测试结果对比
- **分类统计**: 按功能模块分类的测试结果

## 🧪 测试用例设计

### 测试分层
1. **API 层测试**: 验证 HTTP 接口的正确性
2. **数据层测试**: 验证响应数据的完整性和格式
3. **业务层测试**: 验证业务逻辑的正确性
4. **异常层测试**: 验证错误处理的健壮性

### 测试类型
- **功能测试**: 验证 API 功能是否正常
- **边界测试**: 验证边界条件处理
- **异常测试**: 验证异常情况处理
- **性能测试**: 验证响应时间和并发能力
- **数据一致性测试**: 验证不同 API 间数据的一致性

## 🔍 高级功能

### 并行测试
```bash
# 使用并行模式运行测试
make test-parallel
```

### 环境检查
```bash
# 检查环境配置
make check
```

### 清理资源
```bash
# 清理测试结果和报告
make clean
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📝 开发规范

### 代码规范
- 遵循 Go 官方代码规范
- 使用有意义的变量和函数名
- 添加必要的注释和文档
- 保持代码简洁和可读性

### 测试规范
- 每个测试用例应该独立且可重复执行
- 使用描述性的测试名称
- 添加详细的测试步骤和断言
- 包含正常和异常场景的测试

### 提交规范
- 使用清晰的提交信息
- 一次提交只包含一个功能或修复
- 提交前运行所有测试确保通过

## 🐛 故障排除

### 常见问题

1. **网络连接问题**
   ```bash
   # 检查网络连接
   curl -I https://fakestoreapi.com
   ```

2. **依赖安装问题**
   ```bash
   # 清理并重新安装依赖
   go clean -modcache
   go mod download
   ```

3. **Allure 报告问题**
   ```bash
   # 检查 Allure 安装
   allure --version
   
   # 重新安装 Allure
   make install
   ```

4. **测试超时问题**
   - 检查网络连接
   - 增加配置文件中的超时时间
   - 使用 `-timeout` 参数增加测试超时时间

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Platzi Fake Store API](https://fakestoreapi.com/) - 提供测试数据源
- [Testify](https://github.com/stretchr/testify) - Go 测试框架
- [Allure](https://docs.qameta.io/allure/) - 测试报告框架
- [Resty](https://github.com/go-resty/resty) - HTTP 客户端库
- [Viper](https://github.com/spf13/viper) - 配置管理库

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 创建 Issue
- 发送 Pull Request
- 邮件联系: [your-email@example.com]

---

**Happy Testing! 🎉**