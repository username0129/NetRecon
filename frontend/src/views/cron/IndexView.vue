<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { CancelTask, DeleteTask, FetchTasks } from '@/apis/task.js'
import { ElMessage, ElMessageBox } from 'element-plus'
import router from '@/router/index.js'
import { toSQLLine } from '@/utils/stringFun.js'
import { FormatDate } from '@/utils/format.js'
import { useRoute } from 'vue-router'

defineOptions({
  name: 'PortScanIndex'
})

const route = useRoute()
const selectedRows = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  assetUUID: ref(route.query.uuid || null)
})

// 重置
async function onReset() {
  searchInfo.value = {}
  page.value = 1
  pageSize.value = 10
  await getTableData()
}

// 搜索
async function onSubmit() {
  page.value = 1
  pageSize.value = 10
  await getTableData()
}

async function handleCurrentChange(val) {
  page.value = val
  await getTableData()
}

async function handleSizeChange(val) {
  pageSize.value = val
  await getTableData()
}

async function getTableData() {
  try {
    const response = await FetchTasks({
      page: page.value,
      pageSize: pageSize.value,
      type: 'Cron',
      ...searchInfo.value
    })
    if (response.code === 200) {
      tableData.value = response.data.data
      total.value = response.data.total
      page.value = response.data.page
      pageSize.value = response.data.pageSize
    } else if (response.code === 404) {
      tableData.value = []
      total.value = 0
      page.value = 0
      pageSize.value = 0
      ElMessage({
        type: 'info',
        message: response.msg,
        showClose: true
      })
    } else {
      ElMessage({
        type: 'error',
        message: response.msg,
        showClose: true
      })
    }
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '网络错误或数据处理异常',
      showClose: true
    })
  }
}

// 页面加载时获取数据
getTableData()

async function deleteTask(row) {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const response = await DeleteTask({ uuid: row.uuid })
    if (response.code === 200) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}

async function cancelTask(row) {
  ElMessageBox.confirm('确定要取消吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const response = await CancelTask({ uuid: row.uuid })
    if (response.code === 200) {
      ElMessage.success('取消成功')
      await getTableData()
    }
  })
}

function formatStatus(value) {
  switch (value) {
    case '1':
      return '进行中'
    case '2':
      return '已完成'
    case '3':
      return '已取消'
    case '4':
      return '执行失败'
    default:
      return '未知状态'
  }
}

function formatDictType(type, value) {
  if (type === 'Cron/Port') {
    switch (value) {
      case '1':
        return '数据库端口'
      case '2':
        return '企业端口'
      case '3':
        return '高危端口'
      case '4':
        return '全端口'
      case '5':
        return '自定义'
      default:
        return '未知状态'
    }
  } else {
    switch (value) {
      case '1':
        return '小字典(5000)'
      case '2':
        return '中字典(20000)'
      case '3':
        return '大字典(110000)'
      default:
        return '未知'
    }
  }
}

function formatType(type) {
  if (type === 'Cron/Port') {
    return '端口扫描'
  } else {
    return '子域名爆破'
  }
}

function redirectToPortDetailPage(row) {
  // 跳转到任务详情页面
  router.push({
    name: 'PortscanDetail',
    query: {
      uuid: row.uuid
    }
  })
}

function redirectToDomainDetailPage(row) {
  // 跳转到任务详情页面
  router.push({
    name: 'SubdomainDetail',
    query: {
      uuid: row.uuid
    }
  })
}

function redirectToAssetPage(row) {
  // 跳转到任务详情页面
  router.push({
    name: 'Asset',
    query: {
      uuid: row.assetUUID
    }
  })
}

// 获取标签类型
function getTagType(status) {
  switch (status) {
    case '1':
      return 'info'
    case '2':
      return 'success'
    case '3':
      return 'warning'
    case '4':
      return 'danger'
    default:
      return ''
  }
}

