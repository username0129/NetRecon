<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { reactive, ref } from 'vue'
import { CancelTask, DeleteTask, FetchTasks } from '@/apis/task.js'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import router from '@/router/index.js'
import { toSQLLine } from '@/utils/stringFun.js'
import { SubmitSubdomainTask } from '@/apis/subdomain.js'

defineOptions({
  name: 'BruteSubdomainIndex'
})

const selectedRows = ref([])
const radio = ref(1)
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const addTaskForm = ref(null)
const addTaskFormData = ref({
  title: '',
  targets: '',
  dictType: '1',
  timeout: 30,
  threads: 200
})

const rules = reactive({
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
  targets: [{ required: true, message: '目标不能为空', trigger: 'blur' }],
  timeout: [
    { required: true, message: '超时时间不能为空', trigger: 'blur' },
    { pattern: /^[0-9]+$/, message: '超时时间必须为整数', trigger: 'blur' }
  ],
  threads: [
    { required: true, message: '线程数不能为空', trigger: 'blur' },
    { pattern: /^[0-9]+$/, message: '线程数必须为整数', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        const portNumber = parseInt(value, 10)
        if (portNumber >= 1 && portNumber <= 2000) {
          callback() // 如果通过校验，则不传递任何参数给callback
        } else {
          callback(new Error('线程数必须在 1 到 2000 之间')) // 如果校验失败，传递一个 Error 对象给 callback
        }
      },
      trigger: 'blur'
    }
  ]
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
      type: 'BruteSubdomain',
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
  addTaskFormData.value = {
    title: '',
    targets: '',
    dictType: '1',
    timeout: 30,
    threads: 200
  }
}

const addTaskDialog = ref(false)

function showAddTaskDialog() {
  addTaskDialog.value = true
}

function closeAddTaskDialog() {
  initForm()
  addTaskDialog.value = false
}

async function submitAddTaskForm() {
  // 访问 Form 实例
  if (!addTaskForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await addTaskForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交端口扫描任务，请稍候...',
      spinner: 'loading'
    })

    try {
      const response = await SubmitSubdomainTask(addTaskFormData.value)
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '任务提交成功'
        })
        await getTableData()
        closeAddTaskDialog()
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

function formatDictType(value) {
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

function formatDate(value) {
  const date = new Date(value)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
}

function redirectToDetailPage(row) {
  // 跳转到任务详情页面
  router.push({
    name: 'SubdomainDetail',
    query: {
      uuid: row.uuid
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

function updateDictType() {
  addTaskFormData.value.dictType = String(radio.value)
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
        if (row.status !== '1') {
          const response = await DeleteTask({ uuid: row.uuid })
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
    await getTableData()
  })
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
        <el-form-item label="任务标题">
          <el-input v-model="searchInfo.title" placeholder="任务标题" />
        </el-form-item>
        <el-form-item label="任务目标">
          <el-input v-model="searchInfo.targets" placeholder="任务目标" />
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
        <el-button type="primary" icon="plus" @click="showAddTaskDialog">
          添加子域名扫描任务
        </el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
        :default-sort="{ prop: 'CreatedAt', order: 'descending' }"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column fixed label="任务 UUID" min-width="250" sortable="custom" prop="uuid">
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
        <el-table-column label="任务标题" min-width="150" sortable="custom" prop="title" />
        <el-table-column label="任务目标" min-width="150" sortable="custom" prop="targets" />
        <el-table-column label="字典类型" min-width="120" sortable="custom" prop="dictType">
          <template #default="scope">
            {{ formatDictType(scope.row.dictType) }}
          </template>
        </el-table-column>
        <el-table-column label="执行状态" min-width="100" sortable="custom" prop="status">
          <template #default="scope">
            <el-tag :type="getTagType(scope.row.status)" disable-transitions>
              {{ formatStatus(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建者" min-width="150" sortable="custom" prop="creator.username" />
        <el-table-column label="创建时间" min-width="150" sortable="custom" prop="CreatedAt">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
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

    <el-drawer
      v-model="addTaskDialog"
      size="40%"
      :before-close="closeAddTaskDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">添加子域名扫描任务</span>
          <div>
            <el-button @click="closeAddTaskDialog">取 消</el-button>
            <el-button type="primary" @click="submitAddTaskForm">确 定</el-button>
          </div>
        </div>
      </template>
      <warning-bar title="新增子域名扫描任务" />
      <el-form ref="addTaskForm" :model="addTaskFormData" :rules="rules" label-width="auto">
        <el-form-item label="任务标题" prop="title">
          <el-input v-model="addTaskFormData.title" />
        </el-form-item>

        <el-form-item label="域名列表:" prop="targets">
          <el-input type="textarea" rows="3" v-model="addTaskFormData.targets" resize="none" />
        </el-form-item>
        <el-form-item label="预设字典:" prop="dictType">
          <el-radio-group v-model="radio" @change="updateDictType">
            <el-radio :value="1">小字典(5000)</el-radio>
            <el-radio :value="2">中字典(20000)</el-radio>
            <el-radio :value="3">大字典(110000)</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="线程数量:" prop="threads">
          <el-input-number
            controls-position="right"
            v-model="addTaskFormData.threads"
            :min="1"
            :max="30000"
          />
        </el-form-item>
        <el-form-item label="超时时长(s):" prop="timeout">
          <el-input-number controls-position="right" v-model="addTaskFormData.timeout" :min="1" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<style lang="scss" scoped></style>
