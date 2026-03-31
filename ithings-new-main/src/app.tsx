import type { MenuDataItem, Settings as LayoutSettings } from '@ant-design/pro-layout';
import { history } from '@umijs/max';
import type { UserInfo } from './services/competitionPlatform';
import { getCurrentUser } from './services/competitionPlatform';

const loginPath = '/user/login';
const registerPath = '/user/register';

const buildMenu = (user: UserInfo): MenuDataItem[] => {
  const menu: MenuDataItem[] = [
    {
      path: '/competition',
      name: '比赛列表',
    },
    {
      path: '/record',
      name: '提交记录',
    },
    {
      path: '/user/profile',
      name: '个人中心',
    },
  ];

  if (user.roles === 'admin') {
    menu.push({
      path: '/admin',
      name: '管理中心',
      children: [
        {
          path: '/admin/user',
          name: '用户管理',
        },
        {
          path: '/admin/competition',
          name: '比赛管理',
        },
      ],
    });
  }

  return menu;
};

export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: { userInfo: UserInfo; menuInfo: MenuDataItem[] };
  fetchUserInfo?: () => Promise<{ userInfo: UserInfo; menuInfo: MenuDataItem[] } | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const { data: userInfo } = await getCurrentUser();
      return {
        userInfo,
        menuInfo: buildMenu(userInfo),
      };
    } catch (error) {
      return undefined;
    }
  };

  if (history.location.pathname !== loginPath && history.location.pathname !== registerPath) {
    const currentUser = await fetchUserInfo();
    if (!currentUser) {
      history.push(loginPath);
    }
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
