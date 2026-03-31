import { PageContainer } from '@ant-design/pro-components';
import { Button, Card, Form, Input, message } from 'antd';
import React, { useEffect } from 'react';
import { useModel } from '@umijs/max';
import { getUserProfile, updateUser } from '@/services/competitionPlatform';

const ProfilePage: React.FC = () => {
  const [form] = Form.useForm();
  const { initialState, setInitialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser?.userInfo;

  useEffect(() => {
    const loadProfile = async () => {
      if (!currentUser?.id) return;
      try {
        const { data } = await getUserProfile(currentUser.id);
        form.setFieldsValue({
          username: data.username,
          email: data.email,
          school: data.school,
          role: data.roles,
        });
      } catch (error) {
        message.error('获取个人信息失败');
      }
    };
    loadProfile();
  }, [currentUser?.id, form]);

  return (
    <PageContainer>
      <Card title="个人中心">
        <Form
          form={form}
          layout="vertical"
          onFinish={async (values) => {
            if (!currentUser?.id) return;
            try {
              await updateUser(currentUser.id, values);
              const { data } = await getUserProfile(currentUser.id);
              await setInitialState((state) =>
                state
                  ? {
                      ...state,
                      currentUser: {
                        ...(state.currentUser as any),
                        userInfo: data,
                      },
                    }
                  : state,
              );
              message.success('个人信息已更新');
            } catch (error) {
              message.error('更新失败');
            }
          }}
        >
          <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="school" label="学校">
            <Input />
          </Form.Item>
          <Form.Item name="password" label="新密码">
            <Input.Password placeholder="不修改可留空" />
          </Form.Item>
          <Form.Item name="role" label="角色">
            <Input disabled />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            保存
          </Button>
        </Form>
      </Card>
    </PageContainer>
  );
};

export default ProfilePage;
