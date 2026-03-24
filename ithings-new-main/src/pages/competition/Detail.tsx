import React, { useEffect, useState } from 'react';
import { Card, Descriptions, Button, message, Table } from 'antd';
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'umi';
import type { ColumnsType } from 'antd/es/table';
import { getCompetition } from '@/services/ladder/competition';
import { getCompetitionRecords } from '@/services/ladder/record';

interface CompetitionItem {
  id: string;
  title: string;
  abstract: string;
  sort_order: string;
  time_limit: number;
  status: string;
  create_time: string;
  latest_time: string;
}

interface RecordItem {
  id: string;
  user?: { id: string; username: string };
  competition?: { id: string };
  score: number;
  run_time: number;
  status: string;
  create_time: string;
  rank?: number;
}

const Detail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [competition, setCompetition] = useState<CompetitionItem | null>(null);
  const [records, setRecords] = useState<RecordItem[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const fetchCompetitionDetail = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const res: any = await getCompetition(id);
      if (res?.code === 200) {
        setCompetition(res?.data);
      } else {
        message.error(res?.message || '获取比赛详情失败');
      }
    } catch (error) {
      message.error('获取比赛详情失败');
    } finally {
      setLoading(false);
    }
  };

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
      fetchCompetitionDetail();
      fetchRecords();
    }
  }, [id]);

  if (loading && !competition) {
    return <div>加载中...</div>;
  }

  if (!competition) {
    return <div>无法获取比赛详情</div>;
  }

  const columns: ColumnsType<RecordItem> = [
    {
      title: '排名',
      dataIndex: 'rank',
      key: 'rank',
      render: (rank, record, index) => (record.rank != null ? record.rank : index + 1),
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
        <Button
          type="link"
          onClick={() => navigate(`/competition/${id}/record/${record.id}`)}
        >
          查看详情
        </Button>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/competition')}
        style={{ marginBottom: 16 }}
      >
        返回比赛列表
      </Button>

      <Card title="比赛详情">
        <Descriptions bordered>
          <Descriptions.Item label="标题">{competition.title}</Descriptions.Item>
          <Descriptions.Item label="简介">{competition.abstract}</Descriptions.Item>
          <Descriptions.Item label="排序方式">{competition.sort_order}</Descriptions.Item>
          <Descriptions.Item label="时间限制">{competition.time_limit}s</Descriptions.Item>
          <Descriptions.Item label="状态">{competition.status}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{competition.create_time}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card
        title="提交排名"
        extra={
          <Button
            type="primary"
            icon={<UploadOutlined />}
            onClick={() => navigate(`/record/submit/${id}`)}
          >
            提交代码
          </Button>
        }
        style={{ marginTop: 16 }}
      >
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

export default Detail;
