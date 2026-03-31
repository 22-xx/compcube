import styles from '@/pages/user/Login/index.less';
import { history } from '@umijs/max';
import { Button, Col, Row } from 'antd';
import React from 'react';

const NoFoundPage: React.FC = () => (
  <div className={styles.container}>
    <Row align="middle">
      <img width={36} height={36} src="/favicon.ico" />
      <div className={styles.title}>Ladder Competition Platform</div>
    </Row>
    <Row className={styles['row-content']}>
      <Col flex="auto">
        <div className={styles.content}>
          <h1 style={{ fontSize: '50px' }}>404</h1>
          <h2>页面不存在或功能尚未开放</h2>
          <Button type="primary" onClick={() => history.push('/')}>
            返回首页
          </Button>
        </div>
      </Col>
    </Row>
  </div>
);

export default NoFoundPage;
