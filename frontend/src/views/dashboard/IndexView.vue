<script setup>
import EchartsLine from '@/views/dashboard/dashboardCharts/echartsLine.vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/modules/user.js'
import { FetchTaskCount } from '@/apis/task.js'
import { ElMessage } from 'element-plus'
import { Avatar, Comment, Sort } from '@element-plus/icons-vue'
// import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({
  name: 'DashboardIndex'
})

const portTaskCount = ref(0)
const domainTaskCount = ref(0)
const cronTaskCount = ref(0)

const userStore = useUserStore()
const toolCards = ref([
  {
    label: 'IP 扫描',
    icon: 'monitor',
    name: 'PortScan',
    color: '#ff9c6e',
    bg: 'rgba(255, 156, 110,.3)'
  },
  {
    label: '子域名扫描',
    icon: 'setting',
    name: 'Subdomain',
    color: '#69c0ff',
    bg: 'rgba(105, 192, 255,.3)'
  },
  {
    label: '资产管理',
    icon: 'menu',
    name: 'Asset',
    color: '#b37feb',
    bg: 'rgba(179, 127, 235,.3)'
  },
  {
    label: '定时任务',
    icon: 'cpu',
    name: 'Cron',
    color: '#ffd666',
    bg: 'rgba(255, 214, 102,.3)'
  }
])

async function fetchTaskCount(value) {
  try {
    const response = await FetchTaskCount({
      taskType: value
    })
    if (response.code === 200) {
      if (value === 'PortScan') {
        portTaskCount.value = response.data.total
      }
      if (value === 'Subdomain') {
        domainTaskCount.value = response.data.total
      }
      if (value === 'Cron') {
        cronTaskCount.value = response.data.total
      }
    }
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '网络错误或数据处理异常',
      showClose: true
    })
  }
}

fetchTaskCount('PortScan')
fetchTaskCount('Subdomain')
fetchTaskCount('Cron')

const router = useRouter()

const toTarget = (name) => {
  router.push({ name })
}
</script>

<template>
  <div class="page">
    <div class="my-card-box">
      <div class="my-card my-top-card">
        <div class="my-top-card-left">
          <div class="my-top-card-left-title">Hello，{{ userStore.userInfo.nickname }}</div>
          <div class="my-top-card-left-dot">您的任务列表如下：</div>
          <el-row class="my-8 w-[500px]">
            <el-col :span="8" :xs="24" :sm="8">
              <div class="flex items-center">
                <el-icon class="dashboard-icon">
                  <Sort />
                </el-icon>
                子域名扫描任务 ({{ domainTaskCount }})
              </div>
            </el-col>
            <el-col :span="8" :xs="24" :sm="8">
              <div class="flex items-center">
                <el-icon class="dashboard-icon">
                  <Avatar />
                </el-icon>
                端口扫描任务 ({{ portTaskCount }})
              </div>
            </el-col>
            <el-col :span="8" :xs="24" :sm="8">
              <div class="flex items-center">
                <el-icon class="dashboard-icon">
                  <Comment />
                </el-icon>
                监控任务 ({{ cronTaskCount }})
              </div>
            </el-col>
          </el-row>
        </div>
        <img src="@/assets/dashboard.png" class="my-top-card-right" alt />
      </div>
    </div>
    <div class="my-card-box">
      <div class="my-card quick-entrance">
        <div class="my-card-title">快捷入口</div>
        <el-row :gutter="20">
          <el-col
            v-for="(card, key) in toolCards"
            :key="key"
            :span="4"
            :xs="8"
            class="quick-entrance-items"
            @click="toTarget(card.name)"
          >
            <div class="quick-entrance-item">
              <div class="quick-entrance-item-icon" :style="{ backgroundColor: card.bg }">
                <el-icon>
                  <component :is="card.icon" :style="{ color: card.color }" />
                </el-icon>
              </div>
              <p>{{ card.label }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <div class="my-card-box">
      <div class="my-card">
        <div class="my-card-title">数据统计</div>
        <div class="p-4">
          <echarts-line />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.page {
  @apply p-0;
  .my-card-box {
    @apply p-4;
    & + .my-card-box {
      @apply pt-0;
    }
  }

  .my-card {
    @apply box-border bg-white rounded h-auto px-6 py-8 overflow-hidden shadow-sm;
    .my-card-title {
      @apply pb-5 border-t-0 border-l-0 border-r-0 border-b border-solid border-gray-100;
    }
  }

  .my-top-card {
    @apply h-72 flex items-center justify-between text-gray-500;
    &-left {
      @apply h-full flex flex-col w-auto;
      &-title {
        @apply text-3xl text-gray-600;
      }

      &-dot {
        @apply mt-4 text-gray-600 text-lg;
      }

      &-item {
        + .my-top-card-left-item {
          margin-top: 24px;
        }

        margin-top: 14px;
      }
    }

    &-right {
      height: 600px;
      width: 600px;
      margin-top: 28px;
    }
  }

  :deep(.el-card__header) {
    @apply p-0  border-gray-200;
  }

  .card-header {
    @apply pb-5 border-b border-solid border-gray-200 border-t-0 border-l-0 border-r-0;
  }

  .quick-entrance-items {
    @apply flex items-center justify-center text-center text-gray-800;
    .quick-entrance-item {
      @apply px-8 py-6 flex items-center flex-col transition-all duration-100 ease-in-out rounded-lg cursor-pointer;
      &:hover {
        @apply shadow-lg;
      }

      &-icon {
        @apply flex items-center h-16 w-16 rounded-lg justify-center mx-0 my-auto text-2xl;
      }

      p {
        @apply mt-2.5;
      }
    }
  }
}

.dashboard-icon {
  @apply flex items-center text-xl mr-2 text-blue-400;
}
</style>
