FROM python:3.9
COPY ./requirements.txt /opt/competition/requirements.txt
RUN pip install -r /opt/competition/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

ADD ./limit.py /opt/competition/limit.py
ADD ./judge.py /opt/competition/judge.py
ADD ./unzip.sh /opt/competition/unzip.sh
ADD ./Config.json /opt/competition/Config.json
RUN chmod 777 /opt/competition/limit.py
RUN chmod 777 /opt/competition/judge.py
RUN chmod 777 /opt/competition/unzip.sh
# ADD ./其它需要脚本文件 /opt/competition/其他需要脚本文件
ENTRYPOINT ["/opt/competition/limit.py"]