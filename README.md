# 🐦 Twitter账号状态批量检查工具

[![Build Status](https://github.com/Fooyao/authTokenCheck/workflows/Build%20and%20Release/badge.svg)](https://github.com/Fooyao/authTokenCheck/actions)
[![Release](https://img.shields.io/github/release/Fooyao/authTokenCheck.svg)](https://github.com/Fooyao/authTokenCheck/releases)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

一个高性能的Twitter账号状态批量检查工具，使用Go语言开发，支持智能解析、并发处理和增量检查。

**🚀 所有版本通过 GitHub Actions 自动编译，支持多平台，下载即用**

## 🚀 主要特性

- **🔄 智能解析**: 自动识别 auth_token 和 ct0，支持多种格式排列
- **⚡ 高性能并发**: 10个 goroutine 并发查询，大幅提升检查速度
- **📈 增量检查**: 自动跳过已检查账号，支持断点续传
- **🏷️ 状态分类**: 自动将不同状态账号分类保存到不同文件
- **📝 格式灵活**: 支持多种输入格式，容错能力强
- **🏗️ 自动构建**: 通过 GitHub Actions 自动编译多平台版本

## 📦 快速开始

### 方法一：下载预编译版本（推荐）

**所有版本均由 GitHub Actions 自动编译，无需本地环境配置**

1. 前往 [Releases 页面](../../releases)
2. 下载对应平台的二进制文件：
   - **Windows**: `twitter-checker-windows-amd64.exe` 或 `twitter-checker-windows-386.exe`
   - **Linux**: `twitter-checker-linux-amd64` 或 `twitter-checker-linux-arm64`
   - **macOS**: `twitter-checker-macos-amd64` 或 `twitter-checker-macos-arm64`
3. **Linux/macOS 用户**：添加执行权限
   ```bash
   chmod +x twitter-checker-*
   ```
4. 直接运行：
   ```bash
   # Linux/macOS
   ./twitter-checker-linux-amd64 accounts.txt
   
   # Windows
   twitter-checker-windows-amd64.exe accounts.txt
   ```

### 方法二：本地编译（开发者）

如果需要自定义编译或参与开发：

```bash
# 1. 确保已安装 Go 1.21 或更高版本
go version

# 2. 克隆项目
git clone <repository-url>
cd authTokenCheck

# 3. 安装依赖
go mod tidy

# 4. 直接运行
go run main.go accounts.txt

# 5. 或编译为可执行文件
go build -o twitter-checker main.go
./twitter-checker accounts.txt
```

### 方法三：多平台批量编译

使用项目提供的编译脚本：

```bash
# 编译所有平台版本
chmod +x scripts/release.sh
./scripts/release.sh 1.0.0

# 编译结果在 dist/ 目录
ls -la dist/
```

## 📋 账号状态说明

程序会将账号分类为以下状态：

| 状态 | 文件名 | 说明 |
|------|--------|------|
| **GOOD** | `good_accounts.txt` | 账号状态正常，可正常使用 |
| **BAD_TOKEN** | `bad_token_accounts.txt` | 认证令牌无效，需要重新获取 |
| **SUSPENDED** | `suspended_accounts.txt` | 账号已被平台暂停 |
| **LOCKED** | `locked_accounts.txt` | 账号已被锁定，需要验证 |
| **ERROR** | `error_accounts.txt` | 网络错误或其他异常，建议重新检查 |

## 📄 输入文件格式

### 基本要求

- 文本文件，每行一个账号信息
- 使用 `----` 作为字段分隔符
- **auth_token**: 40位小写十六进制字符串（必需）
- **ct0**: 160位小写十六进制字符串（可选）
- 支持注释行（以 `#` 开头）

### 支持的格式示例

```text
# 标准格式：auth_token----ct0
1234567890abcdef1234567890abcdef12345678----1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef

# 颠倒顺序：ct0----auth_token
abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890----abcdef1234567890abcdef1234567890abcdef12

# 只有auth_token（ct0可选）
fedcba0987654321fedcba0987654321fedcba09

# 包含其他数据
username----1111222233334444555566667777888899990000----email@domain.com----other_data

# 复杂混合格式
user_info----aaaa1111bbbb2222cccc3333dddd4444eeee5555----extra_data----ccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbb----more_info
```

### 智能解析规则

1. **自动识别**: 根据字符串长度和格式自动识别字段类型
2. **顺序无关**: auth_token 和 ct0 可以任意顺序排列
3. **容错处理**: 自动跳过格式错误的行，继续处理其他行
4. **重复字段**: 如果同一行有多个相同类型字段，使用第一个
5. **必需字段**: auth_token 是必需的，ct0 是可选的
6. **注释支持**: 以 `#` 开头的行会被忽略

## 🔧 使用方法

### 快速使用（预编译版本）

```bash
# 1. 从 GitHub Releases 下载对应平台的文件
# 2. 准备账号文件 accounts.txt
# 3. 直接运行

# Linux/macOS 示例
chmod +x twitter-checker-linux-amd64
./twitter-checker-linux-amd64 accounts.txt

# Windows 示例  
twitter-checker-windows-amd64.exe accounts.txt
```

### 命令行选项

```bash
# 显示版本信息
./twitter-checker --version

# 显示帮助信息
./twitter-checker
```

### 示例输出
```
Twitter 账号状态检查工具 v1.0.0
正在从文件 accounts.txt 读取账号信息...
发现 150 个已检查的账号，将跳过重复检查
跳过了 50 个已经检查过的账号
成功读取 100 个账号，开始并发检查...

Worker 1 完成检查: 1234567890... -> GOOD
Worker 2 完成检查: abcdef1234... -> SUSPENDED
Worker 3 完成检查: fedcba0987... -> BAD_TOKEN
...

检查完成，共处理 100 个账号

状态统计:
GOOD: 45 个账号
SUSPENDED: 30 个账号
BAD_TOKEN: 15 个账号
LOCKED: 8 个账号
ERROR: 2 个账号

正在将结果写入文件...
已将 45 个 GOOD 状态账号写入 good_accounts.txt
已将 30 个 SUSPENDED 状态账号写入 suspended_accounts.txt
已将 15 个 BAD_TOKEN 状态账号写入 bad_token_accounts.txt
已将 8 个 LOCKED 状态账号写入 locked_accounts.txt
已将 2 个 ERROR 状态账号写入 error_accounts.txt

任务完成！
```

## 🎯 核心功能

### 1. 增量检查
- 自动读取现有结果文件
- 跳过已经检查过的账号
- 支持断点续传，提高效率
- 适合大批量处理

### 2. 并发处理
- 10个 goroutine 并发检查
- 显著提升处理速度
- 合理控制并发数，避免过载

### 3. 智能解析
- 自动识别 auth_token 和 ct0
- 支持任意顺序和格式
- 强大的容错能力

### 4. 结果分类
- 按状态自动分类输出
- 便于后续处理和分析
- 详细的统计信息

## ⚙️ 技术规格

### 系统要求
- **Go版本**: 1.19 或更高
- **内存**: 建议 512MB 以上
- **网络**: 稳定的互联网连接

### 依赖包
```go
require (
	github.com/go-resty/resty/v2 v2.11.0  // HTTP客户端
	github.com/tidwall/gjson v1.17.0      // JSON解析
)
```

### 性能参数
- **并发数**: 10个 goroutine
- **超时时间**: 60秒/请求
- **处理速度**: 约100-500账号/分钟（取决于网络）

## 📝 注意事项

### 使用建议
1. **网络环境**: 确保网络连接稳定，避免频繁超时
2. **文件格式**: 仔细检查输入文件格式，程序会跳过错误行
3. **频率控制**: 大量账号检查时注意API调用频率限制
4. **结果备份**: 重要结果建议及时备份

### 错误处理
1. **格式错误**: 程序会显示警告并跳过错误行
2. **网络错误**: 自动归类到 ERROR 状态，建议重新检查
3. **文件权限**: 确保程序有读写文件的权限

### 文件管理
- 结果文件会覆盖同名的现有文件
- ERROR 状态的账号建议重命名文件后重新检查
- 可以合并多次检查的结果文件

## 🔄 工作流程

1. **初始化**: 读取现有结果文件，建立已检查账号索引
2. **解析**: 智能解析输入文件，提取有效账号信息
3. **过滤**: 跳过已检查的账号，准备待检查列表
4. **并发**: 启动10个worker并发检查账号状态
5. **分类**: 根据检查结果分类存储到不同文件
6. **统计**: 输出详细的处理统计信息

## 🛠️ 故障排除

### 常见问题

**Q: 提示 "no required module provides package"**
```bash
A: 运行 go mod tidy 安装依赖包
```

**Q: 大量账号显示 ERROR 状态**
```bash
A: 检查网络连接，考虑降低并发数或增加延时
```

**Q: 程序运行很慢**
```bash
A: 检查网络速度，确认是否有代理或防火墙影响
```

**Q: 结果文件被覆盖**
```bash
A: 程序会覆盖同名文件，重要结果请及时备份
```

---

## 📞 技术支持

如果遇到问题或需要功能定制，请提供：
1. 错误信息截图
2. 输入文件示例（脱敏）
3. 运行环境信息
4. 具体需求描述 

## 🚀 自动化构建说明

### GitHub Actions 工作流

本项目采用 GitHub Actions 实现自动化构建和发布，具有以下特点：

- **多平台支持**: 自动编译 Windows、Linux、macOS 三个平台的版本
- **架构完整**: 支持 amd64、arm64、386 等多种架构
- **版本管理**: 通过 Git tag 触发自动发布（格式：v1.0.0）
- **质量保证**: 自动运行测试、依赖验证等检查步骤
- **即时可用**: 编译完成后立即可下载使用，无需额外配置

### 支持的平台

| 平台 | 架构 | 文件名 |
|------|------|--------|
| Windows | amd64 | `twitter-checker-windows-amd64.exe` |
| Windows | 386 | `twitter-checker-windows-386.exe` |
| Linux | amd64 | `twitter-checker-linux-amd64` |
| Linux | arm64 | `twitter-checker-linux-arm64` |
| macOS | amd64 | `twitter-checker-macos-amd64` |
| macOS | arm64 | `twitter-checker-macos-arm64` |

### 发布流程

1. 开发者推送代码并创建版本标签：`git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions 自动触发构建流程
3. 自动编译所有平台版本
4. 运行测试和质量检查
5. 生成 SHA256 校验和
6. 创建 GitHub Release 并上传所有文件
7. 用户可立即下载使用

### 优势

- **零配置**: 用户无需安装 Go 环境或处理依赖
- **多平台**: 一次构建，支持所有主流操作系统
- **自动化**: 完全自动化的构建和发布流程
- **可追溯**: 每个版本都有完整的构建日志和校验和
- **持续集成**: 代码质量通过自动化测试保障

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

**📧 问题反馈**: 如遇到问题请提交 [Issue](../../issues)  
**⭐ 支持项目**: 如果觉得有用请给个星标  
**🔧 贡献代码**: 欢迎提交 Pull Request 
