import request from '@/utils/request';
import { jsonToFormData } from '@/utils/formDataConverter';

export interface RecordListParams {
  pageNum?: number;
  pageSize?: number;
}

export function getRecords(params: RecordListParams) {
  return request('/record', {
    method: 'GET',
    params,
  });
}

export function getCompetitionRecords(competitionID: string, params: RecordListParams) {
  return request(`/competition/${competitionID}/record`, {
    method: 'GET',
    params,
  });
}

export function submitRecord(competitionID: string, file: File) {
  const formData = new FormData();
  formData.append('submission', file);

  return request(`/competition/${competitionID}/record`, {
    method: 'POST',
    data: formData,
  });
}

export function getRecord(competitionID: string, recordID: string) {
  return request(`/competition/${competitionID}/record/${recordID}`, {
    method: 'GET',
  });
}

export interface UpdateRecordParams {
  score: string;
  error: string;
  runTime: string;
  status: string;
}

export function updateRecord(
  competitionID: string,
  recordID: string,
  data: UpdateRecordParams,
) {
  return request(`/competition/${competitionID}/record/${recordID}`, {
    method: 'PUT',
    data: jsonToFormData(data),
  });
}

export function deleteRecord(competitionID: string, recordID: string) {
  return request(`/competition/${competitionID}/record/${recordID}`, {
    method: 'DELETE',
  });
}

