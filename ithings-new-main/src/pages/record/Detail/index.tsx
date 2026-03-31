import { PageContainer } from '@ant-design/pro-components';
import { Card, Descriptions, message } from 'antd';
import React, { useEffect, useState } from 'react';
import { useParams } from '@umijs/max';
import type { RecordInfo } from '@/services/competitionPlatform';
import { getRecord } from '@/services/competitionPlatform';

const RecordDetailPage: React.FC = () => {
  const params = useParams<{ cid: string; rid: string }>();
  const [data, setData] = useState<RecordInfo>();

  useEffect(() => {
    const loadData = async () => {
      if (!params.cid || !params.rid) return;
      try {
        const { data } = await getRecord(params.cid, params.rid);
        setData(data);
      } catch (error) {
        message.error('获取提交详情失败');
      }
    };
    loadData();
  }, [params.cid, params.rid]);

  return (
    <PageContainer>
      <Card>
        <Descriptions column={1} bordered>
          <Descriptions.Item label="比赛">{data?.competition?.title}</Descriptions.Item>
          <Descriptions.Item label="用户">{data?.user?.username}</Descriptions.Item>
          <Descriptions.Item label="状态">{data?.status}</Descriptions.Item>
          <Descriptions.Item label="得分">{data?.score}</Descriptions.Item>
          <Descriptions.Item label="运行时间">{data?.run_time}</Descriptions.Item>
          <Descriptions.Item label="错误信息">{data?.errors || '-'}</Descriptions.Item>
          <Descriptions.Item label="提交时间">{data?.create_time}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{data?.latest_time}</Descriptions.Item>
          <Descriptions.Item label="完成时间">{data?.finish_time}</Descriptions.Item>
        </Descriptions>
      </Card>
    </PageContainer>
  );
};

export default RecordDetailPage;
