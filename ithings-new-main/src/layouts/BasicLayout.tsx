import RightContent from '@/components/RightContent';
import type { MenuDataItem } from '@ant-design/pro-layout';
import { ProLayout } from '@ant-design/pro-layout';
import { history, Outlet, useLocation, useModel } from '@umijs/max';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/es/locale/zh_CN';
import moment from 'moment';
import 'moment/locale/zh-cn';
import React from 'react';
import FooterCom from '@/components/FooterCom';
import defaultSettings from '../../config/defaultSettings';

moment.locale('zh-cn');

const BasicLayout: React.FC = (props) => {
  const { initialState } = useModel('@@initialState');
  const location = useLocation();
  const menuTree = initialState?.currentUser?.menuInfo ?? [];

  return (
    <ProLayout
      location={{ pathname: location.pathname }}
      title="Ladder Competition Platform"
      siderWidth={220}
      rightContentRender={() => <RightContent />}
      footerRender={() => <FooterCom />}
      menuItemRender={(item, dom) => (
        <a
          onClick={() => {
            history.push(item.path ?? '/competition');
          }}
        >
          {dom}
        </a>
      )}
      menuDataRender={() => menuTree as MenuDataItem[]}
      {...props}
      {...defaultSettings}
      logo={<img src="/favicon.ico" alt="" />}
    >
      <ConfigProvider locale={zhCN}>
        <Outlet />
      </ConfigProvider>
    </ProLayout>
  );
};

export default BasicLayout;
