import React, { useEffect, useState } from 'react';
import { Card, Table, Button, message } from 'antd';
import { useNavigate } from 'umi';
import type { ColumnsType } from 'antd/es/table';
import { getCompetitions } from '@/services/ladder/competition';

// 后端返回字段为下划线，前端统一映射
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

const List: React.FC = () => {
  const [competitions, setCompetitions] = useState<CompetitionItem[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const fetchCompetitions = async () => {
    setLoading(true);
    try {
      const res: any = await getCompetitions({ pageNum: 1, pageSize: 100 });
      if (res?.code === 200) {
        setCompetitions(res?.data?.competitionList ?? []);
      } else {
        message.error(res?.message || '获取比赛列表失败');
      }
    } catch (error) {
      message.error('获取比赛列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCompetitions();
  }, []);

  const columns: ColumnsType<CompetitionItem> = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (text, record) => (
        <a onClick={() => navigate(`/competition/detail/${record.id}`)}>{text}</a>
      ),
    },
    {
      title: '简介',
      dataIndex: 'abstract',
      key: 'abstract',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      key: 'create_time',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Button type="link" onClick={() => navigate(`/competition/detail/${record.id}`)}>
          查看
        </Button>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card title="比赛列表">
        <Table
          columns={columns}
          dataSource={competitions}
          rowKey="id"
          loading={loading}
        />
      </Card>
    </div>
  );
};

export default List;
