import type { MenuDataItem } from '@ant-design/pro-layout';

/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */
export default function access(initialState: {
  currentUser?: { userInfo: { roles?: string }; menuInfo: MenuDataItem[] } | undefined;
}) {
  const { currentUser } = initialState || {};

  return {
    canAdmin: currentUser?.userInfo?.roles === 'admin',
  };
}
