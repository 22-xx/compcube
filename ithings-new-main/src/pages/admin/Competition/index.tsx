import React, { useEffect, useState } from 'react';
import { Card, Table, Button, message, Modal, Form, Input, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'umi';
import type { ColumnsType } from 'antd/es/table';
import {
  getCompetitions,
  createCompetition,
  updateCompetition,
  deleteCompetition,
} from '@/services/ladder/competition';

interface CompetitionItem {
  id: string;
  title: string;
  abstract: string;
  sort_order: string;
  time_limit: number;
  status: string;
  create_time: string;
}

const CompetitionManagement: React.FC = () => {
  const [competitions, setCompetitions] = useState<CompetitionItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingCompetition, setEditingCompetition] = useState<CompetitionItem | null>(null);
  const [form] = Form.useForm();
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

  const handleCreate = () => {
    setEditingCompetition(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const handleEdit = (competition: CompetitionItem) => {
    setEditingCompetition(competition);
    form.setFieldsValue({
      title: competition.title,
      abstract: competition.abstract,
      sortOrder: competition.sort_order,
      timeLimit: String(competition.time_limit ?? 20),
      dockerImage: (competition as any).docker_image || '',
      status: competition.status,
    });
    setIsModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      const res: any = await deleteCompetition(id);
      if (res?.code === 200) {
        message.success('删除成功');
        fetchCompetitions();
      } else {
        message.error(res?.message || '删除失败');
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      const payload = {
        title: values.title,
        abstract: values.abstract,
        sortOrder: values.sortOrder || '降序',
        timeLimit: values.timeLimit || '20',
        dockerImage: values.dockerImage,
        status: values.status,
      };
      let res: any;
      if (editingCompetition) {
        res = await updateCompetition(editingCompetition.id, payload);
      } else {
        res = await createCompetition(payload);
      }
      if (res?.code === 200) {
        message.success(editingCompetition ? '更新成功' : '创建成功');
        setIsModalVisible(false);
        fetchCompetitions();
      } else {
        message.error(res?.message || (editingCompetition ? '更新失败' : '创建失败'));
      }
    } catch (error) {
      message.error(editingCompetition ? '更新失败' : '创建失败');
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
        <div>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title="比赛管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            新建比赛
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={competitions}
          rowKey="id"
          loading={loading}
        />
      </Card>

      <Modal
        title={editingCompetition ? '编辑比赛' : '新建比赛'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={onFinish} layout="vertical">
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入比赛标题！' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="abstract"
            label="简介（文档链接）"
            rules={[{ required: true, message: '请输入比赛简介！' }]}
          >
            <Input.TextArea rows={3} placeholder="比赛简介或文档链接" />
          </Form.Item>
          <Form.Item name="sortOrder" label="成绩排序方式">
            <Select>
              <Select.Option value="升序">升序</Select.Option>
              <Select.Option value="降序">降序</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="timeLimit" label="时间要求(秒)">
            <Input placeholder="如 20" />
          </Form.Item>
          <Form.Item
            name="dockerImage"
            label="Docker 镜像名"
            rules={editingCompetition ? [] : [{ required: true, message: '请输入 Docker 镜像名！' }]}
          >
            <Input placeholder="赛题 docker 镜像名（新建时必填）" />
          </Form.Item>
          {editingCompetition && (
            <Form.Item name="status" label="状态">
              <Select>
                <Select.Option value="准备中">准备中</Select.Option>
                <Select.Option value="进行中">进行中</Select.Option>
                <Select.Option value="已结束">已结束</Select.Option>
              </Select>
            </Form.Item>
          )}
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {editingCompetition ? '更新' : '创建'}
            </Button>
            <Button style={{ marginLeft: 8 }} onClick={() => setIsModalVisible(false)}>
              取消
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CompetitionManagement;
