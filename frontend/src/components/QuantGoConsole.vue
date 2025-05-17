
<template>
  <div>
    <div id="chart" style="width: 100%; height: 400px;"></div>
    <input id="startDate" type="text" placeholder="2025-01-01"/>
    <input id="endDate" type="text" placeholder="2025-02-01"/>
    <input id="stockCode" type="text" placeholder="Enter stock code..."/>
    <button class="btn" @click="updateData">查看</button>
    <button class="btn" @click="syncData">同步最新数据</button>
    <button class="btn"  @click="updateConsoleLog">更新运行日志</button>
    <br/>
    <br/>
    <textarea id="consoleLog" placeholder="运行日志" style="width: 80%; height: 200px;"></textarea>
  </div>
</template>

<script lang="ts" setup>

import * as echarts from 'echarts';
import { onMounted } from 'vue';
import { GetChartData, GetChartDataLabels, SyncStockHisData, GetConsoleLogs} from '../../wailsjs/go/main/App'

async function updateConsoleLog() {
  const textarea = document.getElementById('consoleLog');
  textarea.value = await GetConsoleLogs();
}

async function syncData() {
    const inputElement = document.getElementById('stockCode') as HTMLInputElement | null;
    const stockCode = inputElement?.value || '';
    SyncStockHisData(stockCode, 'YJM2Y8V6DA3XZ7UN')
}

async function updateData() {
  const chart = echarts.init(document.getElementById('chart'));

  const inputElement = document.getElementById('stockCode') as HTMLInputElement | null;
  const stockCode = inputElement?.value || '';

  const startDateInput = document.getElementById('startDate') as HTMLInputElement | null;
  const endDateInput = document.getElementById('endDate') as HTMLInputElement | null;

  const startDate = startDateInput?.value || '';
  const endDate = endDateInput?.value || '';

  var goData = await GetChartData(stockCode, startDate, endDate);
  var xData = await GetChartDataLabels(stockCode, startDate, endDate);

  const option = {
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

onMounted(async () => {
  const chart = echarts.init(document.getElementById('chart'));
  var goData = await GetChartData('AAPL', '2025-01-01', '2025-02-01');
  var xData = await GetChartDataLabels('AAPL', '2025-01-01', '2025-02-01');

  const option = {
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
})


</script>