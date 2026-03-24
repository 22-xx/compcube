import request from '@/utils/request';
import { jsonToFormData } from '@/utils/formDataConverter';

export interface UserListParams {
  pageNum?: number;
  pageSize?: number;
}

export interface CreateUserParams {
  username: string;
  password: string;
  role?: 'admin' | 'user';
  email: string;
  school: string;
}

export interface UpdateUserParams {
  username?: string;
  password?: string;
  role?: 'admin' | 'user';
  email?: string;
  school?: string;
}

export function getUsers(params: UserListParams) {
  return request('/user', {
    method: 'GET',
    params,
  });
}

export function createUser(data: CreateUserParams) {
  return request('/user', {
    method: 'POST',
    data: jsonToFormData(data),
  });
}

export function getUser(userID: string) {
  return request(`/user/${userID}`, {
    method: 'GET',
  });
}

export function updateUser(userID: string, data: UpdateUserParams) {
  return request(`/user/${userID}`, {
    method: 'PUT',
    data: jsonToFormData(data),
  });
}

export function deleteUser(userID: string) {
  return request(`/user/${userID}`, {
    method: 'DELETE',
  });
}

