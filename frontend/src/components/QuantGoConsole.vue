
<template>
  <div>
    <br/>
    <text>选基金/股票: </text>
    <select id="daysSelect">
      <option value="30">30</option>
      <option value="60">60</option>
      <option value="90">90</option>
    </select>
    <input id="code" type="text" placeholder="输入基金代码..."/>
    <input id="smaWin" type="text" placeholder="输入SMA窗口..."/>
    <button class="btn" @click="updateData">执行策略</button>
    <button class="btn" @click="syncData">同步最新数据</button>
    <button class="btn"  @click="updateConsoleLog">更新运行日志</button>
    <button class="btn" @click="startAgentBot">启用交易机器人</button>
    <button class="btn" @click="stopAgentBot">停止</button>
    <br/>
    <button class="btn" @click="seeReturnRateDistribution">收益率分布</button>
    <button class="btn" @click="seeBacktestingAsset">回测资产曲线</button>
    <pre id="indicators" style="white-space: pre-line;"></pre>
    <div id="chart" style="width: 100%; height: 400px;"></div>
    <textarea id="consoleLog" placeholder="运行日志" style="width: 80%; height: 200px;"></textarea>
  </div>
</template>

<script lang="ts" setup>

import * as echarts from 'echarts';
import { onMounted } from 'vue';
import { GetChartData, GetChartDataLabels, SyncStockHisData, SyncFundHisData, GetFundChartData, GetFundChartDataLabels, GetFundName, GetConsoleLogs, GetFundChartDataLastDaysLabels, GetFundChartLastDaysData, SimpleMovingAverage, MovingAverageCrossover, CalPricesVaR, MaximumDrawdown, NormalDistribution, MovingAverageCrossoverProfit} from '../../wailsjs/go/main/App'

var timerId = 0;

async function updateConsoleLog() {
  const textarea = document.getElementById('consoleLog');
  textarea.value = await GetConsoleLogs();
}

async function syncData() {
    const inputElement = document.getElementById('code') as HTMLInputElement | null;
    const code = inputElement?.value || '';
    SyncFundHisData(code);
}

async function updateStockData() {
  const chart = echarts.init(document.getElementById('chart'));
  const inputElement = document.getElementById('code') as HTMLInputElement | null;
  const stockCode = inputElement?.value || '';

  const startDateInput = document.getElementById('startDate') as HTMLInputElement | null;
  const endDateInput = document.getElementById('endDate') as HTMLInputElement | null;

  const startDate = startDateInput?.value || '';
  const endDate = endDateInput?.value || '';

  var goData = await GetChartData(stockCode, startDate, endDate);
  var xData = await GetChartDataLabels(stockCode, startDate, endDate);

  const option = {
    title: { text: stockCode, left: 'center'},
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        },
        formatter: (params: any) => {
            if(Array.isArray(params)) {
                let tar = params[0];
                if(tar.seriesType === 'candlestick') {
                    return `<div style="border-bottom: 1px solid rgba(255,255,255,.3); font-size: 18px;padding-bottom: 7px;margin-bottom: 7px">
                        ${tar.name}</div>
                        开盘 : ${tar.value[1]}<br/>
                        收盘 : ${tar.value[2]}<br/>
                        最高 : ${tar.value[4]}<br/>
                        最低 : ${tar.value[3]}`;
                }
            }
            return '';
        }
    },
    xAxis: {
      type: 'category',
      data: xData 
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      data: goData,
      type: 'candlestick'
    }]
  };

  chart.setOption(option);
}

async function startAgentBot() {
  let count = 0;
  timerId = setInterval(() => {
    console.log(`第 ${count + 1} 次执行`);
    count++;
    updateConsoleLog();
    updateData();
    }, 10 * 1000);
}

async function stopAgentBot() {
   clearInterval(timerId);
}

