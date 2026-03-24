import React, { useEffect, useState } from 'react';
import { Card, Table, Button, message } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'umi';
import type { ColumnsType } from 'antd/es/table';
import { getCompetitionRecords } from '@/services/ladder/record';

interface RecordItem {
  id: string;
  user?: { id: string; username: string };
  score: number;
  run_time: number;
  status: string;
  create_time: string;
  rank?: number;
}

const Rank: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [records, setRecords] = useState<RecordItem[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const fetchRecords = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const res: any = await getCompetitionRecords(id, { pageNum: 1, pageSize: 100 });
      if (res?.code === 200) {
        setRecords(res?.data?.recordList ?? []);
      } else {
        message.error(res?.message || '获取排名失败');
      }
    } catch (error) {
      message.error('获取排名失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (id) {
      fetchRecords();
    }
  }, [id]);

  const columns: ColumnsType<RecordItem> = [
    {
      title: '排名',
      dataIndex: 'rank',
      key: 'rank',
      render: (_, record, index) => (record.rank != null ? record.rank : index + 1),
    },
    {
      title: '用户',
      key: 'username',
      render: (_, record) => record.user?.username ?? '-',
    },
    {
      title: '分数',
      dataIndex: 'score',
      key: 'score',
      sorter: (a, b) => a.score - b.score,
      defaultSortOrder: 'descend',
    },
    {
      title: '提交时间',
      dataIndex: 'create_time',
      key: 'create_time',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <a onClick={() => navigate(`/competition/${id}/record/${record.id}`)}>
          查看详情
        </a>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate(`/competition/detail/${id}`)}
        style={{ marginBottom: 16 }}
      >
        返回比赛详情
      </Button>

      <Card title="比赛排名">
        <Table
          columns={columns}
          dataSource={records}
          rowKey="id"
          loading={loading}
        />
      </Card>
    </div>
  );
};

export default Rank;
