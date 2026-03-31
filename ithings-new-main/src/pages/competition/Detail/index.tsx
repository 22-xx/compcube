import { PageContainer } from '@ant-design/pro-components';
import { Button, Card, Descriptions, Space, message } from 'antd';
import React, { useEffect, useState } from 'react';
import { history, useParams } from '@umijs/max';
import type { CompetitionInfo } from '@/services/competitionPlatform';
import { getCompetition } from '@/services/competitionPlatform';

const CompetitionDetailPage: React.FC = () => {
  const params = useParams<{ id: string }>();
  const [data, setData] = useState<CompetitionInfo>();

  useEffect(() => {
    const loadData = async () => {
      if (!params.id) return;
      try {
        const { data } = await getCompetition(params.id);
        setData(data);
      } catch (error) {
        message.error('获取比赛详情失败');
      }
    };
    loadData();
  }, [params.id]);

  return (
    <PageContainer
      extra={[
        <Space key="actions">
          <Button onClick={() => history.push(`/record/rank/${params.id}`)}>查看排行</Button>
          <Button type="primary" onClick={() => history.push(`/record/submit/${params.id}`)}>
            提交代码
          </Button>
        </Space>,
      ]}
    >
      <Card>
        <Descriptions column={1} bordered>
          <Descriptions.Item label="标题">{data?.title}</Descriptions.Item>
          <Descriptions.Item label="作者">{data?.author?.username || '-'}</Descriptions.Item>
          <Descriptions.Item label="状态">{data?.status}</Descriptions.Item>
          <Descriptions.Item label="时间限制">{data?.time_limit}</Descriptions.Item>
          <Descriptions.Item label="简介">{data?.abstract}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{data?.create_time}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{data?.latest_time}</Descriptions.Item>
        </Descriptions>
      </Card>
    </PageContainer>
  );
};

export default CompetitionDetailPage;