async function seeReturnRateDistribution() {
  const chart = echarts.init(document.getElementById('chart'));
  chart.clear();
  const inputElement = document.getElementById('code') as HTMLInputElement | null;
  const code = inputElement?.value || '';
  const selectElement = document.getElementById("daysSelect");
  const selectedValue = Number(selectElement.value);
  var fundName = await GetFundName(code);

  var xLabels = await GetFundChartDataLastDaysLabels(code, selectedValue);
  var fundData = await GetFundChartLastDaysData(code, selectedValue);

  var distribution = await NormalDistribution(fundData.map(row => row[0]).reverse())

  let xDatas = [];
  let yDatas = [];
  for (let point in distribution) {
    xDatas.push(point.X);
    yDatas.push(point.Y);
  }

  let seriesData = distribution.map(item => [item.X, item.Y]);

  chart.setOption({
        title: {
            text: fundName + '日累积收益 CDF 图', left: 'center', padding: [0, 0, 20, 20]
        },
        tooltip: {
            trigger: 'axis',
            formatter: function(params) {
                return `值: ${params[0].value[0]}<br/>累积概率: ${params[0].value[1].toFixed(4)}`;
            }
        },
        xAxis: {
            type: 'value',
            name: '值'
        },
        yAxis: {
            type: 'value',
            name: '累积概率',
            min: 0,
            max: 1
        },
        series: [{
            type: 'line',
            data: seriesData,
            showSymbol: false,
            lineStyle: {
                color: '#5470c6'
            }
        }]
    })
}

async function seeBacktestingAsset() {
  const chart = echarts.init(document.getElementById('chart'));
  chart.clear();

  const inputElement = document.getElementById('code') as HTMLInputElement | null;
  const code = inputElement?.value || '';
  const selectElement = document.getElementById("daysSelect");
  const selectedValue = Number(selectElement.value);
  const smaWinElement = document.getElementById("smaWin");
  const smaWinValue = Number(smaWinElement.value);

  const smaWinValueMultiply3 = smaWinValue * 3;

  var xLabels = await GetFundChartDataLastDaysLabels(code, selectedValue);
  var fundData = await GetFundChartLastDaysData(code, selectedValue);
  var singals = await MovingAverageCrossover(fundData.map(row => row[0]).reverse(), smaWinValue, smaWinValueMultiply3);
  var profits = await MovingAverageCrossoverProfit(fundData.map(row => row[0]).reverse(), smaWinValue, smaWinValueMultiply3);
  var fundName = await GetFundName(code);
  // console.log(profits.map(profit => profit.assets))

  chart.setOption({
    title: { text: fundName + ' 回测预估资产变化（初始资金 10000）', left: 'center'},
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: xLabels.reverse()
    },
    yAxis: { type: 'value'},
    series: [{
      name: '数值',
      type: 'line',
      data: profits.map(profit => profit.assets),
      smooth: true,
      itemStyle: { color: '#a3c6ed' },
      lineStyle: { color: '#a3c6ed' }
    }]
  })
    
}

