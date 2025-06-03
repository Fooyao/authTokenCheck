# 发布指南

## 自动发布 (推荐)

### 1. 更新版本信息
在发布新版本前，请先更新以下文件：

```bash
# 更新 CHANGELOG.md
# 将 [Unreleased] 部分的内容移动到新版本下
# 添加发布日期

# 提交更改
git add CHANGELOG.md
git commit -m "准备发布 v1.0.1"
git push origin main
```

### 2. 创建并推送标签
```bash
# 创建标签 (遵循语义化版本)
git tag v1.0.1

# 推送标签到远程仓库
git push origin v1.0.1
```

### 3. 自动构建和发布
推送标签后，GitHub Actions 将自动：
- ✅ 为多个平台编译二进制文件
- ✅ 创建压缩包
- ✅ 生成校验和文件
- ✅ 创建 GitHub Release
- ✅ 上传所有文件到 Release

## 手动发布 (备用)

### 1. 本地编译
```bash
# 使用发布脚本
chmod +x scripts/release.sh
./scripts/release.sh 1.0.1

# 或者手动编译
go build -ldflags="-s -w -X main.version=v1.0.1" -o twitter-checker main.go
```

### 2. 手动创建 Release
1. 前往 GitHub 仓库页面
2. 点击 "Releases" -> "Create a new release"
3. 选择或创建标签
4. 上传编译好的文件
5. 填写发布说明

## 版本号规则

遵循 [语义化版本](https://semver.org/lang/zh-CN/) 规则：

- **主版本号** (X.0.0): 不兼容的 API 修改
- **次版本号** (0.X.0): 向下兼容的功能性新增
- **修订号** (0.0.X): 向下兼容的问题修正

### 示例
- `v1.0.0` - 首个稳定版本
- `v1.0.1` - 修复 bug
- `v1.1.0` - 新增功能
- `v2.0.0` - 重大更新，可能不兼容

## 支持的平台

自动构建支持以下平台：
- Windows (amd64, 386)
- Linux (amd64, arm64)
- macOS (amd64, arm64)

## 发布检查清单

发布前请确认：

- [ ] 代码已通过测试
- [ ] CHANGELOG.md 已更新
- [ ] 版本号符合语义化版本规则
- [ ] 所有更改已提交到 main 分支
- [ ] GitHub Actions 配置正确

## 故障排除

### GitHub Actions 失败
如果自动构建失败：

1. 检查 Actions 页面的错误日志
2. 确认 go.mod 和 go.sum 文件正确
3. 验证代码能在本地编译
4. 检查网络连接和依赖下载

### 手动构建失败
```bash
# 检查 Go 环境
go version
go env

# 验证依赖
go mod verify
go mod tidy

# 测试编译
go build main.go
```

## 发布后

1. 验证 Release 页面文件完整
2. 测试下载的二进制文件
3. 更新相关文档
4. 通知用户新版本发布 