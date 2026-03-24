import React, { useEffect, useState } from 'react';
import { Card, Descriptions, Button, message } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'umi';
import { getRecord } from '@/services/ladder/record';

// 后端返回：id, user, competition, status, run_time, score, errors, create_time, latest_time, finish_time
interface RecordItem {
  id: string;
  user?: { id: string; username: string };
  competition?: { id: string; title: string };
  score: number;
  run_time: number;
  status: string;
  errors: string;
  create_time: string;
  finish_time: string;
}

const Detail: React.FC = () => {
  const { cid, rid } = useParams<{ cid: string; rid: string }>();
  const [record, setRecord] = useState<RecordItem | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const fetchRecordDetail = async () => {
    if (!cid || !rid) return;
    setLoading(true);
    try {
      const res: any = await getRecord(cid, rid);
      if (res?.code === 200) {
        setRecord(res?.data);
      } else {
        message.error(res?.message || '获取提交详情失败');
      }
    } catch (error) {
      message.error('获取提交详情失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (cid && rid) {
      fetchRecordDetail();
    }
  }, [cid, rid]);

  if (loading && !record) {
    return <div>加载中...</div>;
  }

  if (!record) {
    return (
      <div style={{ padding: '24px' }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/record')} style={{ marginBottom: 16 }}>
          返回提交记录
        </Button>
        <div>无法获取提交详情</div>
      </div>
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/record')}
        style={{ marginBottom: 16 }}
      >
        返回提交记录
      </Button>

      <Card title="提交详情">
        <Descriptions bordered>
          <Descriptions.Item label="比赛">
            {cid ? (
              <a onClick={() => navigate(`/competition/detail/${cid}`)}>
                {record.competition?.title ?? '-'}
              </a>
            ) : (
              record.competition?.title ?? '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label="用户">{record.user?.username ?? '-'}</Descriptions.Item>
          <Descriptions.Item label="分数">{record.score}</Descriptions.Item>
          <Descriptions.Item label="运行时间">{record.run_time}ms</Descriptions.Item>
          <Descriptions.Item label="状态">{record.status}</Descriptions.Item>
          <Descriptions.Item label="提交时间">{record.create_time}</Descriptions.Item>
          {record.errors && (
            <Descriptions.Item label="错误信息" span={3}>
              {record.errors}
            </Descriptions.Item>
          )}
        </Descriptions>
      </Card>
    </div>
  );
};

export default Detail;
