import { HeartTwoTone, SmileTwoTone } from '@ant-design/icons';
import { Alert, Card, Typography } from 'antd';
import React from 'react';

export default (): React.ReactNode => {
  return (
    <div>
      <Typography.Paragraph>这个页面只有 admin 权限才能查看。</Typography.Paragraph>
      <Card>
        <Alert
          message="更快更强的重型组件，已经发布。"
          type="success"
          showIcon
          banner
          style={{
            margin: -12,
            marginBottom: 48,
          }}
        />
        <Typography.Title
          level={2}
          style={{
            textAlign: 'center',
          }}
        >
          <SmileTwoTone /> iThings <HeartTwoTone twoToneColor="#eb2f96" /> You
        </Typography.Title>
      </Card>
      <p
        style={{
          textAlign: 'center',
          marginTop: 24,
        }}
      >
        Want to add more pages? Please refer to{' '}
        <a href="https://pro.ant.design/docs/block-cn" target="_blank" rel="noopener noreferrer">
          use block
        </a>
        。
      </p>
    </div>
  );
};
