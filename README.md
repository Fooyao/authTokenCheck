# Twitter 账号状态批量检查工具

这是一个高效的 Twitter 账号状态批量检查工具，使用 Go 语言开发，支持并发处理、智能解析和增量检查。

## 🚀 主要特性

- **🔍 智能解析**: 自动识别 auth_token 和 ct0，支持任意顺序排列
- **⚡ 高效并发**: 10个协程并发检查，大幅提升处理速度
- **📊 状态分类**: 自动按账号状态分类输出到不同文件
- **🔄 增量检查**: 自动跳过已检查的账号，支持断点续传
- **🛡️ 容错处理**: 完善的错误处理机制，程序稳定可靠
- **📝 格式灵活**: 支持多种输入格式，容错能力强

## 📦 编译和安装

### 方法一：下载预编译版本（推荐）

从 [GitHub Releases](https://github.com/your-username/twitter-checker/releases) 下载最新版本：

```bash
# 1. 访问 Release 页面下载对应平台的压缩包
# Windows: twitter-checker-windows-amd64.zip
# Linux: twitter-checker-linux-amd64.tar.gz  
# macOS: twitter-checker-macos-amd64.tar.gz

# 2. 解压文件
# Windows
unzip twitter-checker-windows-amd64.zip

# Linux/macOS
tar -xzf twitter-checker-linux-amd64.tar.gz

# 3. 运行程序
./twitter-checker accounts.txt
```

### 方法二：从源码编译

```bash
# 1. 确保已安装 Go 1.19 或更高版本
go version

# 2. 克隆或下载项目文件
# 确保包含：main.go, go.mod, go.sum

# 3. 安装依赖
go mod tidy

# 4. 直接运行
go run main.go <账号文件>
```

### 方法三：编译为可执行文件

```bash
# 编译为当前平台可执行文件
go build -o twitter-checker main.go

# 运行编译后的程序
./twitter-checker <账号文件>
```

### 方法四：交叉编译

```bash
# 编译为 Windows 64位
GOOS=windows GOARCH=amd64 go build -o twitter-checker.exe main.go

# 编译为 Linux 64位
GOOS=linux GOARCH=amd64 go build -o twitter-checker-linux main.go

# 编译为 macOS 64位
GOOS=darwin GOARCH=amd64 go build -o twitter-checker-mac main.go
```

### 方法五：本地批量编译

```bash
# 使用提供的脚本编译所有平台
chmod +x scripts/release.sh
./scripts/release.sh 1.0.0
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
1234567890abcdef1234567890abcdef12345678----1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef

# 颠倒顺序：ct0----auth_token
abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890----abcdef1234567890abcdef1234567890abcdef12

# 只有auth_token（ct0可选）
fedcba0987654321fedcba0987654321fedcba09

# 包含其他数据
username----1111222233334444555566667777888899990000----email@domain.com----other_data

# 复杂混合格式
user_info----aaaa1111bbbb2222cccc3333dddd4444eeee5555----extra_data----ccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbbccccddddeeeeffffaaaabbbb----more_info
```

### 智能解析规则

1. **自动识别**: 根据字符串长度和格式自动识别字段类型
2. **顺序无关**: auth_token 和 ct0 可以任意顺序排列
3. **容错处理**: 自动跳过格式错误的行，继续处理其他行
4. **重复字段**: 如果同一行有多个相同类型字段，使用第一个
5. **必需字段**: auth_token 是必需的，ct0 是可选的
6. **注释支持**: 以 `#` 开头的行会被忽略

## 🔧 使用方法

### 基本用法

```bash
# 使用 go run
go run main.go accounts.txt

# 使用编译后的程序
./twitter-checker accounts.txt
```

### 帮助信息

```bash
go run main.go
```

输出：
```
使用方法: go run main.go <accounts_file.txt>
文件格式: 每行按----分割，至少包含auth_token
---------------------------------------------------
good_accounts.txt， 正常账号
bad_token_accounts.txt， 无效token
suspended_accounts.txt， 封禁账号
locked_accounts.txt， 锁定账号
error_accounts.txt， 错误账号，请重命名后重新查询
```

## 📊 程序输出示例

```
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