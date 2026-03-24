import React, { useState } from 'react';
import { Card, Button, Upload, message } from 'antd';
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'umi';
import type { UploadProps } from 'antd';
import { submitRecord } from '@/services/ladder/record';

const Submit: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(false);
  const [fileList, setFileList] = useState<any[]>([]);
  const navigate = useNavigate();

  const handleSubmit = async () => {
    if (!id) return;
    const file = fileList[0]?.originFileObj;
    if (!file) {
      message.error('请选择 zip 文件');
      return;
    }
    setLoading(true);
    try {
      const res: any = await submitRecord(id, file);
      if (res?.code === 200) {
        message.success('提交成功');
        navigate(`/competition/detail/${id}`);
      } else {
        message.error(res?.message || '提交失败');
      }
    } catch (error) {
      message.error('提交失败');
    } finally {
      setLoading(false);
    }
  };

  const uploadProps: UploadProps = {
    name: 'submission',
    maxCount: 1,
    accept: '.zip',
    fileList,
    beforeUpload: (file) => {
      const isZip = file.name.endsWith('.zip');
      if (!isZip) {
        message.error('只能上传 zip 文件！');
        return false;
      }
      setFileList([{ uid: file.uid, name: file.name, originFileObj: file }]);
      return false; // 阻止自动上传，改为手动提交
    },
    onRemove: () => setFileList([]),
  };

  return (
    <div style={{ padding: '24px' }}>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate(`/competition/detail/${id}`)}
        style={{ marginBottom: 16 }}
      >
        返回比赛详情
      </Button>

      <Card title="提交代码">
        <p style={{ marginBottom: 16 }}>请上传包含您代码的 zip 文件。</p>
        <Upload {...uploadProps}>
          <Button icon={<UploadOutlined />}>选择 zip 文件</Button>
        </Upload>
        <div style={{ marginTop: 16 }}>
          <Button type="primary" onClick={handleSubmit} loading={loading}>
            提交
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default Submit;
