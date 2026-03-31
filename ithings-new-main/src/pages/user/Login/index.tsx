import { login } from '@/services/competitionPlatform';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-form';
import { history, useModel } from '@umijs/max';
import { Col, message, Row } from 'antd';
import queryString from 'query-string';
import React from 'react';
import styles from './index.less';

const Login: React.FC = () => {
  const { initialState, setInitialState } = useModel('@@initialState');

  const handleSubmit = async (values: { username: string; password: string }) => {
    try {
      await login(values);
      const userInfo = await initialState?.fetchUserInfo?.();
      if (userInfo) {
        await setInitialState((state) => ({ ...state, currentUser: userInfo }));
      }
      message.success('登录成功');
      const query = queryString.parse(history.location.search);
      const { redirect } = query as { redirect?: string };
      history.push(redirect || '/competition');
    } catch (error) {
      message.error('登录失败，请检查用户名和密码');
    }
  };

  return (
    <div className={styles.container}>
      <Row align="middle">
        <img width={36} height={36} src="/favicon.ico" />
        <div className={styles.title}>Ladder Competition Platform</div>
      </Row>
      <Row className={styles['row-content']}>
        <Col flex="auto">
          <div className={styles['content-wrap']}>
            <div className={styles.content}>
              <p className={styles.title}>用户登录</p>
              <p className={styles['sub-title']}>使用后端当前配置完成本地联调</p>
              <LoginForm
                initialValues={{
                  username: 'administrator',
                  password: 'iThings666',
                }}
                submitter={{
                  searchConfig: {
                    submitText: '登录',
                  },
                }}
                onFinish={async (values) => {
                  await handleSubmit(values as { username: string; password: string });
                }}
              >
                <ProFormText
                  name="username"
                  fieldProps={{
                    size: 'large',
                    prefix: <UserOutlined className={styles.prefixIcon} />,
                  }}
                  placeholder="用户名"
                  rules={[
                    {
                      required: true,
                      message: '请输入用户名',
                    },
                  ]}
                />
                <ProFormText.Password
                  name="password"
                  fieldProps={{
                    size: 'large',
                    prefix: <LockOutlined className={styles.prefixIcon} />,
                  }}
                  placeholder="密码"
                  rules={[
                    {
                      required: true,
                      message: '请输入密码',
                    },
                  ]}
                />
              </LoginForm>
              <a onClick={() => history.push('/user/register')}>没有账号？去注册</a>
            </div>
          </div>
        </Col>
      </Row>
    </div>
  );
};

export default Login;
