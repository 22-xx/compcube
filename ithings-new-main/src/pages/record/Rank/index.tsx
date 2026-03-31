import { PageContainer } from '@ant-design/pro-components';
import { Table, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import { useParams } from '@umijs/max';
import type { RecordInfo } from '@/services/competitionPlatform';
import { getCompetitionRank } from '@/services/competitionPlatform';

const RankPage: React.FC = () => {
  const params = useParams<{ id: string }>();
  const [dataSource, setDataSource] = useState<RecordInfo[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      if (!params.id) return;
      setLoading(true);
      try {
        const { data } = await getCompetitionRank(params.id);
        setDataSource(data.recordList);
      } catch (error) {
        message.error('获取排行失败');
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, [params.id]);

  const columns: ColumnsType<RecordInfo> = [
    { title: '排名', dataIndex: 'rank' },
    { title: '用户', render: (_, record) => record.user?.username || '-' },
    { title: '得分', dataIndex: 'score' },
    { title: '运行时间', dataIndex: 'run_time' },
    { title: '状态', dataIndex: 'status' },
    { title: '提交时间', dataIndex: 'create_time' },
  ];

  return (
    <PageContainer>
      <Table rowKey="id" loading={loading} columns={columns} dataSource={dataSource} />
    </PageContainer>
  );
};

export default RankPage;
