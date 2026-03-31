## 算力天梯出题模板

[算力天梯用户手册](https://scbnampb7x.feishu.cn/wiki/WzKZw2NtQiUHPGk9oSIcr063nGh?from=from_copylink)

出题者需要准备如下文件：

1. **Config.json**：参数配置文件，应包含输入文件路径、输出文件路径、标准答案路径（如有）、压缩包存放路径、解压后用户文件路径、解压脚本路径、打分脚本路径、日志文件路径和通知接口地址。

2. **judge.py**：打分脚本，根据题目要求输出用户得分。

3. **limit.py**：核心处理脚本，流程为：读入配置——解压——文件规范性检查——执行用户程序——超时检测——执行打分程序——记录日志，通知后端。该文件需要根据Config.json的设置进行修改，主要需要修改output_files规范性检查部分。

   > 注意！！！若是在Windows环境编辑再上传，会遇到镜像创建不成功或不能运行容器错误。请保证在Linux环境编辑，或者直接在本文件上编辑，或者从本地上传后添加 `\#! /usr/bin/env python3` 到文件第一行，并在vim中输入 `:set ff=unix` 后保存并退出。

4. **requirement.txt**：镜像中所有python文件运行所需环境。

5. **unzip.sh**：解压脚本，可自定义逻辑。模板中简单的进入压缩包目录，解压并删除压缩包。

6. **Dockerfile**：创建镜像脚本，流程为：引入python环境——安装依赖——添加所需脚本文件——授予待执行文件权限——配置容器启动后命令。

创建镜像命令为（需先进入Dockerfile所在路径）： `docker build -t <镜像名> .`

运行容器命令为： `docker run -v <输入文件路径>:/input/:ro -v <用户文件路径>:/user-data/ -v <输出文件路径>:/output/ -v <日志文件路径>:/log/ -v <标准答案路径>:/answer/:ro <镜像名> <record主键> <超时时间>`