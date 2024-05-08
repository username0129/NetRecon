<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import { toSQLLine } from '@/utils/stringFun.js'
import { FormatDate } from '@/utils/format.js'
import { useRoute } from 'vue-router'
import { DeleteOperationResults, FetchOperationResult } from '@/apis/operation.js'

defineOptions({
  name: 'OperationIndex'
})

const route = useRoute()
const selectedRows = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  uuid: ref(route.query.uuid || null)
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
    const response = await FetchOperationResult({
      page: page.value,
      pageSize: pageSize.value,
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
      const response = await DeleteOperationResults({ uuids: uuids })
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

const fmtBody = (value) => {
  try {
    return JSON.parse(value)
  } catch (err) {
    return value
  }
}

// 获取标签类型
function getTagType(value) {
  const firstDigit = value[0] // 获取响应码的第一个数字

  switch (firstDigit) {
    case '1':
      return 'info'
    case '2':
      return 'success'
    case '3':
      return 'warning'
    case '4':
      return 'danger'
    case '5':
      return 'danger'
    default:
      return 'info'
  }
}
</script>

<template>
  <div>
    <warning-bar title="注：没有注释" />
    <div class="my-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="来源 IP">
          <el-input v-model="searchInfo.ip" placeholder="来源 IP" />
        </el-form-item>
        <el-form-item label="请求方法">
          <el-input v-model="searchInfo.method" placeholder="请求方法" />
        </el-form-item>
        <el-form-item label="响应码">
          <el-input v-model="searchInfo.code" placeholder="响应码" />
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
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'CreatedAt', order: 'descending' }"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="用户" min-width="100" prop="user.username" />
        <el-table-column label="来源 IP" min-width="150" prop="ip" />
        <el-table-column label="请求方法" min-width="100" sortable="custom" prop="method" />
        <el-table-column label="请求路径" min-width="300" sortable="custom" prop="path" />
        <el-table-column label="响应码" min-width="100" sortable="custom" prop="code">
          <template #default="scope">
            <el-tag :type="getTagType(scope.row.code)" disable-transitions>
              {{ scope.row.code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="请求参数" min-width="100" prop="body">
          <template #default="scope">
            <el-popover
              v-if="scope.row.body"
              placement="bottom"
              width="400px"
              title="请求参数"
              trigger="click"
            >
              <pre>{{ fmtBody(scope.row.body) }}</pre>
              <template #reference>
                <el-button class="m-2">查看详情</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="响应数据" min-width="100" prop="resp">
          <template #default="scope">
            <el-popover
              v-if="scope.row.resp"
              placement="bottom"
              width="400px"
              title="响应数据"
              trigger="click"
            >
              <pre>{{ fmtBody(scope.row.resp) }}</pre>
              <template #reference>
                <el-button class="m-2">查看详情</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="请求日期" min-width="180" sortable="custom" prop="CreatedAt">
          <template #default="scope">
            {{ FormatDate(scope.row.CreatedAt) }}
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
