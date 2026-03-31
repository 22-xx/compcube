import { PageContainer } from '@ant-design/pro-components';
import { Button, Form, Input, InputNumber, Modal, Popconfirm, Space, Table, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import { history, useModel } from '@umijs/max';
import type { CompetitionInfo } from '@/services/competitionPlatform';
import {
  createCompetition,
  deleteCompetition,
  listCompetitions,
  updateCompetition,
} from '@/services/competitionPlatform';

const CompetitionListPage: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const canAdmin = initialState?.currentUser?.userInfo?.roles === 'admin';
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [dataSource, setDataSource] = useState<CompetitionInfo[]>([]);
  const [total, setTotal] = useState(0);
  const [pageNum, setPageNum] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [editing, setEditing] = useState<CompetitionInfo | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const loadData = async (nextPageNum = pageNum, nextPageSize = pageSize) => {
    setLoading(true);
    try {
      const { data } = await listCompetitions({ pageNum: nextPageNum, pageSize: nextPageSize });
      setDataSource(data.competitionList);
      setTotal(data.total);
    } catch (error) {
      message.error('获取比赛列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [pageNum, pageSize]);

  const columns: ColumnsType<CompetitionInfo> = [
    { title: '标题', dataIndex: 'title' },
    { title: '作者', render: (_, record) => record.author?.username || '-' },
    { title: '状态', dataIndex: 'status' },
    { title: '时限', dataIndex: 'time_limit' },
    { title: '创建时间', dataIndex: 'create_time' },
    {
      title: '操作',
      render: (_, record) => (
        <Space>
          <a onClick={() => history.push(`/competition/detail/${record.id}`)}>详情</a>
          <a onClick={() => history.push(`/record/rank/${record.id}`)}>排行</a>
          <a onClick={() => history.push(`/record/submit/${record.id}`)}>提交</a>
          {canAdmin && (
            <a
              onClick={() => {
                setEditing(record);
                form.setFieldsValue({
                  title: record.title,
                  abstract: record.abstract,
                  dockerImage: '',
                  timeLimit: record.time_limit,
                });
                setModalOpen(true);
              }}
            >
              编辑
            </a>
          )}
          {canAdmin && (
            <Popconfirm
              title="确认删除该比赛？"
              onConfirm={async () => {
                try {
                  await deleteCompetition(record.id);
                  message.success('删除成功');
                  await loadData();
                } catch (error) {
                  message.error('删除失败');
                }
              }}
            >
              <a>删除</a>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <PageContainer
      extra={
        canAdmin
          ? [
              <Button
                key="create"
                type="primary"
                onClick={() => {
                  setEditing(null);
                  form.resetFields();
                  form.setFieldsValue({ timeLimit: 20, dockerImage: 'python:3.10' });
                  setModalOpen(true);
                }}
              >
                新建比赛
              </Button>,
            ]
          : undefined
      }
    >
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
      <Modal
        title={editing ? '编辑比赛' : '新建比赛'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={async (values) => {
            try {
              if (editing) {
                await updateCompetition(editing.id, values);
                message.success('更新成功');
              } else {
                await createCompetition(values);
                message.success('创建成功');
              }
              setModalOpen(false);
              await loadData();
            } catch (error) {
              message.error(editing ? '更新失败' : '创建失败');
            }
          }}
        >
          <Form.Item name="title" label="标题" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="abstract" label="简介" rules={[{ required: true }]}>
            <Input.TextArea rows={4} />
          </Form.Item>
          <Form.Item name="dockerImage" label="Docker 镜像" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="timeLimit" label="时间限制" rules={[{ required: true }]}>
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default CompetitionListPage;
