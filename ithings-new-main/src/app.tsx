import type { Settings as LayoutSettings } from '@ant-design/pro-layout';
import { history } from '@umijs/max';
import { getLoginInfo } from './services/ladder/auth';

const loginPath = '/user/login';

/**
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: any;
  fetchUserInfo?: () => Promise<
    any | undefined
  >;
}> {
  const fetchUserInfo = async () => {
    try {
      // Ladder 后端基于 Cookie 认证：只要 cookie 存在，就能拿到当前用户
      const res: any = await getLoginInfo();
      if (res?.code === 200) {
        return res.data;
      }
      history.push(loginPath);
    } catch (error) {
      history.push(loginPath);
    }
    return undefined;
  };
  // 如果是登录页面，不执行
  if (history.location.pathname !== loginPath) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings: {},
    };
  }
  return {
    fetchUserInfo,
    settings: {},
  };
}
