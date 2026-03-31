import { register } from '@/services/competitionPlatform';
import { LoginForm, ProFormText } from '@ant-design/pro-form';
import { history } from '@umijs/max';
import { Card, message } from 'antd';
import React from 'react';

const RegisterPage: React.FC = () => {
  return (
    <div style={{ maxWidth: 520, margin: '48px auto' }}>
      <Card title="注册账号">
        <LoginForm
          submitter={{ searchConfig: { submitText: '注册' } }}
          onFinish={async (values) => {
            try {
              await register(values as any);
              message.success('注册成功，请登录');
              history.push('/user/login');
            } catch (error) {
              message.error('注册失败，请检查用户名或邮箱是否已存在');
            }
          }}
        >
          <ProFormText
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          />
          <ProFormText.Password
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          />
          <ProFormText
            name="email"
            label="邮箱"
            rules={[{ required: true, message: '请输入邮箱' }]}
          />
          <ProFormText
            name="school"
            label="学校"
            rules={[{ required: true, message: '请输入学校' }]}
          />
        </LoginForm>
      </Card>
    </div>
  );
};

export default RegisterPage;
