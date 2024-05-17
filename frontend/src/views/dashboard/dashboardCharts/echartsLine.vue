<template>
  <div class="dashboard-line-box">
    <div class="dashboard-line-title">每日任务</div>
    <div ref="echart" class="dashboard-line" />
  </div>
</template>

<script setup>
import * as echarts from 'echarts'
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useWindowResize } from '@/hooks/use-windows-resize'
import { FetchTaskCount } from '@/apis/echarts.js'
import { ElMessage } from 'element-plus'

const echart = ref(null)
let chart = null
const dataAxis = ref([])
const data = ref([])

const fetchTaskCount = async () => {
  try {
    const response = await FetchTaskCount()
    if (response.code === 200) {
      dataAxis.value = response.data.map(item => item.date)
      data.value = response.data.map(item => item.count)
      initChart()
    }
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '网络错误或数据处理异常',
      showClose: true
    })
  }
}


const initChart = () => {
  if (!echart.value) return
  chart = echarts.init(echart.value)
  chart.setOption({
    grid: { left: '40', right: '20', top: '40', bottom: '20' },
    xAxis: {
      data: dataAxis.value,
      axisTick: { show: false },
      axisLine: { show: false },
      z: 10
    },
    yAxis: {
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { textStyle: { color: '#999' } }
    },
    dataZoom: [{ type: 'inside' }],
    series: [{
      type: 'bar',
      barWidth: '40%',
      itemStyle: { borderRadius: [5, 5, 0, 0], color: '#188df0' },
      emphasis: { itemStyle: { color: '#188df0' } },
      data: data.value
    }]
  })
}


useWindowResize(() => {
  if (chart) chart.resize()
})

onMounted(async () => {
  await nextTick()
  await fetchTaskCount()
})

onUnmounted(() => {
  if (chart) {
    chart.dispose()
    chart = null
  }
})
</script>

<style lang="scss" scoped>
.dashboard-line-box {
  .dashboard-line {
    background-color: #fff;
    height: 360px;
    width: 100%;
  }

  .dashboard-line-title {
    font-weight: 600;
    margin-bottom: 12px;
  }
}
</style>
