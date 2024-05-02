<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import { DeletePortScanResult, FetchPortScanResult } from '@/apis/portscan.js'
import { toSQLLine } from '@/utils/stringFun.js'

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
    const response = await FetchPortScanResult({
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
    for (const row of selectedRows.value) {
      let loadingInstance = ElLoading.service({
        lock: true,
        fullscreen: true,
        text: '正在执行批量删除，请稍候...',
        spinner: 'loading'
      })
      try {
        const response = await DeletePortScanResult({ uuid: row.uuid })
        if (response.code === 200) {
          ElMessage({
            type: 'success',
            message: '删除成功'
          })
          await getTableData()
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
      }
    }
  })
}

function exportData() {

}
</script>

<template>
  <div>
    <warning-bar title="注：没有注释" />

    <div class="my-search-box">
      <el-form
        ref="searchForm"
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="IP 地址">
          <el-input clearable v-model="searchInfo.ip" placeholder="IP 地址" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input clearable v-model.number="searchInfo.port" placeholder="端口" />
        </el-form-item>
        <el-form-item label="开放服务">
          <el-input clearable v-model="searchInfo.service" placeholder="开放服务" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="my-table-box">
      <div class="my-btn-list">
        <el-button icon="Delete" :disabled="selectedRows.length===0" @click="deleteSelectedItems">
          批量删除
        </el-button>
        <el-button icon="Share" @click="exportData">
          导出所有数据
        </el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'ip', order: 'descending' }"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="IP 地址" min-width="150" sortable="custom" prop="ip" />
        <el-table-column align="left" label="端口" min-width="200" sortable="custom" prop="port" />
        <el-table-column align="left" label="开放服务" min-width="150" sortable="custom" prop="service" />
        <el-table-column align="left" label="开放状态" min-width="150" prop="open">
          <el-tag type="success" disable-transitions>Open</el-tag>
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

<style lang="scss" scoped>

</style>
