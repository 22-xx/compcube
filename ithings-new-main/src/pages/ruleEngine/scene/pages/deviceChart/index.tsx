import { Col, Row } from 'antd';
import type { DeviceStatic } from '../../data';
import { Input } from 'antd';
import { Button, Modal } from 'antd';
import React, { useState } from 'react';

import { useHistory } from 'react-router-dom';


interface DeviceChartProps {
  data?: DeviceStatic;
}

const DeviceChart: React.FC<DeviceChartProps> = ({ data ,addFlyLine}) => {
  // const deviceTypeOptions = useMemo(() => {
  //   const deviceTypeData = [
  //     {
  //       name: '设备类型',
  //       value: data?.deviceTypeCount.device || 0,
  //     },
  //     {
  //       name: '网关类型',
  //       value: data?.deviceTypeCount.gateway || 0,
  //     },
  //     {
  //       name: '子设备类型数量',
  //       value: data?.deviceTypeCount.subset || 0,
  //     },
  //   ];
  //   return getChartsOption(deviceTypeData, [
  //     'rgba(0, 136, 255, 0.8)',
  //     'rgba(0, 179, 84, 0.8)',
  //     'rgba(255, 183, 0, 0.8)',
  //   ]);
  // }, [data?.deviceTypeCount]);

  // const deviceOnlineOptions = useMemo(() => {
  //   const deviceOnlineData = [
  //     {
  //       name: '在线设备',
  //       value: data?.deviceCount.online || 0,
  //     },
  //     {
  //       name: '离线设备',
  //       value: data?.deviceCount.offline || 0,
  //     },
  //   ];
  //   return getChartsOption(deviceOnlineData, ['#33C35E', '#AAB3B3']);
  // }, [data?.deviceCount]);

  const [isModalOpen, setIsModalOpen] = useState(false);

  
  const [setShowModal] = useState(false);

  const showModal = () => {
    setIsModalOpen(true);
  };
  // const generateNewData = () => {
  //   // 根据需要生成新的数据
  //   // ...
  //   const newData = {
  //     coords: [
  //       [106.3467, 41.4899],
  //       [107.4543, 31.9222],
  //     ],
  //     lineStyle: {color: "#f34e2b"},
  //   }
  //   return newData;
  // };

//   const handleOk = () => {
//     //setIsModalOpen(false);
//     // 创建新的飞线数据
//     const newFlyLine = {
//       // 飞线数据
//       coords: [[106.44744, 29.556489], [118.873019,34.323843]],
//       lineStyle: { color: "#f34e2b" },
//   };
  
//   // 调用传入的addFlyLine方法来更新飞线数据
//   addFlyLine(newFlyLine);

//   // 关闭模态框
//   setIsModalOpen(false);
// };
const handleOk = () => {
  // 创建新的飞线数据
  const newFlyLine = {
    // 飞线数据
    coords: [generateRandomCoord(), generateRandomCoord()],
    lineStyle: { color: "#f34e2b" },
  };
  
  // 调用传入的addFlyLine方法来更新飞线数据
  addFlyLine(newFlyLine);

  // 关闭模态框
  setIsModalOpen(false);
};

// 生成中国地图范围内的随机坐标
const generateRandomCoord = () => {
  const lng = Math.random() * (115.05 - 83.40) + 83.40; // 经度
  const lat = Math.random() * (41.33 - 28.51) + 28.51; // 纬度
  return [lng.toFixed(6), lat.toFixed(6)]; // 保留6位小数
};


  const handleCancel = () => {
    setIsModalOpen(false);
  };
  /*更改对应文字 */
  const [displayText, setDisplayText] = useState('');

  const handleButtonClick = (text: React.SetStateAction<string>) => {
    setDisplayText(text);
  };

  function generateRandomFloat(min: number, max: number) {
    return (Math.random() * (max - min) + min).toFixed(3);
  }

// 生成两个随机数
  const randomNum1 = generateRandomFloat(100, 999);
  const randomNum2 = generateRandomFloat(10, 99);

// 计算乘积结果
  const ZhiLiangYouXianRandomNum = "质量优先价格：" + randomNum1 + " x " + randomNum2 + "元/月 = " + (randomNum1 * randomNum2).toFixed(2) + "元/月";
  const JiaGeYouXianRandomNum = "价格优先价格：" + randomNum1 + " x " + randomNum2 + "元/月 = " + (randomNum1 * randomNum2).toFixed(2) + "元/月";
  const JunHengYouXianRandomNum = "均衡调度价格：" + randomNum1 + " x " + randomNum2 + "元/月 = " + (randomNum1 * randomNum2).toFixed(2) + "元/月";
//链接跳转
  //const history = useHistory();

  const handleClick = () => {
    //window.location.href = 'https://www.viesc.com/ymir.html';
    window.location.href = 'https://cn.bing.com/?mkt=zh-cn';
  };

  return (
    <Row gutter={20}>
      <Col span={24}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
          <div>
            <span style={{ marginRight: '20px' }}>CPU需求</span>
            <Input style={{ width: '200px' }} placeholder="请输入cpu需求" />
          </div>
          <div>
            <span style={{ marginRight: '20px' }}>带宽需求</span>
            <Input style={{ width: '200px' }} placeholder="请输入带宽需求" />
          </div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
          <div>
            <span style={{ marginRight: '20px' }}>内存需求</span>
            <Input style={{ width: '200px' }} placeholder="请输入内存需求" />
          </div>
          <div>
            <span style={{ marginRight: '20px' }}>时延需求</span>
            <Input style={{ width: '200px' }} placeholder="请输入时延需求" />
          </div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
          <div>
            <span style={{ marginRight: '20px' }}>磁盘需求</span>
            <Input style={{ width: '200px' }} placeholder="请输入磁盘需求" />
          </div>
          <div>
            <span style={{ marginRight: '20px' }}>GPU需求</span>
            <Input style={{ width: '200px' }} placeholder="请输入GPU需求" />
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'center', marginTop: '20px' , marginRight: '10px' }}>
        <div style={{ marginRight: '10px' }}>
          <>
            <Button type="primary" onClick={showModal}>
              节点匹配
            </Button>

            <Modal title="算法选择" open={isModalOpen} onOk={handleOk} onCancel={handleCancel}>
              <div style={{display:'flex', justifyContent:'space-around'}}>
                <Button style={{width:'100px', height:'100px'}} onClick={() => handleButtonClick('用户的请求被赋予更高优先级调度到性能出色的机器上。服务处理的质量较高，但所需要的价格较高，适合延迟敏感型服务。\n' + ZhiLiangYouXianRandomNum)}>质量优先</Button>
                <Button style={{width:'100px', height:'100px'}} onClick={() => handleButtonClick('用户的请求将被赋予较低优先级调度到性能较低的机器上。服务处理的质量较低，但价格便宜，适合离线批处理服务。\n' + JiaGeYouXianRandomNum)}>价格优先</Button>
                <Button style={{width:'100px', height:'100px'}} onClick={() => handleButtonClick('用户的请求将以一般的优先级调度到性能普通的机器上。服务的质量和价格都相对适中。\n' + JunHengYouXianRandomNum)}>均衡调度</Button>
              </div>
              <div style={{display:'flex', justifyContent:'space-around', flexDirection:'column', marginTop:'0px'}}>
                <br/><br/>
                <textarea value={displayText} style={{ display: 'flex', justifyContent: 'space-around', padding: '25px', borderRadius: '5px', border: '1px solid #ccc', wordWrap: 'break-word' }} readOnly />
              </div>
            </Modal>
            
          </>
          
        </div>
        <div style={{ marginRight: '10px' }}> 
        <>
          <Button type="primary" onClick={handleClick}>
              跳转
            </Button>
          </>
          </div>
          </div>
      </Col>
    </Row>
  );
};






export default DeviceChart;
