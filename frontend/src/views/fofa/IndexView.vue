<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { toSQLLine } from '@/utils/stringFun.js'
import { FofaExportData, FofaSearch } from '@/apis/fofa.js'

defineOptions({
  name: 'FofaIndex'
})

const selectedRows = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
async function onReset() {
  searchInfo.value = {}
  page.value = 1
  pageSize.value = 10
  total.value = 0
  tableData.value = []
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
    const response = await FofaSearch({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (response.code === 200) {
      if (response.data.total === 0) {
        ElMessage({
          type: 'info',
          message: '未查询到有效数据',
          showClose: true
        })
        return
      }
      tableData.value = response.data.results
      total.value = response.data.total
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

async function handleSortChange({ prop, order }) {
  if (prop) {
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  await getTableData()
}

// 监听选择项的变化
async function handleSelectionChange(selection) {
  selectedRows.value = selection
}

async function exportData() {
  ElMessageBox.confirm('最多导出 1w 条数据', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 发送请求，携带必要的数据
      const response = await FofaExportData({
        page: 1,
        pageSize: total.value,
        query: searchInfo.value.query
      })
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
      // 网络或其他错误处理
      ElMessage({
        type: 'error',
        message: '网络错误或数据处理异常',
        showClose: true
      })
    }
  })
}
</script>

<template>
  <div>
    <warning-bar title="注：点击资产 UUID可以跳转到资产下的计划任务" />
    <div class="my-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="FOFA 查询语句">
          <el-input
            style="min-width: 600px"
            v-model="searchInfo.query"
            placeholder="FOFA 查询语句"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="my-table-box">
      <div class="my-btn-list">
        <el-button :disabled="tableData.length === 0" icon="Share" @click="exportData">
          导出所有数据
        </el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'CreatedAt', order: 'descending' }"
      >
        <el-table-column fixed label="URL" min-width="350" prop="url" />
        <el-table-column label="标题" min-width="150" prop="title" />
        <el-table-column label="IP 地址" min-width="200" prop="ip" />
        <el-table-column label="端口" min-width="200" prop="port" />
        <el-table-column label="协议" min-width="150" prop="protocol" />
        <el-table-column label="备案信息" min-width="200" prop="icp" />
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
