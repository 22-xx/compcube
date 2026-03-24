import React, { useEffect, useState } from 'react';
import { Card, Table, message } from 'antd';
import { useNavigate } from 'umi';
import type { ColumnsType } from 'antd/es/table';
import { getRecords } from '@/services/ladder/record';

// 后端返回：id, user: {id,username,roles,source}, competition: {...}, status, run_time, score, errors, create_time, latest_time, finish_time
interface RecordItem {
  id: string;
  user?: { id: string; username: string };
  competition?: { id: string; title: string };
  score: number;
  run_time: number;
  status: string;
  errors: string;
  create_time: string;
}

const List: React.FC = () => {
  const [records, setRecords] = useState<RecordItem[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const fetchRecords = async () => {
    setLoading(true);
    try {
      const res: any = await getRecords({ pageNum: 1, pageSize: 100 });
      if (res?.code === 200) {
        setRecords(res?.data?.recordList ?? []);
      } else {
        message.error(res?.message || '获取提交记录失败');
      }
    } catch (error) {
      message.error('获取提交记录失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecords();
  }, []);

  const columns: ColumnsType<RecordItem> = [
    {
      title: '比赛',
      key: 'competitionTitle',
      render: (_, record) => {
        const cid = record.competition?.id;
        const title = record.competition?.title || '-';
        return cid ? (
          <a onClick={() => navigate(`/competition/detail/${cid}`)}>{title}</a>
        ) : (
          title
        );
      },
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
      title: '状态',
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => {
        const cid = record.competition?.id;
        if (!cid) return '-';
        return (
          <a
            onClick={() =>
              navigate(`/competition/${cid}/record/${record.id}`)
            }
          >
            查看详情
          </a>
        );
      },
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card title="提交记录">
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

export default List;
