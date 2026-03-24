/**
 * request 网络请求工具
 * 更详细的 api 文档: https://github.com/umijs/umi-request
 */
import { history } from '@umijs/max';
import { notification } from 'antd';
import { stringify } from 'querystring';
import { extend } from 'umi-request';

/**
 * 异常处理程序
 */
const errorHandler = (error: { response: Response }): Response => {
  const { response } = error;

  if (response && response.status) {
    response
      .clone()
      .text()
      .then((v) => {
        const regex = /"msg":"([^"]+)"/;
        const match = v.match(regex);
        try {
          const data = JSON.parse(v);
          notification.error({
            message: `请求错误, 错误码:${data.code}`,
            description: data.message || data.msg,
          });
        } catch {
          notification.error({
            message: `请求错误, 错误码:${response.status}`,
            description: match?.[1] || v,
          });
        }
      });
  } else if (!response) {
    notification.error({
      description: '您的网络发生异常，无法连接服务器',
      message: '网络异常',
    });
  }

  return response;
};

const redirectLoginPage = () => {
  const queryString = stringify({
    redirect: window.location.href,
  });
  history.push(`/user/login?${queryString}`);
};

// 响应拦截器
const responseInterceptors = (response: any) => {
  if (response.status === 401 && window.location.pathname !== '/user/login') {
    return redirectLoginPage();
  }
  return response;
};

/**
 * 配置request请求时的默认参数
 */
const request = extend({
  errorHandler, // 默认错误处理
  credentials: 'include', // 默认请求是否带上cookie
  timeout: 200000,
});

export const stream = extend({
  credentials: 'include',
  parseResponse: false,
});

request.interceptors.response.use(responseInterceptors, { global: false });

export default request;
