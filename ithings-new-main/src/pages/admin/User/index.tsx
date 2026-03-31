import { PageContainer } from '@ant-design/pro-components';
import { Button, Form, Input, Modal, Popconfirm, Select, Space, Table, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import type { UserInfo } from '@/services/competitionPlatform';
import { createUser, deleteUser, listUsers } from '@/services/competitionPlatform';

const AdminUserPage: React.FC = () => {
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [dataSource, setDataSource] = useState<UserInfo[]>([]);

  const loadData = async () => {
    setLoading(true);
    try {
      const { data } = await listUsers();
      setDataSource(data.userList);
    } catch (error) {
      message.error('获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const columns: ColumnsType<UserInfo> = [
    { title: '用户名', dataIndex: 'username' },
    { title: '角色', dataIndex: 'roles' },
    { title: '学校', dataIndex: 'school' },
    { title: '邮箱', dataIndex: 'email' },
    {
      title: '操作',
      render: (_, record) => (
        <Space>
          <Popconfirm
            title="确认删除该用户？"
            onConfirm={async () => {
              try {
                await deleteUser(record.id);
                message.success('删除成功');
                await loadData();
              } catch (error) {
                message.error('删除失败');
              }
            }}
          >
            <a>删除</a>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <PageContainer
      extra={[
        <Button key="create" type="primary" onClick={() => setOpen(true)}>
          新建用户
        </Button>,
      ]}
    >
      <Table rowKey="id" loading={loading} columns={columns} dataSource={dataSource} />
      <Modal
        title="新建用户"
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={async (values) => {
            try {
              await createUser(values);
              message.success('创建成功');
              setOpen(false);
              form.resetFields();
              await loadData();
            } catch (error) {
              message.error('创建失败');
            }
          }}
        >
          <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="school" label="学校" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="role" label="角色" initialValue="user">
            <Select
              options={[
                { label: '普通用户', value: 'user' },
                { label: '管理员', value: 'admin' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default AdminUserPage;
