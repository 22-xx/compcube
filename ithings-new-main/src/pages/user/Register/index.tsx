import React, { useState } from 'react';
import { Button, Checkbox, Form, Input, message, Card } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'umi';
import styles from './index.less';

const Register: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      // 调用注册接口
      const response = await fetch('/register', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        body: new URLSearchParams(values as any),
      });
      
      const data = await response.json();
      if (data.code === 0) {
        message.success('注册成功');
        navigate('/user/login');
      } else {
        message.error(data.msg || '注册失败');
      }
    } catch (error) {
      message.error('注册失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card} title="用户注册">
        <Form
          name="register"
          className={styles.form}
          initialValues={{ remember: true }}
          onFinish={onFinish}
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名！' }]}
          >
            <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="用户名" />
          </Form.Item>
          
          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱！' },
              { type: 'email', message: '请输入有效的邮箱地址！' },
            ]}
          >
            <Input prefix={<MailOutlined className="site-form-item-icon" />} placeholder="邮箱" />
          </Form.Item>
          
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码！' }]}
          >
            <Input
              prefix={<LockOutlined className="site-form-item-icon" />}
              type="password"
              placeholder="密码"
            />
          </Form.Item>
          
          <Form.Item
            name="confirmPassword"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码！' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致！'));
                },
              }),
            ]}
          >
            <Input
              prefix={<LockOutlined className="site-form-item-icon" />}
              type="password"
              placeholder="确认密码"
            />
          </Form.Item>
          
          <Form.Item>
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>我已阅读并同意用户协议</Checkbox>
            </Form.Item>
          </Form.Item>
          
          <Form.Item>
            <Button type="primary" htmlType="submit" className={styles.button} loading={loading}>
              注册
            </Button>
            已有账号？ <Link to="/user/login">立即登录</Link>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Register;