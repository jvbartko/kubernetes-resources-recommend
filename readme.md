# 🎯 Kubernetes Resources Recommend

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

> 🚀 基于 Prometheus 监控数据的 Kubernetes Deployment 内存资源智能推荐系统

## ✨ 功能特性

🔍 **智能分析** - 基于过去7天的内存使用数据进行深度分析  
📊 **科学算法** - 使用P90值和指数衰减权重算法计算精准推荐值  
🏢 **多租户支持** - 支持多命名空间独立分析  
⚡ **高性能** - 并发处理，大幅提升分析效率  
📑 **友好报告** - 导出美观的Excel格式分析报告  

## 🎯 推荐范围说明

> ⚠️ **重要说明**：本工具当前**仅提供内存(Memory)资源推荐**

### 💾 为什么只推荐内存资源？

**内存(Memory) - 不可压缩资源**  
- 🚨 **严格限制** - 当容器内存使用超过limit时，Kubernetes会立即终止(OOMKilled)Pod
- 📈 **准确预测** - 内存使用模式相对稳定，历史数据具有很好的预测价值
- ⚖️ **关键平衡** - 设置过低导致频繁OOM，设置过高浪费资源
- 🎯 **精确控制** - 需要基于真实使用数据进行精确推荐

**CPU - 可压缩资源(暂不推荐)**  
- 🔄 **可压缩性** - CPU资源不足时，容器会被限流而不会被杀死
- 📊 **复杂模式** - CPU使用受业务峰值、并发量等多因素影响，预测复杂
- 🌊 **动态调整** - Kubernetes HPA可以根据CPU使用率自动扩缩容
- ⏰ **时机考虑** - 计划在后续版本中加入CPU推荐功能

### 🔮 未来规划

- **v2.0** - 增加CPU资源推荐


## 📁 项目结构

```
kubernetes-resources-recommend/
├── 📂 cmd/kubernetes-resources-recommend/   # 🚪 主程序入口
│   └── main.go
├── 📂 internal/                            # 🔒 内部包，不对外暴露
│   ├── 📂 exporter/                        # 📊 导出功能
│   │   └── excel.go
│   ├── 📂 prometheus/                      # 📈 Prometheus客户端
│   │   ├── client.go
│   │   └── metrics.go
│   ├── 📂 recommender/                     # 🧠 推荐算法核心
│   │   └── recommender.go
│   └── 📂 types/                          # 📋 类型定义
│       ├── prometheus.go
│       └── recommendation.go
├── 📂 pkg/                                # 📦 公共包
│   └── 📂 config/                         # ⚙️ 配置管理
│       ├── config.go
│       └── errors.go
├── 🔧 Makefile                           # 🛠️ 构建脚本
├── 📄 go.mod & go.sum                    # 📚 依赖管理
└── 📖 README.md                          # 📝 项目文档
```

## 🚀 快速开始

### 📦 安装

```bash
# 克隆项目
git clone https://github.com/your-username/kubernetes-resources-recommend.git
cd kubernetes-resources-recommend

# 安装依赖
go mod download
```

### 🔨 构建

```bash
# 使用 Makefile (推荐)
make build

# 或手动构建
go build -o bin/kubernetes-resources-recommend cmd/kubernetes-resources-recommend/main.go

# 跨平台构建
make build-all
```

### 💡 使用示例

```bash
# 基本使用
./bin/kubernetes-resources-recommend \
  -prometheusUrl="https://prometheus.your-domain.com" \
  -checkNamespace="production" \
  -limits=1.5

# 使用 Makefile 运行
make run ARGS="-prometheusUrl=https://prometheus.example.com -checkNamespace=staging"
```

### 📋 参数说明

| 参数 | 类型 | 默认值 | 说明 |
|-----|------|--------|------|
| `-prometheusUrl` | string | `https://prometheus.example.com` | 🌐 Prometheus服务器地址 |
| `-checkNamespace` | string | `default` | 🏷️ 要分析的Kubernetes命名空间 |
| `-limits` | float64 | `1.5` | 📏 内存限制相对于请求的倍数 |

## 🧮 内存推荐算法

我们的内存资源智能推荐算法基于以下步骤：

```mermaid
graph TD
    A[🔍 收集7天历史数据] --> B[⏰ 按小时粒度分析]
    B --> C[📊 计算每日P90值]
    C --> D[⚖️ 应用指数衰减权重]
    D --> E[🎯 生成最终推荐值]
    
    F[💾 container_memory_rss] --> A
    G[🔢 weight = 0.5^(day+1)] --> D
    E --> H[📄 Excel报告输出]
```

**算法详解：**

1. **📈 数据收集** - 分析过去7天的内存使用数据
2. **⏱️ 精细粒度** - 以小时为单位收集 `container_memory_rss` 指标
3. **📊 统计分析** - 计算每日P90值（排除异常峰值）
4. **⚖️ 智能权重** - 应用指数衰减权重：`weight = 0.5^(day+1)`
5. **🎯 最终计算** - 加权求和得到推荐的内存请求值
6. **🛡️ 安全边界** - 内存限制 = 内存请求 × limits倍数