async function handleSortChange({ prop, order }) {
  if (prop) {
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  await getTableData()
}

const statusOptions = [
  {
    value: '1',
    label: '进行中'
  },
  {
    value: '2',
    label: '已完成'
  },
  {
    value: '3',
    label: '已取消'
  },
  {
    value: '4',
    label: '执行失败'
  }
]

const typesOptions = [
  {
    value: 'Cron/Port',
    label: '端口扫描'
  },
  {
    value: 'Cron/Domain',
    label: '子域名爆破'
  }
]


// 监听选择项的变化
async function handleSelectionChange(selection) {
  selectedRows.value = selection
}
</script>

<template>
  <div>
    <warning-bar title="注：没有注释" />
    <div class="my-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="任务 UUID">
          <el-input v-model="searchInfo.uuid" placeholder="任务 UUID" />
        </el-form-item>
        <el-form-item label="资产 UUID">
          <el-input v-model="searchInfo.assetUUID" placeholder="任务 UUID" />
        </el-form-item>
        <el-form-item label="任务标题">
          <el-input v-model="searchInfo.title" placeholder="任务标题" />
        </el-form-item>
        <el-form-item label="任务目标">
          <el-input v-model="searchInfo.targets" placeholder="任务目标" />
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select v-model="searchInfo.type" clearable placeholder="请选择">
            <el-option
              v-for="item in typesOptions"
              :key="item.value"
              :label="`${item.label}(${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="执行状态">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="`${item.label}(${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="my-table-box">
      <div class="my-btn-list"></div>
      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'CreatedAt', order: 'descending' }"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column fixed label="任务 UUID" min-width="300" sortable="custom" prop="uuid">
          <template v-slot="scope">
            <a
              v-if="scope.row.type === 'Cron/Port'"
              href="#"
              @click="redirectToPortDetailPage(scope.row)"
              style="color: #00c5dc; text-decoration: none"
            >
              {{ scope.row.uuid }}
            </a>
            <a
              v-if="scope.row.type === 'Cron/Domain'"
              href="#"
              @click="redirectToDomainDetailPage(scope.row)"
              style="color: #00c5dc; text-decoration: none"
            >
              {{ scope.row.uuid }}
            </a>
          </template>
        </el-table-column>
        <el-table-column label="所属资产" min-width="300" sortable="custom" prop="assetUUID">
          <template v-slot="scope">
            <a
              href="#"
              @click="redirectToAssetPage(scope.row)"
              style="color: #00c5dc; text-decoration: none"
            >
              {{ scope.row.assetUUID }}
            </a>
          </template>
        </el-table-column>
        <el-table-column label="任务标题" min-width="200" sortable="custom" prop="title" />
        <el-table-column label="任务目标" min-width="150" sortable="custom" prop="targets" />
        <el-table-column label="任务类型" min-width="150" sortable="custom" prop="type">
          <template #default="scope">
            {{ formatType(scope.row.type) }}
          </template>
        </el-table-column>
        <el-table-column label="字典类型" min-width="120" sortable="custom" prop="dictType">
          <template #default="scope">
            {{ formatDictType(scope.row.type, scope.row.dictType) }}
          </template>
        </el-table-column>
        <el-table-column label="执行状态" min-width="130" sortable="custom" prop="status">
          <template #default="scope">
            <el-tag :type="getTagType(scope.row.status)" disable-transitions>
              {{ formatStatus(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建者" min-width="150" prop="creator.username" />
        <el-table-column label="上一次运行日期" min-width="200" sortable="custom" prop="lastTime" />
        <el-table-column label="下一次运行日期" min-width="200" sortable="custom" prop="nextTime" />
        <el-table-column label="创建时间" min-width="200" sortable="custom" prop="CreatedAt">
          <template #default="scope">
            {{ FormatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="200" fixed="right">
          <template #default="scope">
            <el-button
              :disabled="scope.row.status !== '1'"
              icon="Close"
              @click="cancelTask(scope.row)"
            >取消
            </el-button>
            <el-button
              type="danger"
              :disabled="scope.row.status === '1'"
              icon="Delete"
              @click="deleteTask(scope.row)"
            >删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="my-pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
