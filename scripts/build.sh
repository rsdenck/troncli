#!/bin/bash
set -e

APP_NAME="troncli"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR="bin"

echo "Building $APP_NAME version $VERSION..."

platforms=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm/v7"
)

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    GOARM=${platform_split[2]}
    
    output_name=$APP_NAME'_'$GOOS'_'$GOARCH
    if [ "$GOARCH" = "arm" ]; then
        output_name=$output_name'v'$GOARM
    fi
    
    env_str="GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0"
    if [ "$GOARCH" = "arm" ]; then
        env_str="$env_str GOARM=$GOARM"
    fi

    echo "Building for $GOOS/$GOARCH..."
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    export CGO_ENABLED=0
    if [ "$GOARCH" = "arm" ]; then
        export GOARM=$GOARM
    fi

    go build -ldflags "-s -w -X main.version=$VERSION" -trimpath -o "$BUILD_DIR/$output_name" ./cmd/troncli
    
    if [ $? -eq 0 ]; then
        echo "✅ Built $output_name"
    else
        echo "❌ Failed to build for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Build complete! Artifacts in $BUILD_DIR/"
ls -lh $BUILD_DIR/
