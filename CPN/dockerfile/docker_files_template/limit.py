#! /usr/bin/env python3
# coding:utf-8
import json
import logging
import os
import subprocess
import sys
import time
import requests


def loginit(logname, logpath):
    logger = logging.getLogger(logname)
    logger.setLevel(logging.DEBUG)
    # 创建处理器
    loghandler = logging.FileHandler(logpath)
    loghandler.setLevel(logging.DEBUG)
    # 创建格式器(时间-级别-程序名-函数名-行号: 信息)
    logformatter = logging.Formatter('%(asctime)s-%(levelname)s-%(filename)s-%(lineno)d: %(message)s')
    # 将格式器添加到处理器中
    loghandler.setFormatter(logformatter)
    # 将处理器添加到记录器中
    logger.addHandler(loghandler)
    return logger


def notice(url, record_id, competition_id, score=-1, error="", run_time=-1, status=""):
    try:
        full_url = f"{url}/competition/{competition_id}/record/{record_id}"
        data = {
            "score": score,
            "error": error,
            "runTime": run_time,
            "status": status,
        }
        requests.put(full_url, data=data)
    except Exception as e:
        logger.error("通知后端失败: %s", e)


# 读入配置文件
with open("/opt/competition/Config.json", "r", encoding="utf-8") as f:
    content = json.load(f)

input_path = content["input_path"]  # 用户输入文件路径
output_path = content["output_path"]  # 用户输出文件路径
answer_path = content["answer_path"]  # 标准答案路径
zipfile_path = content["zipfile_path"]  # 待解压文件路径
zip_path = content["zip_path"]  # 解压脚本路径
judge_path = content["judge_path"]  # 打分脚本路径
codefile_path = content["codefile_path"]  # 解压后实际执行程序路径
log_path = content["log_path"]  # 日志路径
url = content["url"]  # 通知接口的url

fileId = sys.argv[1]
timeout = int(sys.argv[2])
competition_id = sys.argv[3]

logger = loginit('limit-log', log_path)


def function(fileId, timeout, competition_id):
    # timeout 赛题要求运行时间上限
    # fileId record主键

    # code.zip规范性检查
    if not os.path.exists(zipfile_path + "/code.zip"):
        logger.error("Record:%s 待解压文件不存在", fileId)
        notice(url, fileId, competition_id, status="解压失败", error="待解压文件不存在")
        return

    # 解压文件
    zip_command = 'sh ' + zip_path + ' ' + zipfile_path
    unzip_status = subprocess.call(zip_command, shell=True)
    if unzip_status != 0:
        logger.error("Record:%s 解压失败", fileId)
        notice(url, fileId, competition_id, status="解压失败", error="解压失败")
        return
    logger.info("Record:%s 解压成功", fileId)

    # code/main.py规范性检查
    if not os.path.exists(codefile_path):
        logger.error("Record:%s 待执行文件不存在", fileId)
        notice(url, fileId, competition_id, status="格式错误", error="待执行文件不存在")
        return

    codePopen = None
    start = time.time()
    try:
        codePopen = subprocess.Popen(['python', codefile_path], stderr=subprocess.PIPE)
        _, error = codePopen.communicate(timeout=timeout)
        if error:
            logger.error("Record:%s 用户提交运行出错 %s", fileId, error.decode('utf-8'))
            notice(url, fileId, competition_id, status="运行失败", error=error.decode('utf-8'))
            return
    except subprocess.TimeoutExpired:
        codePopen.kill()
        logger.error("Record:%s 运行超时", fileId)
        notice(url, fileId, competition_id, status="运行超时", error="运行超时")
        return
    except Exception as e:
        logger.error("Record:%s 系统运行出错 %s", fileId, e)
        notice(url, fileId, competition_id, status="运行失败", error=str(e))
        return

    runtime = time.time() - start
    logger.info("Record:%s 用户提交运行成功", fileId)

    # output_files规范性检查
    for i in range(1, 4):
        file_name = '/output' + str(i) + '.csv'
        output_file = output_path + file_name
        if not os.path.exists(output_file):
            logger.error("Record:%s 答案文件%s不存在", fileId, file_name)
            notice(url, fileId, competition_id, status="格式错误", error="答案文件%s不存在" % file_name)
            return

    # 执行打分脚本
    try:
        judgePopen = subprocess.Popen(['python', judge_path], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        result, error = judgePopen.communicate()
        if error:
            logger.error("Record:%s 输出格式不规范或错误", fileId)
            notice(url, fileId, competition_id, status="格式错误", error=error.decode('utf-8'))
            return
        else:
            score_str = result.decode('utf-8').strip()
            logger.info("Record:%s 运行成功, score=%s, 运行时间为%s", fileId, score_str, round(runtime, 4))
            try:
                score_val = int(float(score_str))
            except ValueError:
                score_val = -1
            notice(url, fileId, competition_id, score=score_val, run_time=int(round(runtime, 4)), status="运行成功")
            return
    except Exception as e:
        logger.error("Record:%s %s", fileId, e)
        notice(url, fileId, competition_id, status="运行失败", error=str(e))
        return


function(fileId, timeout, competition_id)



