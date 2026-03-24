import { stream } from '@/utils/request';

export function downloadFile(filename: string) {
  return stream(`/files/${encodeURIComponent(filename)}`, {
    method: 'GET',
  });
}

export function getFileUrl(filename: string) {
  return `/files/${encodeURIComponent(filename)}`;
}

