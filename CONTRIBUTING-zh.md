# 🤝 为 Kubernetes Resources Recommend 做贡献

感谢您对为 Kubernetes Resources Recommend 做贡献的兴趣！我们欢迎来自每个人的贡献，无论经验水平如何。

## 🌍 语言支持

- **[English](CONTRIBUTING.md)**
- **中文** (当前)

## 📋 目录

- [行为准则](#行为准则)
- [入门指南](#入门指南)
- [开发环境设置](#开发环境设置)
- [如何贡献](#如何贡献)
- [Pull Request 流程](#pull-request-流程)
- [Issue 指南](#issue-指南)
- [代码风格](#代码风格)
- [测试](#测试)
- [文档](#文档)

## 📜 行为准则

本项目及其参与者都受我们行为准则的约束。通过参与，您需要遵守此准则。请向[项目维护者](mailto:your-email@example.com)报告不当行为。

### 我们的标准

- **相互尊重**: 以尊重和善意对待每个人
- **包容性**: 欢迎所有背景和身份的人
- **协作性**: 建设性地合作
- **专业性**: 在所有互动中保持专业行为

## 🚀 入门指南

### 前置要求

- Go 1.23.9 或更高版本
- Git
- 基本的 Kubernetes 和 Prometheus 知识
- 熟悉 Go 编程语言

### 首次设置

1. **在 GitHub 上 Fork 仓库**
2. **本地克隆您的 fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/kubernetes-resources-recommend.git
   cd kubernetes-resources-recommend
   ```

3. **添加原始仓库为上游**:
   ```bash
   git remote add upstream https://github.com/luozijian1990/kubernetes-resources-recommend.git
   ```

4. **安装依赖**:
   ```bash
   go mod download
   ```

5. **构建项目**:
   ```bash
   make build
   # 或者
   go build -o bin/kubernetes-resources-recommend cmd/kubernetes-resources-recommend/main.go
   ```

## 🛠️ 开发环境设置

### 项目结构

```
kubernetes-resources-recommend/
├── cmd/                     # 主应用程序
├── internal/                # 私有应用程序代码
│   ├── exporter/           # 导出功能
│   ├── prometheus/         # Prometheus 客户端
│   ├── recommender/        # 核心推荐逻辑
│   └── types/              # 类型定义
├── pkg/                    # 公共库代码
└── docs/                   # 文档
```

### 可用命令

```bash
make help          # 查看所有可用命令
make build         # 构建项目
make test          # 运行测试
make test-coverage # 运行测试并生成覆盖率
make fmt           # 格式化代码
make lint          # 代码检查
make clean         # 清理构建产物
```

### 环境设置

为本地开发创建 `.env` 文件:
```bash
PROMETHEUS_URL=https://your-prometheus.example.com
CHECK_NAMESPACE=default
LIMITS=1.5
```

## 🔄 如何贡献

### 贡献类型

我们欢迎多种类型的贡献：

- 🐛 **Bug 修复**
- ✨ **新功能**
- 📝 **文档改进**
- 🧪 **测试**
- 🌍 **翻译**
- 🎨 **UI/UX 改进**
- 📊 **性能优化**

### 开始之前

1. **检查现有 issues** 以避免重复工作
2. **为新功能或重大更改创建 issue**
3. **如需要，与维护者讨论您的方法**

### 开发工作流

1. **创建功能分支**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **进行更改**:
   - 编写清洁、可读的代码
   - 遵循现有代码风格
   - 为新功能添加测试
   - 根据需要更新文档

3. **测试您的更改**:
   ```bash
   make test
   make test-coverage
   make lint
   ```

4. **提交您的更改**:
   ```bash
   git add .
   git commit -m "feat: 添加您的功能描述"
   ```

5. **推送到您的 fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **在 GitHub 上创建 Pull Request**

## 📥 Pull Request 流程

### 提交前

- [ ] 代码遵循项目风格指南
- [ ] 测试在本地通过
- [ ] 文档已更新
- [ ] 提交消息遵循约定
- [ ] 与主分支没有合并冲突

### PR 标题格式

使用约定式提交格式：
- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更改
- `refactor:` 代码重构
- `test:` 添加测试
- `chore:` 维护任务

### PR 描述模板

```markdown
## 描述
简要描述更改

## 更改类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 重构
- [ ] 性能改进

## 测试
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 手动测试完成

## 检查清单
- [ ] 代码遵循风格指南
- [ ] 完成自我审查
- [ ] 文档已更新
- [ ] 无破坏性更改（或已记录）
```

### 审查流程

1. **自动检查** 必须通过
2. **维护者代码审查**
3. **在不同环境中测试**
4. **项目维护者批准**
5. **合并** 到主分支

## 🐛 Issue 指南

### Bug 报告

报告 bug 时，请包含：

- **清晰的标题** 和描述
- **重现问题的步骤**
- **预期 vs 实际行为**
- **环境详情**:
  - Go 版本
  - Kubernetes 版本
  - Prometheus 版本
  - 操作系统
- **日志和错误消息**
- **截图**（如适用）

### 功能请求

对于功能请求，请提供：

- **功能的清晰描述**
- **用例** 和动机
- **建议的实现**（如有）
- **考虑的替代方案**
- **额外的上下文**

### Issue 标签

- `bug`: 某些功能不正常
- `enhancement`: 新功能请求
- `documentation`: 文档改进
- `good first issue`: 适合新手
- `help wanted`: 需要额外关注
- `priority/high`: 高优先级问题
- `priority/low`: 低优先级问题

## 🎨 代码风格

### Go 风格指南

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 进行格式化
- 使用有意义的变量和函数名
- 为导出的函数添加注释
- 保持函数小而专注

### 命名约定

- **包**: 小写，单个单词
- **函数**: camelCase，如果导出则以大写开头
- **变量**: camelCase
- **常量**: UPPER_CASE 或 camelCase
- **文件**: 小写带下划线

### 代码组织

- 分组相关功能
- 将关注点分离到不同的包中
- 使用接口进行抽象
- 适当处理错误
- 编写自文档化的代码

## 🧪 测试

### 测试策略

- **单元测试**: 测试单个函数
- **集成测试**: 测试组件交互
- **端到端测试**: 测试完整工作流

### 编写测试

```go
func TestRecommenderFunction(t *testing.T) {
    // 安排
    input := setupTestData()
    
    // 行动
    result := functionUnderTest(input)
    
    // 断言
    assert.Equal(t, expected, result)
}
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率
make test-coverage

# 运行特定测试
go test -v ./internal/recommender -run TestSpecificFunction
```

### 测试覆盖率

- 目标至少 80% 的测试覆盖率
- 专注于关键业务逻辑
- 测试边界情况和错误条件

## 📝 文档

### 文档类型

- **代码注释**: 解释复杂逻辑
- **README 文件**: 项目概述和设置
- **API 文档**: 函数和类型文档
- **用户指南**: 操作指南和教程

### 文档标准

- 使用清晰、简洁的语言
- 在有帮助的地方提供示例
- 保持文档最新
- 使用适当的 markdown 格式

### 更新文档

进行更改时：
- 更新相关的 README 部分
- 添加/更新代码注释
- 更新 API 文档
- 如需要，创建/更新用户指南

## 🏷️ 版本控制

我们使用 [语义化版本](https://semver.org/):
- **主版本**: 不兼容的 API 更改
- **次版本**: 新功能（向后兼容）
- **补丁版本**: Bug 修复（向后兼容）

## 🎯 路线图

当前优先级：
- [ ] CPU 资源推荐
- [ ] GPU 资源推荐
- [ ] HPA 集成
- [ ] Web UI 界面
- [ ] API 端点
- [ ] Helm chart

## 🙋‍♀️ 获取帮助

如果您需要帮助：

1. **首先检查文档**
2. **搜索 GitHub 上的现有 issues**
3. **创建包含详细信息的新 issue**
4. **参与 issues 和 PRs 中的讨论**
5. **如需要，直接联系维护者**

## 🏆 认可

贡献者将在以下地方得到认可：
- **README.md** 贡献者部分
- **重要贡献的发布说明**
- **项目更新中的特别提及**

感谢您为 Kubernetes Resources Recommend 做贡献！🎉

---

## 📧 联系方式

- **维护者**: [luozijian1990](https://github.com/luozijian1990)
- **Issues**: [GitHub Issues](https://github.com/luozijian1990/kubernetes-resources-recommend/issues)
- **讨论**: [GitHub Discussions](https://github.com/luozijian1990/kubernetes-resources-recommend/discussions)
