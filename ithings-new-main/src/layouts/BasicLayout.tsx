import { Link, history, Outlet, useModel } from '@umijs/max';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/es/locale/zh_CN';
import React, { useState } from 'react';

const BasicLayout: React.FC = (props) => {
  const { initialState } = useModel('@@initialState');
  const [pathname, setPathname] = useState(
    window.location.pathname || '/competition'
  );

  const menuItems = [
    {
      path: '/competition',
      name: '比赛列表',
    },
    {
      path: '/record',
      name: '提交记录',
    },
  ];

  if (initialState?.currentUser?.roles === 'admin') {
    menuItems.push({
      path: '/admin',
      name: '管理中心',
    });
  }

  return (
    <div style={{ minHeight: '100vh' }}>
      <div
        style={{
          height: 64,
          background: '#001529',
          display: 'flex',
          alignItems: 'center',
          padding: '0 24px',
          justifyContent: 'space-between',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Link to="/" style={{ color: '#fff', fontSize: 18, fontWeight: 'bold' }}>
            竞赛平台
          </Link>
        </div>
        <div style={{ display: 'flex', gap: 16 }}>
          {menuItems.map((item) => (
            <a
              key={item.path}
              onClick={() => {
                setPathname(item.path);
                history.push(item.path);
              }}
              style={{
                color: pathname.startsWith(item.path) ? '#1890ff' : '#fff',
              }}
            >
              {item.name}
            </a>
          ))}
          {initialState?.currentUser ? (
            <a
              onClick={() => {
                history.push('/user/profile');
              }}
              style={{ color: '#fff' }}
            >
              {initialState.currentUser.username}
            </a>
          ) : (
            <Link to="/user/login" style={{ color: '#fff' }}>
              登录
            </Link>
          )}
        </div>
      </div>
      <div style={{ padding: 16 }}>
        <ConfigProvider locale={zhCN}>
          <Outlet />
        </ConfigProvider>
      </div>
    </div>
  );
};

export default BasicLayout;
