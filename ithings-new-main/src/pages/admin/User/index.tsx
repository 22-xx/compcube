import React, { useEffect, useState } from 'react';
import { Card, Table, Button, message, Modal, Form, Input, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
} from '@/services/ladder/user';

// 后端返回：id, username, school, email, roles, source, is_delete, create_time, latest_time
interface UserItem {
  id: string;
  username: string;
  email: string;
  school: string;
  roles: string;
  create_time: string;
}

const UserManagement: React.FC = () => {
  const [users, setUsers] = useState<UserItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<UserItem | null>(null);
  const [form] = Form.useForm();

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res: any = await getUsers({ pageNum: 1, pageSize: 100 });
      if (res?.code === 200) {
        setUsers(res?.data?.userList ?? []);
      } else {
        message.error(res?.message || '获取用户列表失败');
      }
    } catch (error) {
      message.error('获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingUser(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const handleEdit = (user: UserItem) => {
    setEditingUser(user);
    form.setFieldsValue({
      username: user.username,
      email: user.email,
      school: user.school,
      role: user.roles,
    });
    setIsModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      const res: any = await deleteUser(id);
      if (res?.code === 200) {
        message.success('删除成功');
        fetchUsers();
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
      let res: any;
      if (editingUser) {
        res = await updateUser(editingUser.id, {
          username: values.username,
          email: values.email,
          school: values.school,
          role: values.role,
          ...(values.password ? { password: values.password } : {}),
        });
      } else {
        res = await createUser({
          username: values.username,
          password: values.password,
          email: values.email,
          school: values.school,
          role: values.role || 'user',
        });
      }
      if (res?.code === 200) {
        message.success(editingUser ? '更新成功' : '创建成功');
        setIsModalVisible(false);
        fetchUsers();
      } else {
        message.error(res?.message || (editingUser ? '更新失败' : '创建失败'));
      }
    } catch (error) {
      message.error(editingUser ? '更新失败' : '创建失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const columns: ColumnsType<UserItem> = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '学校',
      dataIndex: 'school',
      key: 'school',
    },
    {
      title: '角色',
      dataIndex: 'roles',
      key: 'roles',
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
        title="用户管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            新建用户
          </Button>
        }
      >
        <Table columns={columns} dataSource={users} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title={editingUser ? '编辑用户' : '新建用户'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={onFinish} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名！' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱！' },
              { type: 'email', message: '请输入有效的邮箱地址！' },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="school"
            label="学校"
            rules={[{ required: true, message: '请输入学校！' }]}
          >
            <Input />
          </Form.Item>
          {!editingUser && (
            <Form.Item
              name="password"
              label="密码"
              rules={[{ required: true, message: '请输入密码！' }]}
            >
              <Input.Password />
            </Form.Item>
          )}
          {editingUser && (
            <Form.Item name="password" label="新密码（不修改请留空）">
              <Input.Password placeholder="留空则不修改密码" />
            </Form.Item>
          )}
          <Form.Item name="role" label="角色" rules={[{ required: true, message: '请选择角色！' }]}>
            <Select>
              <Select.Option value="user">普通用户</Select.Option>
              <Select.Option value="admin">管理员</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {editingUser ? '更新' : '创建'}
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

export default UserManagement;
