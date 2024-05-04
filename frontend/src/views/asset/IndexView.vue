<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { reactive, ref } from 'vue'

import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import router from '@/router/index.js'
import { toSQLLine } from '@/utils/stringFun.js'
import { FormatDate } from '@/utils/format.js'
import { AddAsset, DeleteAsset, FetchAsset, UpdateAsset } from '@/apis/asset.js'
import { DeletePortScanResult } from '@/apis/portscan.js'


defineOptions({
  name: 'PortScanIndex'
})

const dialogFlag = ref('add')
const selectedRows = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const addAssetDialog = ref(false)
const addAssetForm = ref(null)
const addAssetFormData = ref({
  title: '',
  targets: '',
  dictType: '',
  ports: '',
  timeout: 30,
  threads: 200,
  checkAlive: true
})

const rules = reactive({
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }]
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
    const response = await FetchAsset({
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

function initForm() {
  addAssetFormData.value = {
    title: '',
    domains: '',
    ips: ''
  }
}


function showAddAssetDialog() {
  dialogFlag.value = 'add'
  addAssetDialog.value = true
}

const showUpdateAssetDialog = (row) => {
  dialogFlag.value = 'update'
  addAssetFormData.value = JSON.parse(JSON.stringify(row))
  addAssetDialog.value = true
}


function closeAddAssetDialog() {
  initForm()
  addAssetDialog.value = false
}

async function submitAddAssetForm() {
  // 访问 Form 实例
  if (!addAssetForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await addAssetForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交端口扫描任务，请稍候...',
      spinner: 'loading'
    })

    try {
      let response
      if (dialogFlag.value === 'add') {
        response = await AddAsset(addAssetFormData.value)
      } else {
        response = await UpdateAsset(addAssetFormData.value)
      }
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '任务提交成功'
        })
        await getTableData()
        closeAddAssetDialog()
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
  } else {
    // 表单验证失败，显示错误消息
    ElMessage({
      type: 'error',
      message: '请正确填写表单信息',
      showClose: true
    })
  }
}

async function deleteTask(row) {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const response = await DeleteAsset({ uuid: row.uuid })
    if (response.code === 200) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}


function redirectToDetailPage(row) {
  // 跳转到任务详情页面
  router.push({
    name: 'AssetDetail',
    query: {
      uuid: row.uuid
    }
  })
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
        const response = await DeleteAsset({ uuid: row.uuid })
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

</script>

<template>
  <div>
    <warning-bar title="注：没有注释" />
    <div class="my-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="资产 UUID">
          <el-input v-model="searchInfo.uuid" placeholder="资产 UUID" />
        </el-form-item>
        <el-form-item label="资产标题">
          <el-input v-model="searchInfo.title" placeholder="资产标题" />
        </el-form-item>
        <el-form-item label="域名">
          <el-input v-model="searchInfo.domains" placeholder="域名" />
        </el-form-item>
        <el-form-item label="IP 地址">
          <el-input v-model="searchInfo.ips" placeholder="IP 地址" />
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
        <el-button type="primary" icon="plus" @click="showAddAssetDialog">
          添加资产信息
        </el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'CreatedAt', order: 'descending' }"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column fixed label="资产 UUID" min-width="320" sortable="custom" prop="uuid">
          <template v-slot="scope">
            <a
              href="#"
              @click="redirectToDetailPage(scope.row)"
              style="color: #00c5dc; text-decoration: none"
            >
              {{ scope.row.uuid }}
            </a>
          </template>
        </el-table-column>
        <el-table-column label="资产标题" min-width="150" sortable="custom" prop="title" />
        <el-table-column label="域名列表" min-width="200" sortable="custom" prop="domains" />
        <el-table-column label="IP 列表" min-width="200" sortable="custom" prop="ips" />
        <el-table-column label="创建者" min-width="150" sortable="custom" prop="creator.username" />
        <el-table-column label="创建时间" min-width="200" sortable="custom" prop="CreatedAt">
          <template #default="scope">
            {{ FormatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="420" fixed="right">
          <template #default="scope">
            <el-button
              icon="Edit"
              @click="showUpdateAssetDialog(scope.row)"
            >修改
            </el-button>
            <el-button
              icon="Location"
              @click="deleteTask(scope.row)"
            >添加站点监控
            </el-button>
            <el-button
              icon="MagicStick"
              @click="deleteTask(scope.row)"
            >添加 IP 监控
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

    <el-drawer
      v-model="addAssetDialog"
      size="40%"
      :before-close="closeAddAssetDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span v-if="dialogFlag==='add'" class="text-lg">添加资产</span>
          <span v-if="dialogFlag==='update'" class="text-lg">更新资产</span>
          <div>
            <el-button @click="closeAddAssetDialog">取 消</el-button>
            <el-button type="primary" @click="submitAddAssetForm">确 定</el-button>
          </div>
        </div>
      </template>
      <warning-bar v-if="dialogFlag==='add'" title="新增资产" />
      <warning-bar v-if="dialogFlag==='update'" title="更新资产" />
      <el-form ref="addAssetForm" :model="addAssetFormData" :rules="rules" label-width="auto">
        <el-form-item label="资产标题:" prop="title">
          <el-input v-model="addAssetFormData.title" />
        </el-form-item>
        <el-form-item label="资产域名列表:" prop="domains">
          <el-input type="textarea" rows="3" resize="none" v-model="addAssetFormData.domains" />
        </el-form-item>
        <el-form-item label="资产 IP 列表:" prop="ips">
          <el-input type="textarea" rows="3" resize="none" v-model="addAssetFormData.ips" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<style lang="scss" scoped></style>
