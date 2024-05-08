<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import { toSQLLine } from '@/utils/stringFun.js'
import {
  FetchSubdomainResult,
  ExportSubdomainResult,
  DeleteSubdomainResults
} from '@/apis/subdomain.js'

const route = useRoute()
const taskUUID = ref(route.query.uuid || '')
const selectedRows = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 获取表格数据
const getTableData = async () => {
  try {
    const response = await FetchSubdomainResult({
      page: page.value,
      pageSize: pageSize.value,
      taskUUID: taskUUID.value,
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

// 分页改变时更新数据
async function handleCurrentChange(val) {
  page.value = val
  await getTableData()
}

// 每页条数改变时更新数据
async function handleSizeChange(val) {
  pageSize.value = val
  await getTableData()
}

// 监听选择项的变化
async function handleSelectionChange(selection) {
  selectedRows.value = selection
}

// 排序
async function handleSortChange({ prop, order }) {
  if (prop) {
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  await getTableData()
}

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

async function deleteSelectedItems() {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const uuids = []
    selectedRows.value.forEach((item) => {
      uuids.push(item.uuid)
    })
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在执行批量删除，请稍候...',
      spinner: 'loading'
    })
    try {
      const response = await DeleteSubdomainResults({ uuids: uuids })
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '删除成功'
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
    } finally {
      loadingInstance.close()
      await getTableData()
    }
  })
}

async function exportData() {
  try {
    const response = await ExportSubdomainResult({ uuid: taskUUID.value })
    // 检查响应状态码是否为 200（OK）
    if (response.status === 200) {
      // 从响应头中提取文件名
      const contentDisposition = response.headers['content-disposition']
      let filename = 'default.csv' // 默认文件名
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename="(.+)*"/)
        if (filenameMatch) {
          filename = filenameMatch[1]
        }
      }
      // 创建一个 Blob URL，并使用临时 <a> 标签下载文件
      const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' })
      const downloadUrl = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = downloadUrl
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      // 显示成功消息
      ElMessage({
        type: 'success',
        message: '文件下载成功'
      })
    } else {
      // 如果响应状态码不是 200，处理错误
      ElMessage({
        type: 'error',
        message: '下载失败: 服务器处理异常',
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

function formatDate(value) {
  const date = new Date(value)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
}
</script>

<template>
  <div>
    <warning-bar title="注：没有注释" />

    <div class="my-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="域名">
          <el-input clearable v-model="searchInfo.subDomain" placeholder="域名" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input clearable v-model="searchInfo.title" placeholder="标题" />
        </el-form-item>
        <el-form-item label="IP 地址">
          <el-input clearable v-model="searchInfo.ips" placeholder="IP 地址" />
        </el-form-item>
        <el-form-item label="响应码">
          <el-input clearable v-model.number="searchInfo.code" placeholder="响应码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="my-table-box">
      <div class="my-btn-list">
        <el-button icon="Delete" :disabled="selectedRows.length === 0" @click="deleteSelectedItems">
          批量删除
        </el-button>
        <el-button :disabled="tableData.length === 0" icon="Share" @click="exportData">
          导出所有数据
        </el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'subDomain', order: 'descending' }"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column
          fixed
          align="left"
          label="子域名"
          min-width="150"
          sortable="custom"
          prop="subDomain"
        >
          <template v-slot="scope">
            <a
              :href="`http://${scope.row.subDomain}`"
              style="color: #00c5dc; text-decoration: none"
            >
              {{ scope.row.subDomain }}
            </a>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="响应状态"
          min-width="200"
          sortable="custom"
          prop="code"
        />
        <el-table-column
          align="left"
          label="网页标题"
          min-width="200"
          sortable="custom"
          prop="title"
        />
        <el-table-column align="left" label="CNAME 解析记录" min-width="200" prop="cname" />
        <el-table-column
          align="left"
          label="IP 解析记录"
          min-width="300"
          sortable="custom"
          prop="ips"
        />
        <el-table-column
          align="left"
          label="扫描时间"
          min-width="200"
          sortable="custom"
          prop="CreatedAt"
        >
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="备注" min-width="200" prop="notes" />
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
