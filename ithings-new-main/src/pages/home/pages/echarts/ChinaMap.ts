// import worldJson from "@/assets/china.json";
// import {Component} from "react";
//
// class WorldMap extends Component {
//
//   componentDidMount() {
//     this.createMap();
//   }
//
//   createMap = () => {
//     constructor(props) {
//       super(props);
//       this.state = {
//         options: {},
//         ipOptions: {},
//         mapChart: null,
//       };
//     }
//     const options = {
//       radius: '100%',
//       tooltip: {
//         trigger: 'item',
//       },
//       legend: {
//         orient: 'horizontal', //图例的排列方向
//         // textStyle: { color: '#1a1e45' },
//         x: 'left', //图例的位置
//         y: '-20000000000000',
//       },
//       visualMap: {
//         //颜色的设置  dataRange
//         // textStyle: { color: '#1a1e45' },
//         x: 'left',
//         y: 'bottom',
//         // splitList: [{ start: 0, end: 150000 }],
//         show: false,
//         // text:['高','低'],// 文本，默认为数值文本
//         color: [flyColor], // 组件配置变量
//       },
//       geo: {
//         map: 'world', // 此处需与注册地图命名保持一致
//         type: 'map',
//         zoom: 1.2,
//         label: {
//           normal: {
//             show: false,
//             textStyle: {
//               color: '#FFFFFF',
//             },
//           },
//           emphasis: {
//             show: false,
//           },
//         },
//         roam: false, //是否允许缩放
//         itemStyle: {
//           normal: {
//             color: bgColor, //地图背景色，用到组件配置的变量
//             borderColor: borderColor, //省市边界线
//             borderWidth: 1,
//             textStyle: '#fff',
//           },
//           emphasis: {
//             areaColor: selectColor, //悬浮背景
//           },
//         },
//         data: [],
//       },
//       series: this.getSeries(centerPoint, mainData, flyLineArr, coordData, flyDirection),  // 飞线的配置
//     }
//
//     //飞线配置
//     getSeries = (centerPoint, mainData, flyLineArr, coordData, flyDirection) => {
//       let series = [];
//       let centerPointName = Object.keys(centerPoint)[0] || '北京区域中心';
//       let centerPointValue = Object.values(centerPoint)[0] || [116, 39];
//       [[centerPointName, flyLineArr]].forEach((item, i) => {
//         series.push(
//           {
//             type: 'lines',
//             coordinateSystem: 'geo',
//             zlevel: 2,
//             effect: {
//               show: true,
//               period: 5, //箭头指向速度，值越小速度越快
//               trailLength: 0, //特效尾迹长度[0,1]值越大，尾迹越长重
//               symbol: 'arrow', //箭头图标
//               symbolSize: 5, //图标大小
//               color: mainData.iconColor, // 图标颜色
//             },
//             lineStyle: {
//               normal: {
//                 show: true,
//                 width: 1, //尾迹线条宽度
//                 opacity: 1, //尾迹线条透明度
//                 curveness: 0.3, //尾迹线条曲直度
//                 color: mainData.flyColor, // 飞线颜色 - 细线
//               },
//             },
//             data: this.convertData(item[1], coordData, flyDirection, centerPointValue),
//           },
//           {
//             type: 'effectScatter',
//             radius: '100%',
//             coordinateSystem: 'geo',
//             zlevel: 2,
//             rippleEffect: { //涟漪特效
//               period: 4, //动画时间，值越小速度越快
//               brushType: 'stroke', //波纹绘制方式 stroke, fill
//               scale: 3, //波纹圆环最大限制，值越大波纹越大
//               color: mainData.rippleColor,
//             },
//             label: {
//               normal: {
//                 show: false,
//                 position: 'right', //显示位置
//                 offset: [5, 0], //偏移设置
//                 formatter: (params) => {
//                   return params.data.name //圆环显示文字
//                 },
//                 fontSize: 13,
//               },
//               emphasis: {
//                 show: false,
//               },
//             },
//             symbol: 'circle',
//             symbolSize: (val) => {
//               return 5 //圆环大小
//             },
//             itemStyle: {
//               normal: {
//                 show: true,
//                 // areaColor: mainData.pointColor,
//                 // color: mainData.pointColor,
//                 areaColor: "#ade9f4",
//                 color: "#ade9f4",
//               },
//               emphasis: {
//                 // areaColor: mainData.pointColor,
//                 areaColor: "#ade9f4",
//               },
//             },
//             data: item[1].map((dataItem) => {
//               return {
//                 //在这里定义你所要展示的数据
//                 name: dataItem[0].name,
//                 value: coordData[dataItem[0].name]?.concat([dataItem[0].value]),
//               }
//             }),
//           },
//           //中心点
//           {
//             type: 'effectScatter',
//             radius: '100%',
//             coordinateSystem: 'geo',
//             zlevel: 15,
//             rippleEffect: {
//               period: 4,
//               brushType: 'stroke',
//               scale: 4,
//               color: '#FFD246',
//             },
//             label: {
//               normal: {
//                 show: false,
//                 position: 'right',
//                 //offset:[5, 0],
//                 color: '#FFD246',
//                 formatter: '{b}',
//                 textStyle: {
//                   color: '#FFD246',
//                 },
//               },
//               emphasis: {
//                 show: false,
//                 color: '#FFD246',
//               },
//             },
//             symbol: 'circle',
//             symbolSize: 5,
//             itemStyle: {
//               color: '#FFD246',
//             },
//             data: [
//               {
//                 name: item[0],
//                 value: coordData[item[0]]?.concat([10]),
//               },
//             ],
//           }
//         )
//       })
//       return series
//     };
//
//
//
//     const dom = document.getElementById(componentConfig.id); // 保证每个地图唯一
//     var mapChart = echarts.init(dom); // 初始化
//     echarts.registerMap("world", worldJson); // 注册地图
//     this.setState({ options, ipOptions }); // 飞线、IP都存state
//     mapChart.setOption(displayMode === 0 ? options : ipOptions); // 显示方式切换;配置地图选项
//     this.setState({ mapChart })
//   };
// }
