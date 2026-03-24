import { defineConfig } from '@umijs/max';
import defaultSettings from './defaultSettings';
import proxy from './proxy';
import routes from './routes';

// const { REACT_APP_ENV } = process.env;

export default defineConfig({
  hash: true,
  history: {type: 'hash'},
  antd: {},
  clientLoader: {},
  /*
    front/iThingsCore/    iThings核心前端路由
    front/iThingsCompany/  iThings企业版前端路由
    front/custom/chengde/  企业版定制版前端路由
  */
  // publicPath: '/front/iThingsCore/',
  publicPath: '/',
  dva: {},
  // dynamicImport: {
  //   loading: '@ant-design/pro-layout/es/PageLoading',
  // },
  routes,
  // Theme for antd: https://ant.design/docs/react/customize-theme-cn
  theme: {
    'primary-color': defaultSettings.colorPrimary,
  },
  // esbuild is father build tools
  // https://umijs.org/plugins/plugin-esbuild
  // esbuild: {},
  title: 'iThings',
  request: {},
  initialState: {},
  presets: ['umi-presets-pro'],
  model: {},
  ignoreMomentLocale: true,
  // proxy: proxy[REACT_APP_ENV || 'dev'],
  proxy: proxy['dev'],
  manifest: {
    basePath: '/',
  },
  fastRefresh: true,
  requestRecord: {},
  mfsu: {},
});
