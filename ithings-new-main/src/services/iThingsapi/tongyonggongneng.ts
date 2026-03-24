import request from '@/utils/request';

export type UploadFileParams = {
  business: string;
  scene: string;
  filePath: string;
  file: File;
};

export async function postApiV1SystemCommonUploadFile(params: UploadFileParams) {
  const formData = new FormData();
  formData.append('business', params.business);
  formData.append('scene', params.scene);
  formData.append('filePath', params.filePath);
  formData.append('file', params.file);

  return request('/api/v1/system/common/uploadFile', {
    method: 'POST',
    data: formData,
  });
}
