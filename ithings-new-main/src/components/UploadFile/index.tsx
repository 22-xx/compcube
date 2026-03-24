import { postApiV1SystemCommonUploadFile } from '@/services/iThingsapi/tongyonggongneng';
import { ResponseCode } from '@/utils/base';
import { TOKENKEY } from '@/utils/const';
import { getToken } from '@/utils/utils';
import { ProFormUploadButton } from '@ant-design/pro-components';
import { useRequest } from 'ahooks';
import { message, Spin } from 'antd';
import { useEffect, useState } from 'react';

import type { UploadProps } from 'antd';
import type { RcFile } from 'antd/es/upload';

type UploadResponse = {
  code: number;
  data?: {
    filePath?: string;
  };
};

const UploadFile: React.FC<{
  accept: string;
  filePathPrefix: string;
  scene: string;
  business: string;
  getUploadFileData: (path: string) => void;
  handleClick?: () => void;
  label?: React.ReactNode;
}> = ({ accept, filePathPrefix, scene, business, handleClick, label, getUploadFileData }) => {
  const [processLoading, setProcessLoading] = useState(false);

  const { data, run } = useRequest<UploadResponse, [Parameters<typeof postApiV1SystemCommonUploadFile>[0]]>(
    (params) => postApiV1SystemCommonUploadFile(params),
    {
      manual: true,
      onSuccess: (res) => {
        setProcessLoading(false);
        if (res?.code === ResponseCode.SUCCESS) {
          message.success('上传成功');
        } else {
          message.error('上传失败');
        }
      },
      onError: () => {
        setProcessLoading(false);
        message.error('上传错误');
      },
    },
  );

  const customRequest: UploadProps['customRequest'] = async (info) => {
    setProcessLoading(true);
    const file = info.file as RcFile;
    run({
      business,
      scene,
      filePath: `${filePathPrefix}/${file.name}`,
      file,
    });
  };

  useEffect(() => {
    if (data?.code === ResponseCode.SUCCESS && data?.data?.filePath) {
      getUploadFileData(data.data.filePath);
    }
  }, [data, getUploadFileData]);

  return (
    <Spin spinning={processLoading}>
      <ProFormUploadButton
        width="md"
        label={label}
        accept={accept}
        fieldProps={{
          headers: {
            [TOKENKEY]: getToken(),
          },
          customRequest,
          showUploadList: false,
          progress: {
            strokeColor: {
              '0%': '#108ee9',
              '100%': '#87d068',
            },
            strokeWidth: 3,
            format: (percent) => (percent ? `${parseFloat(percent.toFixed(2))}%` : ''),
          },
        }}
        buttonProps={{
          onClick: handleClick,
        }}
      />
    </Spin>
  );
};

export default UploadFile;
