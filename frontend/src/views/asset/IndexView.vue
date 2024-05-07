<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { reactive, ref } from 'vue'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import router from '@/router/index.js'
import { toSQLLine } from '@/utils/stringFun.js'
import { FormatDate } from '@/utils/format.js'
import { AddAsset, DeleteAsset, FetchAsset, UpdateAsset } from '@/apis/asset.js'
import { SubmitSubdomainTask } from '@/apis/subdomain.js'
import { AddCron } from '@/apis/cron.js'
import { DeleteTask } from '@/apis/task.js'

defineOptions({
  name: 'PortScanIndex'
})

const radio = ref(1)
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
  domains: '',
  ips: ''
})

const addPortScanDialog = ref(false)
const addPortScanForm = ref(null)
const addPortScanFormData = ref({
  assetUUID: '',
  title: '',
  targets: '',
  dictType: '',
  ports: '',
  timeout: 30,
  threads: 200,
  checkAlive: true
})

const addSubdomainDialog = ref(false)
const addSubdomainForm = ref(null)
const addSubdomainFormData = ref({
  assetUUID: '',
  title: '',
  targets: '',
  dictType: '1',
  timeout: 30,
  threads: 200
})

const rules = reactive({
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
  checkAlive: [{ required: true, message: '', trigger: 'blur' }],
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

function initPortScanForm() {
  addPortScanForm.value = {
    assetUUID: '',
    title: '',
    targets: '',
    dictType: '',
    ports: '',
    timeout: 30,
    threads: 200,
    checkAlive: true
  }
}

function initSubdomainForm() {
  addPortScanForm.value = {
    title: '',
    targets: '',
    dictType: '1',
    timeout: 30,
    threads: 200
  }
}

function showAddAssetDialog() {
  dialogFlag.value = 'add'
  addAssetDialog.value = true
}

function showAddPortScanDialog(row) {
  updatePorts()
  addPortScanFormData.value.assetUUID = row.uuid
  addPortScanFormData.value.targets = row.ips
  addPortScanFormData.value.title = '资产 ' + row.title + ' IP 监控任务'
  addPortScanDialog.value = true
}

function showAddSubdomainDialog(row) {
  addSubdomainFormData.value.assetUUID = row.uuid
  addSubdomainFormData.value.targets = row.domains
  addSubdomainFormData.value.title = '资产 ' + row.title + ' 站点监控任务'
  addSubdomainDialog.value = true
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

function closePortScanDialog() {
  initPortScanForm()
  radio.value = 1
  addPortScanDialog.value = false
}

function closeSubdomainDialog() {
  initSubdomainForm()
  radio.value = 1
  addSubdomainDialog.value = false
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

async function submitAddPortScanForm() {
  // 访问 Form 实例
  if (!addPortScanForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await addPortScanForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交端口扫描任务，请稍候...',
      spinner: 'loading'
    })

    try {
      const response = await AddCron({
        taskType: 'PortScan',
        ...addPortScanFormData.value
      })
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '任务提交成功'
        })
        closeSubdomainDialog()
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

async function submitAddSubdomainForm() {
  // 访问 Form 实例
  if (!addSubdomainForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await addSubdomainForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交端口扫描任务，请稍候...',
      spinner: 'loading'
    })

    try {
      const response = await AddCron({
        taskType: 'Subdomain',
        ...addSubdomainFormData.value
      })
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '任务提交成功'
        })
        closeSubdomainDialog()
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

const options = [
  {
    text: '数据库端口',
    value: '1433,1521,3306,5432,6379,9200,11211,27017'
  },
  {
    text: '企业端口',
    value:
      '21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017,80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880'
  },
  {
    text: '高危端口',
    value: '21,22,23,53,80,443,8080,8000,139,445,3389,1521,3306,6379,7001,2375,27017,11211'
  },
  {
    text: '全端口',
    value: '1-65535'
  },
  {
    text: '自定义',
    value: ''
  }
]

// 更新端口输入框内容
function updatePorts() {
  const index = Number(radio.value) - 1
  if (index >= 0 && index < options.length) {
    addPortScanFormData.value.ports = options[index].value
    addPortScanFormData.value.dictType = String(radio.value)
  }
}

function updateDictType() {
  addSubdomainFormData.value.dictType = String(radio.value)
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
        <el-table-column label="创建者" min-width="150" prop="creator.username" />
        <el-table-column label="创建时间" min-width="200" sortable="custom" prop="CreatedAt">
          <template #default="scope">
            {{ FormatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="420" fixed="right">
          <template #default="scope">
            <el-button icon="Edit" @click="showUpdateAssetDialog(scope.row)">修改</el-button>
            <el-button icon="Location" @click="showAddSubdomainDialog(scope.row)"
            >添加站点监控
            </el-button>
            <el-button icon="MagicStick" @click="showAddPortScanDialog(scope.row)"
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
          <span v-if="dialogFlag === 'add'" class="text-lg">添加资产</span>
          <span v-if="dialogFlag === 'update'" class="text-lg">更新资产</span>
          <div>
            <el-button @click="closeAddAssetDialog">取 消</el-button>
            <el-button type="primary" @click="submitAddAssetForm">确 定</el-button>
          </div>
        </div>
      </template>
      <warning-bar v-if="dialogFlag === 'add'" title="新增资产" />
      <warning-bar v-if="dialogFlag === 'update'" title="更新资产" />
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

    <el-drawer
      v-model="addPortScanDialog"
      size="40%"
      :before-close="closePortScanDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">添加端口扫描监控任务</span>
          <div>
            <el-button @click="closePortScanDialog">取 消</el-button>
            <el-button type="primary" @click="submitAddPortScanForm">确 定</el-button>
          </div>
        </div>
      </template>
      <warning-bar title="新增端口扫描任务" />
      <el-form ref="addPortScanForm" :model="addPortScanFormData" :rules="rules" label-width="auto">
        <el-form-item label="任务标题:" prop="title">
          <el-input disabled v-model="addPortScanFormData.title" />
        </el-form-item>

        <el-form-item label="IP:" prop="targets">
          <el-input
            disabled
            type="textarea"
            rows="3"
            v-model="addPortScanFormData.targets"
            resize="none"
          />
        </el-form-item>

        <el-form-item label="预设端口:" prop="dictType">
          <el-radio-group v-model="radio" @change="updatePorts">
            <el-radio :value="1">数据库端口</el-radio>
            <el-radio :value="2">企业端口</el-radio>
            <el-radio :value="3">高危端口</el-radio>
            <el-radio :value="4">全端口</el-radio>
            <el-radio :value="5">自定义</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="端口列表:" prop="ports">
          <el-input type="textarea" rows="3" v-model="addPortScanFormData.ports" resize="none" />
        </el-form-item>
        <el-form-item label="存活探测:" prop="checkAlive">
          <el-checkbox v-model="addPortScanFormData.checkAlive" />
        </el-form-item>
        <el-form-item label="线程数量:" prop="threads">
          <el-input-number
            controls-position="right"
            v-model="addPortScanFormData.threads"
            :min="1"
            :max="30000"
          />
        </el-form-item>
        <el-form-item label="超时时长(s):" prop="timeout">
          <el-input-number
            controls-position="right"
            v-model="addPortScanFormData.timeout"
            :min="1"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-model="addSubdomainDialog"
      size="40%"
      :before-close="closeSubdomainDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">添加子域名扫描任务</span>
          <div>
            <el-button @click="closeSubdomainDialog">取 消</el-button>
            <el-button type="primary" @click="submitAddSubdomainForm">确 定</el-button>
          </div>
        </div>
      </template>
      <warning-bar title="新增子域名扫描监控任务" />
      <el-form
        ref="addSubdomainForm"
        :model="addSubdomainFormData"
        :rules="rules"
        label-width="auto"
      >
        <el-form-item label="任务标题" prop="title">
          <el-input disabled v-model="addSubdomainFormData.title" />
        </el-form-item>

        <el-form-item label="域名列表:" prop="targets">
          <el-input
            disabled
            type="textarea"
            rows="3"
            v-model="addSubdomainFormData.targets"
            resize="none"
          />
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
            v-model="addSubdomainFormData.threads"
            :min="1"
            :max="30000"
          />
        </el-form-item>
        <el-form-item label="超时时长(s):" prop="timeout">
          <el-input-number
            controls-position="right"
            v-model="addSubdomainFormData.timeout"
            :min="1"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<style lang="scss" scoped></style>
