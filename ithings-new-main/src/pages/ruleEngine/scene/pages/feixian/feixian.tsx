import chinaJson from "@/assets/china.json"
import ReactEcharts from "echarts-for-react"

import * as echarts from 'echarts/core';
import React, {Component } from 'react';

export default class FeiXian extends Component {
  //echartsInstance = null; // 初始化echarts实例为null
  echartsReactRef = React.createRef(); // 使用React.createRef()创建ref
  // componentDidMount() {
  //   echarts.registerMap("china", chinaJson); // 注册地图
  // }


  getOption = () => {
    echarts.registerMap("china", chinaJson); // 注册地图
    let option = {
      backgroundColor: "",
      geo: {
        map: "china",
        aspectScale: 0.75, //长宽比
        zoom: 1.20,
        roam: false,
        itemStyle: {
          normal: {
            areaColor: {
              type: "radial",
              x: 0.5,
              y: 0.5,
              r: 0.8,
              colorStops: [
                {
                  offset: 0,
                  color: "#09132c", // 0% 处的颜色
                },
                {
                  offset: 1,
                  color: "#274d68", // 100% 处的颜色
                },
              ],
              globalCoord: true, // 缺省为 false
            },
            shadowColor: "rgb(58,115,192)",
            shadowOffsetX: 10,
            shadowOffsetY: 11,
          },
          emphasis: {
            areaColor: "#2AB8FF",
            borderWidth: 0,
            color: "green",
            label: {
              show: false,
            },
          },
        },
        regions: [
          {
            name: "南海诸岛",
            itemStyle: {
              areaColor: "rgba(0, 10, 52, 1)",

              borderColor: "rgba(0, 10, 52, 1)",
              normal: {
                opacity: 0,
                label: {
                  show: false,
                  color: "#009cc9",
                },
              },
            },
          },
        ],
      },
      series: [
        {
          type: "map",
          roam: false,
          label: {
            normal: {
              show: true,
              textStyle: {
                color: "#1DE9B6",
              },
            },
            emphasis: {
              textStyle: {
                color: "rgb(183,185,14)",
              },
            },
          },

          itemStyle: {
            normal: {
              borderColor: "rgb(147, 235, 248)",
              borderWidth: 1,
              areaColor: {
                type: "radial",
                x: 0.5,
                y: 0.5,
                r: 0.8,
                colorStops: [
                  {
                    offset: 0,
                    color: "#09132c", // 0% 处的颜色
                  },
                  {
                    offset: 1,
                    color: "#274d68", // 100% 处的颜色
                  },
                ],
                globalCoord: true, // 缺省为 false
              },
            },
            emphasis: {
              areaColor: "rgb(46,229,206)",
              //    shadowColor: 'rgb(12,25,50)',
              borderWidth: 0.1,
            },
          },
          zoom: 1.25,
          //     roam: false,
          map: "china", //使用
          // data: this.difficultData //热力图数据   不同区域 不同的底色
        },
        {
          type: "effectScatter",
          coordinateSystem: "geo",
          showEffectOn: "render",
          zlevel: 1,
          rippleEffect: {
            period: 15,
            scale: 4,
            brushType: "fill",
          },
          hoverAnimation: true,
          label: {
            normal: {
              formatter: "{b}",
              position: "right",
              offset: [15, 0],
              color: "#1DE9B6",
              show: true,
            },
          },
          itemStyle: {
            normal: {
              color: "#1DE9B6",
              shadowBlur: 10,
              shadowColor: "#333",
            },
          },
          symbolSize: 12,
          data: [],
        }, //地图线的动画效果
        {
          type: "lines",
          zlevel: 2,
          effect: {
            show: true,
            period: 4, //箭头指向速度，值越小速度越快
            trailLength: 0.4, //特效尾迹长度[0,1]值越大，尾迹越长重
            symbol: "arrow", //箭头图标
            symbolSize: 7, //图标大小
          },
          lineStyle: {
            normal: {
              color: "#1DE9B6",
              width: 1, //线条宽度
              opacity: 0.1, //尾迹线条透明度
              curveness: 0.3, //尾迹线条曲直度
            },
          },
          data: [
            {
              coords: [
                [118.8062, 31.9208],
                [109.4543, 29.9222],
              ],
              lineStyle: {color: "#fffff8"},
            },
            {
              coords: [
                [127.9688, 45.368],
                [109.4543, 25.9222],
              ],
              lineStyle: {color: "#fffff8"},
            },
            {
              coords: [
                [110.3467, 41.4899],
                [117.4543, 31.9222],
              ],
              lineStyle: {color: "#fffff8"},
              /*#f34e2b */
            }
          ],
        },
      ],
    };
    return option;
  };
  componentDidUpdate(prevProps) {
    // 当传入的飞线数据更新时，重新设置图表的option
    if (prevProps.flyLineData !== this.props.flyLineData) {
      this.setOptionWithNewData(this.props.flyLineData);
    }
  }

  setOptionWithNewData = (newData) => {
    if (this.echartsInstance && this.echartsInstance.getEchartsInstance) {
      const echartsInstance = this.echartsInstance.getEchartsInstance();
      const option = this.getOption();
      //option.series[2].data = newData; 
      option.series[2].data = [...option.series[2].data, ...newData];
      //option.series[2].data.push(newData);
      echartsInstance.setOption(option);
    }
  };


  render() {
    return (
      <div>
        {/* <ReactEcharts option={this.getOption()} style={{
          marginTop: '40px',
          height: '80vh',
          width: '100%'
        }}/> */}
        <ReactEcharts
           ref={(e) => { this.echartsInstance = e; }}
           option={this.getOption()}
           style={{
             marginTop: '40px',
             height: '80vh',
             width: '100%'
           }}
         />
      </div>
    )
  }
}
