import request from '@/utils/request';

export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

export interface UserInfo {
  id: string;
  username: string;
  roles: 'admin' | 'user' | string;
  source?: string;
  school?: string;
  email?: string;
}

export interface CompetitionInfo {
  id: string;
  author?: UserInfo;
  title: string;
  abstract: string;
  sort_order?: string;
  time_limit: number;
  status: string;
  create_time: string;
  latest_time: string;
}

export interface RecordInfo {
  id: string;
  user: UserInfo;
  competition: CompetitionInfo;
  status: string;
  run_time: number;
  score: number;
  errors: string;
  create_time: string;
  latest_time: string;
  finish_time: string;
  rank?: number;
}

export const login = (payload: { username: string; password: string }) =>
  request<ApiResponse<UserInfo>>('/login', {
    method: 'POST',
    requestType: 'form',
    data: payload,
  });

export const register = (payload: {
  username: string;
  password: string;
  email: string;
  school: string;
}) =>
  request<ApiResponse<string>>('/register', {
    method: 'POST',
    requestType: 'form',
    data: payload,
  });

export const getCurrentUser = () =>
  request<ApiResponse<UserInfo>>('/getInfo', {
    method: 'GET',
  });

export const logout = () =>
  request<ApiResponse<string>>('/logout', {
    method: 'POST',
  });

export const listCompetitions = (params?: { pageNum?: number; pageSize?: number }) =>
  request<ApiResponse<{ total: number; competitionList: CompetitionInfo[] }>>('/competition', {
    method: 'GET',
    params,
  });

export const getCompetition = (competitionID: string) =>
  request<ApiResponse<CompetitionInfo>>(`/competition/${competitionID}`, {
    method: 'GET',
  });

export const createCompetition = (payload: Record<string, any>) =>
  request<ApiResponse<string>>('/competition', {
    method: 'POST',
    requestType: 'form',
    data: payload,
  });

export const updateCompetition = (competitionID: string, payload: Record<string, any>) =>
  request<ApiResponse<string>>(`/competition/${competitionID}`, {
    method: 'PUT',
    requestType: 'form',
    data: payload,
  });

export const deleteCompetition = (competitionID: string) =>
  request<ApiResponse<string>>(`/competition/${competitionID}`, {
    method: 'DELETE',
  });

export const listRecords = (params?: { pageNum?: number; pageSize?: number }) =>
  request<ApiResponse<{ total: number; recordList: RecordInfo[] }>>('/record', {
    method: 'GET',
    params,
  });

export const getCompetitionRank = (
  competitionID: string,
  params?: { pageNum?: number; pageSize?: number },
) =>
  request<ApiResponse<{ total: number; recordList: RecordInfo[] }>>(
    `/competition/${competitionID}/record`,
    {
      method: 'GET',
      params,
    },
  );

export const getRecord = (competitionID: string, recordID: string) =>
  request<ApiResponse<RecordInfo>>(`/competition/${competitionID}/record/${recordID}`, {
    method: 'GET',
  });

export const submitRecord = async (competitionID: string, file: File) => {
  const formData = new FormData();
  formData.append('submission', file);
  return request<ApiResponse<string>>(`/competition/${competitionID}/record`, {
    method: 'POST',
    data: formData,
  });
};

export const listUsers = (params?: { pageNum?: number; pageSize?: number }) =>
  request<ApiResponse<{ total: number; userList: UserInfo[] }>>('/user', {
    method: 'GET',
    params,
  });

export const getUserProfile = (userID = '0') =>
  request<ApiResponse<UserInfo>>(`/user/${userID}`, {
    method: 'GET',
  });

export const createUser = (payload: Record<string, any>) =>
  request<ApiResponse<string>>('/user', {
    method: 'POST',
    requestType: 'form',
    data: payload,
  });

export const updateUser = (userID: string, payload: Record<string, any>) =>
  request<ApiResponse<string>>(`/user/${userID}`, {
    method: 'PUT',
    requestType: 'form',
    data: payload,
  });

export const deleteUser = (userID: string) =>
  request<ApiResponse<string>>(`/user/${userID}`, {
    method: 'DELETE',
  });