## 📊 输出报告

程序会生成名为 `{namespace}-resource-recommend.xlsx` 的Excel文件：

> 📝 **注意**：当前版本仅包含**内存资源推荐**，CPU资源推荐将在后续版本提供

### 📋 报告字段说明

| 列名 | 说明 | 示例 | 备注 |
|------|------|------|------|
| 🏷️ Namespace | 命名空间 | `production` | Kubernetes命名空间 |
| 🚀 Deployment | 部署名称 | `web-server` | Deployment资源名 |
| 📦 Container | 容器名称 | `app-container` | 容器名称 |
| 💾 Memory Request (MB) | 推荐的内存请求值 | `512` | 基于7天P90算法 |
| 🛡️ Memory Limit (MB) | 推荐的内存限制值 | `768` | Request × limits倍数 |

### 📈 Excel输出示例

生成的Excel文件结构如下：

```
Resource Recommendations.xlsx
┌─────────────┬──────────────┬──────────────┬────────────────────┬───────────────────┐
│  Namespace  │  Deployment  │  Container   │ Memory Request(MB) │ Memory Limit(MB)  │
├─────────────┼──────────────┼──────────────┼────────────────────┼───────────────────┤
│ production  │ web-server   │ nginx        │        256         │        384        │
│ production  │ web-server   │ app          │        512         │        768        │
│ production  │ api-gateway  │ gateway      │        128         │        192        │
│ production  │ redis        │ redis        │        64          │        96         │
│ staging     │ test-app     │ app          │        128         │        192        │
│ staging     │ test-app     │ sidecar      │        32          │        48         │
└─────────────┴──────────────┴──────────────┴────────────────────┴───────────────────┘
```

### 📊 报告特色

- **🎨 美观格式** - 带有标题样式和自动列宽调整
- **📈 数据完整** - 包含所有分析的Deployment和Container
- **🔍 易于筛选** - 支持Excel的筛选和排序功能
- **📝 清晰标识** - 明确的列名和单位标识

### 🎯 实际使用场景

**生成报告示例：**
```bash
$ ./bin/kubernetes-resources-recommend -checkNamespace=production -limits=1.5
2024/01/15 10:30:00 Starting Kubernetes resource recommendation
2024/01/15 10:30:01 All required metrics are available
2024/01/15 10:30:01 Generating memory recommendations...
2024/01/15 10:30:05 Processed namespace: production, deployment: web-server, container: nginx
2024/01/15 10:30:08 Processed namespace: production, deployment: web-server, container: app
2024/01/15 10:30:10 Processed namespace: production, deployment: api-gateway, container: gateway
2024/01/15 10:30:12 Generated 6 recommendations
2024/01/15 10:30:12 Recommendations exported to production-resource-recommend.xlsx
2024/01/15 10:30:12 Process completed in 12.5s
```

**生成的文件：**
- 📄 `production-resource-recommend.xlsx` - 包含详细的资源推荐数据
- 📊 格式化的Excel表格，可直接用于资源配置更新

## 🔧 开发指南

### 🛠️ 可用命令

```bash
make help          # 📖 查看所有可用命令
make build         # 🔨 构建项目
make test          # 🧪 运行测试
make test-coverage # 📊 运行测试并生成覆盖率报告
make fmt           # 🎨 格式化代码
make lint          # 🔍 代码检查
make clean         # 🧹 清理构建产物
```

### 📋 系统要求

- **Go版本**: 1.23.9+ 
- **Prometheus**: 需要安装 kube-state-metrics
- **Kubernetes**: 集群环境（用于获取监控指标）

### 📈 必需的 Prometheus 指标

确保您的 Prometheus 实例收集以下指标：

| 指标名称 | 说明 | 来源 |
|----------|------|------|
| `container_memory_rss` | 容器内存使用量 | cAdvisor |
| `kube_pod_owner` | Pod所有者信息 | kube-state-metrics |
| `kube_replicaset_owner` | ReplicaSet所有者信息 | kube-state-metrics |
| `kube_deployment_created` | Deployment创建时间 | kube-state-metrics |
| `kube_deployment_spec_replicas` | Deployment副本数规格 | kube-state-metrics |

## 🤝 贡献指南

我们欢迎所有形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

### 🐛 问题报告

如果您发现了bug或有功能建议，请 [创建一个Issue](https://github.com/your-username/kubernetes-resources-recommend/issues)。

## 📄 许可证

本项目基于 MIT 许可证开源。详情请查看 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Prometheus](https://prometheus.io/) - 强大的监控系统
- [excelize](https://github.com/qax-os/excelize) - Excel文件处理库
- [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) - Kubernetes指标导出器

---

<div align="center">

**如果这个项目对您有帮助，请给我们一个 ⭐ Star！**

Made with ❤️ by [Your Name](https://github.com/your-username)

</div>