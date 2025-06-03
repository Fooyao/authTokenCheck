#!/bin/bash

# Twitter 账号检查工具 - 本地编译脚本
# 使用方法: ./scripts/release.sh [version]

set -e

VERSION=${1:-"dev"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "🚀 开始编译 Twitter 账号检查工具"
echo "📦 版本: $VERSION"
echo "⏰ 编译时间: $BUILD_TIME"
echo "📝 Git提交: $GIT_COMMIT"
echo ""

# 清理旧的构建文件
echo "🧹 清理旧文件..."
rm -rf dist/
mkdir -p dist/

# 定义编译目标
declare -a targets=(
    "windows/amd64"
    "windows/386" 
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

# LDFLAGS 设置
LDFLAGS="-s -w -X main.version=$VERSION"

echo "🔨 开始编译..."
echo ""

for target in "${targets[@]}"
do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    
    output_name="twitter-checker-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "  📦 编译 ${GOOS}/${GOARCH}..."
    
    GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
        -ldflags="$LDFLAGS" \
        -o "dist/$output_name" \
        main.go
    
    # 创建压缩包
    cd dist
    if [ $GOOS = "windows" ]; then
        zip -q "${output_name%.exe}.zip" "$output_name"
        echo "    ✅ 创建压缩包: ${output_name%.exe}.zip"
    else
        tar -czf "${output_name}.tar.gz" "$output_name"
        echo "    ✅ 创建压缩包: ${output_name}.tar.gz"
    fi
    cd ..
done

echo ""
echo "📊 生成校验和..."
cd dist
sha256sum *.zip *.tar.gz > checksums.txt 2>/dev/null || true
cd ..

echo ""
echo "✅ 编译完成！"
echo ""
echo "📂 输出文件:"
ls -la dist/
echo ""
echo "📋 文件校验和:"
if [ -f "dist/checksums.txt" ]; then
    cat dist/checksums.txt
else
    echo "  (校验和文件未生成)"
fi

echo ""
echo "🎉 所有平台编译完成！"
echo "📁 文件位置: ./dist/"

if [ "$VERSION" != "dev" ]; then
    echo ""
    echo "💡 提示: 可以使用以下命令创建 Git 标签并推送:"
    echo "   git tag v$VERSION"
    echo "   git push origin v$VERSION"
    echo ""
    echo "📤 这将触发 GitHub Actions 自动发布 Release"
fi 