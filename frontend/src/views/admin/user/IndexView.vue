<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { reactive, ref } from 'vue'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import { toSQLLine } from '@/utils/stringFun.js'
import {
  AddUserInfo,
  DeleteUserInfo,
  FetchUserInfo,
  ResetPassword,
  UpdateUserInfo
} from '@/apis/user.js'
import { FormatDate } from '@/utils/format.js'
import { useUserStore } from '@/stores/modules/user.js'
import { QuestionFilled } from '@element-plus/icons-vue'
import ShowImgIndex from '@/components/showimg/IndexView.vue'
import AvatarIndex from '@/components/avatar/IndexView.vue'

defineOptions({
  name: 'UserIndex'
})

const path = 'http://103.228.64.175:8081/'

const dialogFlag = ref('add')

const userStore = useUserStore()
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const addUserForm = ref(null)
const addUserFormData = ref({
  username: '',
  password: '',
  nickname: '',
  mail: '',
  avatar: '',
  authorityId: '2',
  enable: '2'
})

const rules = reactive({
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 30, message: '昵称长度在 2 到 30 个字符', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  authorityId: [{ required: true, message: '请选择用户角色', trigger: 'change' }],
  mail: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
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
    const response = await FetchUserInfo({
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
  addUserFormData.value = {
    username: '',
    password: '',
    nickname: '',
    mail: '',
    avatar: '',
    authorityId: '2',
    enable: '2'
  }
}

const addUserDialog = ref(false)

function showAddUserDialog() {
  dialogFlag.value = 'add'
  addUserDialog.value = true
}

const showUpdateUserDialog = (row) => {
  dialogFlag.value = 'update'
  addUserFormData.value = JSON.parse(JSON.stringify(row))
  addUserDialog.value = true
}

function closeAddUserDialog() {
  initForm()
  addUserDialog.value = false
}

async function submitAddUserForm() {
  // 访问 Form 实例
  if (!addUserForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await addUserForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交用户信息，请稍候...',
      spinner: 'loading'
    })

    try {
      let response
      if (dialogFlag.value === 'add') {
        response = await AddUserInfo(addUserFormData.value)
      } else {
        response = await UpdateUserInfo(addUserFormData.value)
      }
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '用户信息提交成功'
        })
        await getTableData()
        closeAddUserDialog()
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

async function deleteUserInfo(row) {
  if (userStore.userInfo.uuid === row.uuid) {
    ElMessage({
      type: 'error',
      message: '无法删除当前用户',
      showClose: true
    })
    return
  }
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const response = await DeleteUserInfo({ uuid: row.uuid })
    if (response.code === 200) {
      ElMessage.success('删除成功')
      await getTableData()
    } else {
      ElMessage({
        type: 'error',
        message: response.msg,
        showClose: true
      })
    }
  })
}