async function updateData() {
  const chart = echarts.init(document.getElementById('chart'));
  chart.clear();
  const inputElement = document.getElementById('code') as HTMLInputElement | null;
  const code = inputElement?.value || '';
  const selectElement = document.getElementById("daysSelect");
  const selectedValue = Number(selectElement.value);
  const smaWinElement = document.getElementById("smaWin");
  const smaWinValue = Number(smaWinElement.value);

  const smaWinValueMultiply3 = smaWinValue * 3;

  var xLabels = await GetFundChartDataLastDaysLabels(code, selectedValue);
  var fundData = await GetFundChartLastDaysData(code, selectedValue);
  var singals = await MovingAverageCrossover(fundData.map(row => row[0]).reverse(), smaWinValue, smaWinValueMultiply3);
  var profit = await MovingAverageCrossoverProfit(fundData.map(row => row[0]).reverse(), smaWinValue, smaWinValueMultiply3);
  var fundName = await GetFundName(code);
  var sam5 = await SimpleMovingAverage(fundData.map(row => row[0]).reverse(), smaWinValue);
  var sma15 = await SimpleMovingAverage(fundData.map(row => row[0]).reverse(), smaWinValueMultiply3);

  var minValue = Math.min.apply(null, fundData.map(row => row[0]));

  console.log(singals)

  let markPoints = [];
  for (let i = 0; i < singals.length; i++) {
      markPoints.push(
      { 
        name: 'BUY', 
        coord: [xLabels[singals[i].index], fundData.map(row => row[0])[singals[i].index]], 
        value: singals[i].label
      });
  }

  var priceVaR = await CalPricesVaR(fundData.map(row => row[0]).reverse(), 0.95)
  var priceDrawdown = await MaximumDrawdown(fundData.map(row => row[0]).reverse())

  document.getElementById('indicators').textContent = 'VaR(%5): ' + priceVaR + '  MaximumDrawdown: ' + priceDrawdown;

  console.log(SimpleMovingAverage(fundData.map(row => row[0]).reverse(), smaWinValue));
  console.log(fundData.map(row => row[0]).reverse());


  chart.setOption({
    title: { text: fundName + '-净值数据', left: 'center', padding: [20, 0, 0, 0]},
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: xLabels.reverse()
    },
    legend: {
      data: ['原始数据', 'SMA' + smaWinValue.toString(), 'SMA' + smaWinValueMultiply3.toString()]
    },
    yAxis: { type: 'value', min: minValue},
    series: [
      {
        name: '原始数据',
        type: 'line',
        data: fundData.map(row => row[0]).reverse(),
        smooth: true,
        itemStyle: { color: '#ff0000' },
        lineStyle: { color: '#ff0000' },
        markPoint: {
            data: markPoints
        }
      },
      {
        name: 'SMA' + smaWinValue.toString(),
        type: 'line',
        data: sam5,
        smooth: true,
        itemStyle: { color: '#0000ff' },
        lineStyle: { color: '#0000ff' }
      },
      {
        name: 'SMA' + smaWinValueMultiply3.toString(),
        type: 'line',
        data: sma15,
        smooth: true,
        itemStyle: { color: '#000000' },
        lineStyle: { color: '#000000' }
      } 
    ]
  })

}

onMounted(async () => {
  const chart = echarts.init(document.getElementById('chart'));

  var xLabels = await GetFundChartDataLastDaysLabels('021030', 60);
  var fundData = await GetFundChartLastDaysData('021030', 60);
  var fundName = await GetFundName('021030')
  var minValue = Math.min.apply(null, fundData.map(row => row[0]));
  const inputElement = document.getElementById('code') as HTMLInputElement | null;
  inputElement.value = '021030';
  const smaWinElement = document.getElementById("smaWin");
  smaWinElement.value = 3;
  const dayselectElement = document.getElementById("daysSelect");
  dayselectElement.value = 60;

  chart.setOption({
    title: { text: fundName + '净值数据', left: 'center'},
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: xLabels.reverse()
    },
    yAxis: { type: 'value', min: minValue},
    series: [{
      name: '数值',
      type: 'line',
      data: fundData.map(row => row[0]).reverse(),
      smooth: true,
      itemStyle: { color: '#a3c6ed' },
      lineStyle: { color: '#a3c6ed' }
    }]
  })

  // var goData = await GetChartData('AAPL', '2025-01-01', '2025-02-01');
  // var xData = await GetChartDataLabels('AAPL', '2025-01-01', '2025-02-01');

  // const option = {
  //   tooltip: {
  //       trigger: 'axis',
  //       axisPointer: {
  //           type: 'cross'
  //       },
  //       formatter: (params: any) => {
  //           if(Array.isArray(params)) {
  //               let tar = params[0];
  //               if(tar.seriesType === 'candlestick') {
  //                   return `<div style="border-bottom: 1px solid rgba(255,255,255,.3); font-size: 18px;padding-bottom: 7px;margin-bottom: 7px">
  //                       ${tar.name}</div>
  //                       开盘 : ${tar.value[1]}<br/>
  //                       收盘 : ${tar.value[2]}<br/>
  //                       最高 : ${tar.value[4]}<br/>
  //                       最低 : ${tar.value[3]}`;
  //               }
  //           }
  //           return '';
  //       }
  //   },
  //   xAxis: {
  //     type: 'category',
  //     data: xData
  //   },
  //   yAxis: {
  //     type: 'value'
  //   },
  //   series: [{
  //     data: goData,
  //     type: 'candlestick'
  //   }]
  // };

  // chart.setOption(option);
})


</script>