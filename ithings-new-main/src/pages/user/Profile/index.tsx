import React, { useEffect, useState } from 'react';
import { Card, Descriptions, Button, message, Form, Input } from 'antd';
import { UserOutlined, MailOutlined, EditOutlined } from '@ant-design/icons';
import { getLoginInfo } from '@/services/ladder/auth';
import { updateUser } from '@/services/ladder/user';

interface UserInfo {
  id: string;
  username: string;
  email: string;
  roles?: string;
  source?: string;
}

const Profile: React.FC = () => {
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [editing, setEditing] = useState(false);
  const [form] = Form.useForm();

  const fetchUserInfo = async () => {
    setLoading(true);
    try {
      const res: any = await getLoginInfo();
      if (res?.code === 200) {
        setUserInfo(res.data);
        form.setFieldsValue(res.data);
      } else {
        message.error(res?.message || '获取用户信息失败');
      }
    } catch (error) {
      message.error('获取用户信息失败');
    } finally {
      setLoading(false);
    }
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      const res: any = await updateUser(userInfo?.id || 'me', values);
      if (res?.code === 200) {
        message.success('更新成功');
        await fetchUserInfo();
        setEditing(false);
      } else {
        message.error(res?.message || '更新失败');
      }
    } catch (error) {
      message.error('更新失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUserInfo();
  }, []);

  if (loading && !userInfo) {
    return <div>加载中...</div>;
  }

  if (!userInfo) {
    return <div>无法获取用户信息</div>;
  }

  return (
    <div style={{ padding: '24px' }}>
      <Card title="个人中心" extra={
        editing ? (
          <Button onClick={() => setEditing(false)}>取消</Button>
        ) : (
          <Button icon={<EditOutlined />} onClick={() => setEditing(true)}>编辑</Button>
        )
      }>
        {editing ? (
          <Form form={form} onFinish={onFinish} layout="vertical">
            <Form.Item
              name="username"
              label="用户名"
              rules={[{ required: true, message: '请输入用户名！' }]}
            >
              <Input prefix={<UserOutlined />} />
            </Form.Item>
            <Form.Item
              name="email"
              label="邮箱"
              rules={[
                { required: true, message: '请输入邮箱！' },
                { type: 'email', message: '请输入有效的邮箱地址！' },
              ]}
            >
              <Input prefix={<MailOutlined />} />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                保存
              </Button>
            </Form.Item>
          </Form>
        ) : (
          <Descriptions bordered>
            <Descriptions.Item label="用户名">{userInfo.username}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{userInfo.email}</Descriptions.Item>
            <Descriptions.Item label="角色">{userInfo.roles || '-'}</Descriptions.Item>
          </Descriptions>
        )}
      </Card>
    </div>
  );
};

export default Profile;