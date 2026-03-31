import React, { Component } from 'react';

class HtmlPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      timestamp: Date.now(), // 初始化时间戳
    };
  }

  componentDidMount() {
    this.interval = setInterval(() => {
      this.setState({ timestamp: Date.now() }); // 更新时间戳
    }, 20000); // 每20秒更新一次
  }

  componentWillUnmount() {
    clearInterval(this.interval); // 清除定时器
  }

  render() {
    const { timestamp } = this.state;
    const src = `/html/chaotu.html?t=${timestamp}`; // 添加时间戳作为查询参数

    return (
      <div>
        <iframe src={src} width={1000} height={600}></iframe>
      </div>
    );
  }
}

export default HtmlPage;

// import { Component } from 'react';


// export default class HtmlPage extends Component {

//   render() {

//     return (
//       <div>
//         <iframe src="/html/chaotu.html" width={1000} height={600}></iframe>
//       </div>
//     )
//   }
// }
{/* <div>
<iframe src="http://43.143.227.15/index.html" width={1500} height={600}></iframe>
</div> */}
