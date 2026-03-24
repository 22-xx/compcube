/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */
export default function access(initialState: {
  currentUser?: { roles?: string; source?: string } | undefined;
}) {
  const { currentUser } = initialState || {};

  return {
    canAdmin: currentUser?.roles === 'admin',
  };
}
