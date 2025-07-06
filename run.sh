#!/bin/bash

# 检查是否提供了目录路径参数
if [ $# -eq 0 ]; then
    # 如果没有提供目录路径参数，使用脚本所在的目录
    script_dir="$(dirname "$0")"
else
    # 如果提供了目录路径参数，使用参数中指定的目录
    script_dir="$1"
fi

# 切换到指定目录
cd "$script_dir" || exit

# 执行可执行程序
./gohumanloop-feishu
