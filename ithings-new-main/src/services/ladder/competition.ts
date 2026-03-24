import request from '@/utils/request';
import { jsonToFormData } from '@/utils/formDataConverter';

export interface CompetitionListParams {
  pageNum?: number;
  pageSize?: number;
}

export interface CreateCompetitionParams {
  title: string;
  abstract: string;
  sortOrder?: '升序' | '降序';
  timeLimit?: string;
  dockerImage: string;
}

export interface UpdateCompetitionParams {
  title?: string;
  abstract?: string;
  sortOrder?: '升序' | '降序';
  timeLimit?: string;
  dockerImage?: string;
  status?: '准备中' | '进行中' | '已结束';
}

export function getCompetitions(params: CompetitionListParams) {
  return request('/competition', {
    method: 'GET',
    params,
  });
}

export function createCompetition(data: CreateCompetitionParams) {
  return request('/competition', {
    method: 'POST',
    data: jsonToFormData(data),
  });
}

export function getCompetition(competitionID: string) {
  return request(`/competition/${competitionID}`, {
    method: 'GET',
  });
}

export function updateCompetition(competitionID: string, data: UpdateCompetitionParams) {
  return request(`/competition/${competitionID}`, {
    method: 'PUT',
    data: jsonToFormData(data),
  });
}

export function deleteCompetition(competitionID: string) {
  return request(`/competition/${competitionID}`, {
    method: 'DELETE',
  });
}

