#! /usr/bin/env python3
# coding:utf-8
import json
import pandas as pd

with open("/opt/competition/Config.json", "r", encoding="utf-8") as f:
    content = json.load(f)

output_path = content["output_path"]  # 用户答案文件路径
answer_path = content["answer_path"]  # 标准答案文件路径


def main():
    # 判分逻辑
    score = 0
    print(score)  # 子进程的控制台输出会以字符串的形式传给父进程的result，如有换行会替换为'\n'


if __name__ == '__main__':
    main()
