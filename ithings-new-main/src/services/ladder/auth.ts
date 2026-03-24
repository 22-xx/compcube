import request from '@/utils/request';
import { jsonToFormData } from '@/utils/formDataConverter';

export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  password: string;
  email: string;
  school: string;
}

export function login(params: LoginParams) {
  return request('/login', {
    method: 'POST',
    data: jsonToFormData(params),
  });
}

export function register(params: RegisterParams) {
  return request('/register', {
    method: 'POST',
    data: jsonToFormData(params),
  });
}

export function getLoginInfo() {
  return request('/getInfo', {
    method: 'OPTIONS',
  });
}