async function resetPassword(row) {
  ElMessageBox.confirm('确定重置密码吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const response = await ResetPassword({ uuid: row.uuid })
    if (response.code === 200) {
      ElMessage.success('重置成功，提醒用户查看邮箱')
      await getTableData()
    } else {
      ElMessage({
        type: 'error',
        message: response.msg,
        showClose: true
      })
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

const authorityIdOptions = [
  {
    value: '1',
    label: '系统管理员'
  },
  {
    value: '2',
    label: '普通用户'
  }
]

const enableOptions = [
  {
    value: '1',
    label: '启用'
  },
  {
    value: '2',
    label: '冻结'
  }
]

function formatAuthority(value) {
  switch (value) {
    case '1':
      return '系统管理员'
    case '2':
      return '普通用户'
    default:
      return '未知角色'
  }
}
</script>

<template>
  <div>
    <warning-bar title="注：重置密码会生成一串随机密码发送给用户的邮箱" />
    <div class="my-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="用户 UUID">
          <el-input v-model="searchInfo.uuid" placeholder="用户 UUID" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="searchInfo.mail" placeholder="邮箱" />
        </el-form-item>
        <el-form-item label="用户角色">
          <el-select v-model="searchInfo.authorityId" clearable placeholder="请选择">
            <el-option
              v-for="item in authorityIdOptions"
              :key="item.value"
              :label="`${item.label}(${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="启用状态">
          <el-select v-model="searchInfo.enable" clearable placeholder="请选择">
            <el-option
              v-for="item in enableOptions"
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
      <div class="my-btn-list">
        <el-button type="primary" icon="plus" @click="showAddUserDialog"> 添加用户</el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="handleSortChange"
        :default-sort="{ prop: 'createdAt', order: 'descending' }"
      >
        <el-table-column align="left" label="头像" min-width="100">
          <template #default="scope">
            <ShowImgIndex
              style="margin-top:4px"
              img-type="avatar"
              :img-src="scope.row.avatar"
            />
          </template>
        </el-table-column>
        <el-table-column fixed label="用户 UUID" min-width="250" sortable="custom" prop="uuid" />
        <el-table-column label="用户名" min-width="130" sortable="custom" prop="username" />
        <el-table-column label="昵称" min-width="130" sortable="custom" prop="nickname" />
        <el-table-column label="邮箱" min-width="150" sortable="custom" prop="mail" />
        <el-table-column label="用户角色" min-width="100" sortable="custom" prop="authorityId">
          <template #default="scope">
            {{ formatAuthority(scope.row.authorityId) }}
          </template>
        </el-table-column>
        <el-table-column label="启用状态" min-width="90" sortable="custom" prop="enable">
          <template #default="scope">
            <el-switch
              v-model="scope.row.enable"
              disabled
              inline-prompt
              :active-value="'1'"
              :inactive-value="'2'"
              active-text="是"
              inactive-text="否"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="150" sortable="custom" prop="createdAt">
          <template #default="scope">
            {{ FormatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" min-width="300">
          <template #default="scope">
            <el-button icon="Edit" @click="showUpdateUserDialog(scope.row)">编辑</el-button>
            <el-button type="danger" icon="Delete" @click="deleteUserInfo(scope.row)"
            >删除
            </el-button>
            <el-button type="warning" icon="RefreshRight" @click="resetPassword(scope.row)"
            >重置密码
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
      v-model="addUserDialog"
      size="40%"
      :before-close="closeAddUserDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span v-if="dialogFlag === 'add'" class="text-lg">添加用户</span>
          <span v-if="dialogFlag === 'update'" class="text-lg">更新用户</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button type="primary" @click="submitAddUserForm">确 定</el-button>
          </div>
        </div>
      </template>
      <warning-bar v-if="dialogFlag === 'add'" title="新增用户" />
      <warning-bar v-if="dialogFlag === 'update'" title="更新用户" />
      <el-form ref="addUserForm" :model="addUserFormData" :rules="rules" label-width="auto">
        <el-form-item label="昵称:" prop="nickname">
          <el-input v-model="addUserFormData.nickname" />
        </el-form-item>
        <el-form-item label="用户名:" prop="username">
          <el-input :disabled="dialogFlag === 'update'" v-model="addUserFormData.username" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="密码:" prop="password">
          <el-input type="password" show-password v-model="addUserFormData.password" />
        </el-form-item>
        <el-form-item label="用户角色:" prop="authorityId">
          <el-select v-model="addUserFormData.authorityId" placeholder="请选择">
            <el-option
              v-for="item in authorityIdOptions"
              :key="item.value"
              :label="`${item.label}(${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="用户邮箱:" prop="mail">
          <template #label
          >用户邮箱:
            <el-tooltip placement="right-end">
              <template #content>
                邮箱将用于接受重置密码邮件以及任务执行完成通知<br />
                请保证邮箱有效<br />
              </template>
              <el-icon>
                <QuestionFilled size="24" />
              </el-icon>
            </el-tooltip>
          </template>
          <el-input v-model="addUserFormData.mail" />
        </el-form-item>
        <el-form-item label="启用账户:">
          <el-switch
            v-model="addUserFormData.enable"
            inline-prompt
            :active-value="'1'"
            :inactive-value="'2'"
            active-text="是"
            inactive-text="否"
          />
        </el-form-item>
        <el-form-item label="用户头像:" label-width="80px">
          <div>
            <img
              v-if="addUserFormData.avatar"
              alt="头像"
              class="header-img-box"
              :src="path + addUserFormData.avatar"
              @click="addUserFormData.avatar=''"
            >
            <AvatarIndex v-else :target="addUserFormData" :target-key="'avatar'" />
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<style lang="scss">
.header-img-box {
  @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}
</style>
