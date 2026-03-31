import { PageContainer } from '@ant-design/pro-components';
import { Space, Table, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import { history } from '@umijs/max';
import type { RecordInfo } from '@/services/competitionPlatform';
import { listRecords } from '@/services/competitionPlatform';

const RecordListPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [dataSource, setDataSource] = useState<RecordInfo[]>([]);
  const [total, setTotal] = useState(0);
  const [pageNum, setPageNum] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const loadData = async (nextPageNum = pageNum, nextPageSize = pageSize) => {
    setLoading(true);
    try {
      const { data } = await listRecords({ pageNum: nextPageNum, pageSize: nextPageSize });
      setDataSource(data.recordList);
      setTotal(data.total);
    } catch (error) {
      message.error('获取提交记录失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [pageNum, pageSize]);

  const columns: ColumnsType<RecordInfo> = [
    { title: '比赛', render: (_, record) => record.competition?.title || '-' },
    { title: '用户', render: (_, record) => record.user?.username || '-' },
    { title: '状态', dataIndex: 'status' },
    { title: '得分', dataIndex: 'score' },
    { title: '运行时间', dataIndex: 'run_time' },
    { title: '提交时间', dataIndex: 'create_time' },
    {
      title: '操作',
      render: (_, record) => (
        <Space>
          <a
            onClick={() =>
              history.push(`/competition/${record.competition.id}/record/${record.id}`)
            }
          >
            详情
          </a>
          <a onClick={() => history.push(`/record/rank/${record.competition.id}`)}>排行</a>
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={dataSource}
        pagination={{
          current: pageNum,
          pageSize,
          total,
          onChange: (nextPageNum, nextPageSize) => {
            setPageNum(nextPageNum);
            setPageSize(nextPageSize);
          },
        }}
      />
    </PageContainer>
  );
};

export default RecordListPage;
