
<template>
  <div>
    <br/>
    <text>选基金: </text>
    <input id="startDate" type="text" placeholder="2025-01-01"/>
    <input id="endDate" type="text" placeholder="2025-02-01"/>
    <input id="code" type="text" placeholder="输入基金代码..."/>
    <button class="btn" @click="updateData">查看</button>
    <button class="btn" @click="syncData">同步最新数据</button>
    <button class="btn"  @click="updateConsoleLog">更新运行日志</button>
    <br/>
    <div id="chart" style="width: 100%; height: 400px;"></div>
    <textarea id="consoleLog" placeholder="运行日志" style="width: 80%; height: 200px;"></textarea>
  </div>
</template>

<script lang="ts" setup>

import * as echarts from 'echarts';
import { onMounted } from 'vue';
import { GetChartData, GetChartDataLabels, SyncStockHisData, SyncFundHisData, GetFundChartData, GetFundChartDataLabels, GetFundName, GetConsoleLogs} from '../../wailsjs/go/main/App'

async function updateConsoleLog() {
  const textarea = document.getElementById('consoleLog');
  textarea.value = await GetConsoleLogs();
}

async function syncData() {
    const inputElement = document.getElementById('code') as HTMLInputElement | null;
    const code = inputElement?.value || '';
    SyncFundHisData(code);
}

async function updateData() {
  const chart = echarts.init(document.getElementById('chart'));
  const inputElement = document.getElementById('code') as HTMLInputElement | null;
  const code = inputElement?.value || '';
  const startDateInput = document.getElementById('startDate') as HTMLInputElement | null;
  const endDateInput = document.getElementById('endDate') as HTMLInputElement | null;

  const startDate = startDateInput?.value || '';
  const endDate = endDateInput?.value || '';

  var xLabels = await GetFundChartDataLabels(code, startDate, endDate);
  var fundData = await GetFundChartData(code, startDate, endDate);
  var fundName = await GetFundName(code)

  var minValue = Math.min.apply(null, fundData.map(row => row[0]));

  chart.setOption({
    title: { text: fundName + '-净值数据', left: 'center' },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: xLabels
    },
    yAxis: { type: 'value', min: minValue},
    series: [{
      name: '数值',
      type: 'line',
      data: fundData.map(row => row[0]),
      smooth: true,
      itemStyle: { color: '#a3c6ed' },
      lineStyle: { color: '#a3c6ed' }
    }]
  })

  // const inputElement = document.getElementById('stockCode') as HTMLInputElement | null;
  // const stockCode = inputElement?.value || '';

  // const startDateInput = document.getElementById('startDate') as HTMLInputElement | null;
  // const endDateInput = document.getElementById('endDate') as HTMLInputElement | null;

  // const startDate = startDateInput?.value || '';
  // const endDate = endDateInput?.value || '';

  // var goData = await GetChartData(stockCode, startDate, endDate);
  // var xData = await GetChartDataLabels(stockCode, startDate, endDate);

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
}

onMounted(async () => {
  const chart = echarts.init(document.getElementById('chart'));
  var xLabels = await GetFundChartDataLabels('021030', '2025-05-01', '2025-05-30');
  var fundData = await GetFundChartData('021030', '2025-05-01', '2025-05-30');
  var fundName = await GetFundName('021030')
  var minValue = Math.min.apply(null, fundData.map(row => row[0]));

  chart.setOption({
    title: { text: fundName + '净值数据', left: 'center'},
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: xLabels
    },
    yAxis: { type: 'value', min: minValue},
    series: [{
      name: '数值',
      type: 'line',
      data: fundData.map(row => row[0]),
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