export default [
  {
    path: '/user',
    routes: [
      {
        path: '/user/login',
        layout: false,
        name: '登录',
        component: './user/Login',
      },
      {
        path: '/user/register',
        layout: false,
        name: '注册',
        component: './user/Register',
      },
      {
        path: '/user/profile',
        name: '个人中心',
        component: './user/Profile',
      },
      {
        path: '/user',
        redirect: '/user/login',
      },
      {
        component: '404',
      },
    ],
  },
  {
    path: '/',
    component: '../layouts/BasicLayout',
    routes: [
      {
        path: '/competition',
        name: '比赛列表',
        component: './competition/List',
        icon: 'icon_competition',
      },
      {
        path: '/competition/detail/:id',
        name: '比赛详情',
        hideInMenu: true,
        component: './competition/Detail',
      },
      {
        path: '/record',
        name: '提交记录',
        component: './record/List',
        icon: 'icon_record',
      },
      {
        path: '/record/rank/:id',
        name: '比赛排名',
        hideInMenu: true,
        component: './record/Rank',
      },
      {
        path: '/record/submit/:id',
        name: '提交代码',
        hideInMenu: true,
        component: './record/Submit',
      },
      {
        path: '/competition/:cid/record/:rid',
        name: '提交详情',
        hideInMenu: true,
        component: './record/Detail',
      },
      {
        path: '/admin',
        name: '管理中心',
        icon: 'icon_admin',
        access: 'canAdmin',
        routes: [
          {
            name: '用户管理',
            path: '/admin/user',
            component: './admin/User',
          },
          {
            name: '比赛管理',
            path: '/admin/competition',
            component: './admin/Competition',
          },
          {
            component: '404',
          },
        ],
      },
      {
        path: '/',
        redirect: '/competition',
      },
      {
        component: '404',
      },
    ],
  },
  {
    layout: false,
    name: '404',
    path: '/*',
    component: '404',
  },
];
