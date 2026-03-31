import { postApiV1ThingsDeviceInfoCount } from '@/services/iThingsapi/shebeiguanli';
import { useRequest } from 'ahooks';
import { Card, Col, message, Row } from 'antd';
import { useMemo } from 'react';
import DeviceChart from './pages/deviceChart/index';
import DeviceCount from './pages/deviceCount/index';
import DeviceMap from './pages/deviceMap/index';
import Pie from './pages/echarts/index'
import MyChart from "./pages/feixian/feixian";
import React, { useState } from 'react';

const IndexPage = () => {
  /** 获取设备统计 */
  const { data } = useRequest(postApiV1ThingsDeviceInfoCount, {
    onError: (error) => {
      message.error('获取设备统计错误:' + error.message);
    },
  });
// 创建一个状态来保存飞线数据
const [flyLineData, setFlyLineData] = useState([]);

// 创建一个方法来更新飞线数据
const addFlyLine = (newLine) => {
  // setFlyLineData([...flyLineData, newLine]);
  setFlyLineData((prevData) => [...prevData, newLine]);
};


  const deviceTotal = useMemo(() => {
    return (
      (data?.data.deviceCount.unknown || 0) +
      (data?.data.deviceCount.inactive || 0) +
      (data?.data.deviceCount.offline || 0) +
      (data?.data.deviceCount.online || 0) || 0
    );
  }, [data?.data]);

  return (
    <Card>
      <Row gutter={26}>
        <Col span={14}>
          {/*<DeviceMap />*/}
          <MyChart flyLineData={flyLineData} />
        </Col>
        <Col span={10}>
          <div style={{ height: '85vh' }}>
            <DeviceCount data={data?.data} deviceTotal={deviceTotal} />
            <Pie />
            <DeviceChart data={data?.data} addFlyLine={addFlyLine}/>
          </div>
        </Col>
      </Row>
    </Card>
  );
};

export default IndexPage;
