const localTarget = 'http://127.0.0.1:8001';

const backendPaths = [
  '/login',
  '/logout',
  '/register',
  '/getInfo',
  '/competition',
  '/record',
  '/user',
  '/data',
  '/model',
  '/train',
  '/files',
  '/docs',
];

const buildLocalProxy = () =>
  backendPaths.reduce<Record<string, { target: string; changeOrigin: boolean }>>(
    (acc, path) => {
      acc[path] = {
        target: localTarget,
        changeOrigin: true,
      };
      return acc;
    },
    {},
  );

export default {
  dev: buildLocalProxy(),
  test: {
    '/': {
      target: 'http://47.92.240.210:8099',
      changeOrigin: true,
    },
  },
  pre: {
    '/': {
      target: 'http://47.92.240.210:8099',
      changeOrigin: true,
    },
  },
};
