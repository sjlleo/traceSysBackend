#!/bin/bash

# 获取用户输入的端口号
read -p "请输入要查询的端口号：" port

# 使用lsof命令查询占用该端口的进程
result=$(sudo lsof -i :$port)

# 如果查询结果为空，说明没有进程占用该端口
if [ -z "$result" ]; then
  echo "没有进程占用该端口"
else
  # 如果有进程占用该端口，则输出进程信息
  echo "占用该端口的进程信息："
  echo "$result"
  
  # 提示是否杀掉该进程
  read -p "是否要杀掉占用该端口的进程？（y/n）" answer
  if [ "$answer" == "y" ]; then
    # 使用awk命令获取进程ID，并使用kill命令杀掉该进程
    pid=$(echo "$result" | awk 'NR==2 {print $2}')
    kill -9 $pid
    echo "已杀掉占用该端口的进程"
  fi
fi

