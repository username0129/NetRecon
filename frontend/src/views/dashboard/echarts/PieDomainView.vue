<template>
  <div class="pie-chart-box">
    <div class="pie-chart-title">域名资产数据</div>
    <div ref="echartPie" class="pie-chart"></div>
  </div>
</template>

<script setup>
import * as echarts from 'echarts'
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useWindowResize } from '@/hooks/use-windows-resize'
import { FetchDomainCount } from '@/apis/echarts.js'
import { ElMessage } from 'element-plus'

const echartPie = ref(null)
let pieChart = null



const fetchPortCount = async () => {
  try {
    const response = await FetchDomainCount()
    if (response.code === 200) {
      const formattedData = response.data.map((item) => ({
        name: item.target,
        value: item.count
      }))
      initChart(formattedData)
    }
  } catch (error) {
    console.log(error)
    ElMessage({
      type: 'error',
      message: '网络错误或数据处理异常',
      showClose: true
    })
  }
}

const initChart = (chartData) => {
  if (!echartPie.value) return
  pieChart = echarts.init(echartPie.value)
  pieChart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      data: chartData.map((item) => item.name)
    },
    series: [
      {
        name: '端口资产数据',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '20',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: chartData
      }
    ]
  })
}

useWindowResize(() => {
  if (pieChart) pieChart.resize()
})

onMounted(async () => {
  await nextTick()
  await fetchPortCount()
})

onUnmounted(() => {
  if (pieChart) {
    pieChart.dispose()
    pieChart = null
  }
})
</script>

<style lang="scss" scoped>
.pie-chart-box {
  .pie-chart {
    height: 400px;
    width: 100%;
    background-color: #fff;
  }

  .pie-chart-title {
    font-weight: 600;
    margin-bottom: 12px;
  }
}
</style>
