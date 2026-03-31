import { PageContainer } from '@ant-design/pro-components';
import { Button, Card, Upload, message } from 'antd';
import type { UploadFile } from 'antd/es/upload/interface';
import React, { useState } from 'react';
import { history, useParams } from '@umijs/max';
import { submitRecord } from '@/services/competitionPlatform';

const SubmitPage: React.FC = () => {
  const params = useParams<{ id: string }>();
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [submitting, setSubmitting] = useState(false);

  return (
    <PageContainer>
      <Card title="提交代码">
        <Upload.Dragger
          multiple={false}
          accept=".zip"
          beforeUpload={(file) => {
            setFileList([file]);
            return false;
          }}
          fileList={fileList}
          onRemove={() => {
            setFileList([]);
          }}
        >
          <p>将压缩包拖到这里，或点击选择文件</p>
          <p>后端要求提交文件为 `.zip`。</p>
        </Upload.Dragger>
        <Button
          type="primary"
          loading={submitting}
          style={{ marginTop: 16 }}
          onClick={async () => {
            if (!params.id || fileList.length === 0 || !fileList[0].originFileObj) {
              message.warning('请先选择一个 zip 文件');
              return;
            }
            setSubmitting(true);
            try {
              await submitRecord(params.id, fileList[0].originFileObj as File);
              message.success('提交成功');
              history.push('/record');
            } catch (error) {
              message.error('提交失败');
            } finally {
              setSubmitting(false);
            }
          }}
        >
          上传并提交
        </Button>
      </Card>
    </PageContainer>
  );
};

export default SubmitPage;
