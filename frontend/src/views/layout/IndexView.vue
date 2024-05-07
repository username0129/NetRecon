<script setup>
import Aside from '@/views/layout/aside/AsideIndex.vue'
import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { fmtTitle } from '@/utils/fmtRouterTitle.js'
import { useUserStore } from '@/stores/modules/user.js'
import { ArrowDown } from '@element-plus/icons-vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ElLoading, ElMessage } from 'element-plus'
import { UpdatePassword, UpdateUserInfo } from '@/apis/user.js'

defineOptions({
  name: 'LayoutIndex'
})

const route = useRoute()
const userStore = useUserStore()

const userDialog = ref(false)
const userForm = ref(null)
const userFormData = ref({
  username: '',
  password: '',
  nickname: '',
  mail: '',
  authorityId: ''
})

function initForm() {
  userFormData.value = {
    username: '',
    password: '',
    nickname: '',
    mail: '',
    authorityId: '2'
  }
}

const showUserDialog = () => {
  userFormData.value.nickname = userStore.userInfo.nickname
  userFormData.value.username = userStore.userInfo.username
  userFormData.value.mail = userStore.userInfo.mail
  userDialog.value = true
}

function closeUserDialog() {
  initForm()
  userDialog.value = false
}

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
    { min: 6, message: '最少6个字符', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '最少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请输入确认密码', trigger: 'blur' },
    { min: 6, message: '最少6个字符', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdFormData.value.newPassword) {
          callback(new Error('两次密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  mail: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
  ]
})

const matched = computed(() => route.matched)

function getRoleName() {
  switch (userStore.userInfo.valueOf().authorityId) {
    case '1':
      return '系统管理员'
    case '2':
      return '普通用户'
    default:
      return '未知角色'
  }
}

async function submitUserForm() {
  // 访问 Form 实例
  if (!userForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await userForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交用户信息，请稍候...',
      spinner: 'loading'
    })

    try {
      let response
      response = await UpdateUserInfo(userFormData.value)
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '用户信息提交成功'
        })
        closeUserDialog()
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

const pwdDialog = ref(false)
const pwdForm = ref(null)
const pwdFormData = ref({
  password: '',
  newPassword: '',
  confirmPassword: ''
})

function initPwdForm() {
  pwdFormData.value = {
    password: '',
    newPassword: '',
    confirmPassword: ''
  }
}

const showPwdDialog = () => {
  pwdDialog.value = true
}

function closePwdDialog() {
  initPwdForm()
  pwdDialog.value = false
  pwdForm.value.clearValidate()
}

const submitPwdForm = async () => {
  // 访问 Form 实例
  if (!pwdForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  const valid = await pwdForm.value.validate()
  if (valid) {
    let loadingInstance = ElLoading.service({
      lock: true,
      fullscreen: true,
      text: '正在提交用户信息，请稍候...',
      spinner: 'loading'
    })
    try {
      let response
      response = await UpdatePassword(pwdFormData.value)
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '用户信息提交成功'
        })
        closeUserDialog()
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
</script>

<template>
  <el-container class="layout-cont">
    <el-container>
      <el-aside class="main-cont my-aside" width="220px">
        <div class="logo-container">
          <div class="app-title">NetRecon</div>
        </div>
        <Aside class="aside" />
      </el-aside>

      <el-main class="main-cont">
        <transition :duration="{ enter: 800, leave: 100 }" mode="out-in" name="el-fade-in-linear">
          <div class="fixed top-0 box-border z-50" :style="{ width: `calc(100% - 220px)` }">
            <el-header class="header-cont">
              <el-row>
                <el-col>
                  <el-row class="p-0 h-full">
                    <el-col
                      :xs="2"
                      :lg="1"
                      :md="1"
                      :sm="1"
                      :xl="1"
                      class="z-50 flex items-center pl-3"
                    >
                    </el-col>

                    <el-col
                      :xs="10"
                      :lg="14"
                      :md="14"
                      :sm="9"
                      :xl="14"
                      :pull="1"
                      class="flex items-center"
                    >
                      <el-breadcrumb class="breadcrumb">
                        <el-breadcrumb-item
                          v-for="item in matched.slice(1, matched.length)"
                          :key="item.path"
                        >
                          {{ fmtTitle(item.meta.title, route) }}
                        </el-breadcrumb-item>
                      </el-breadcrumb>
                    </el-col>

                    <el-col
                      :xs="12"
                      :lg="9"
                      :md="9"
                      :sm="14"
                      :xl="9"
                      class="flex items-center justify-end"
                    >
                      <div class="mr-1.5 flex items-center">
                        <el-dropdown>
                          <div class="flex justify-center items-center h-full w-full">
                            <span class="cursor-pointer flex justify-center items-center">
                              <!--                              <CustomPic />-->
                              <span style="margin-left: 5px">{{
                                userStore.userInfo.nickname
                              }}</span>
                              <el-icon>
                                <arrow-down />
                              </el-icon>
                            </span>
                          </div>
                          <template #dropdown>
                            <el-dropdown-menu>
                              <el-dropdown-item>
                                <span class="font-bold">当前角色：{{ getRoleName() }}</span>
                              </el-dropdown-item>
                              <el-dropdown-item icon="user" @click="showUserDialog"
                                >个人信息
                              </el-dropdown-item>
                              <el-dropdown-item icon="lock" @click="showPwdDialog"
                                >更新密码
                              </el-dropdown-item>
                              <el-dropdown-item icon="reading-lamp" @click="userStore.logout"
                                >登 出
                              </el-dropdown-item>
                            </el-dropdown-menu>
                          </template>
                        </el-dropdown>
                      </div>
                    </el-col>
                  </el-row>
                </el-col>
              </el-row>
            </el-header>
            <el-drawer
              v-model="userDialog"
              size="40%"
              :before-close="closeUserDialog"
              :show-close="false"
            >
              <template #header>
                <div class="flex justify-between items-center">
                  <span class="text-lg">更新用户</span>
                  <div>
                    <el-button @click="closeUserDialog">取 消</el-button>
                    <el-button type="primary" @click="submitUserForm">确 定</el-button>
                  </div>
                </div>
              </template>
              <warning-bar title="更新用户" />
              <el-form ref="userForm" :model="userFormData" :rules="rules" label-width="auto">
                <el-form-item label="昵称:" prop="nickname">
                  <el-input v-model="userFormData.nickname" />
                </el-form-item>
                <el-form-item label="用户名:" prop="username">
                  <el-input disabled v-model="userFormData.username" />
                </el-form-item>
                <el-form-item label="用户邮箱:" prop="mail">
                  <el-input v-model="userFormData.mail" />
                </el-form-item>
              </el-form>
            </el-drawer>

            <el-drawer
              v-model="pwdDialog"
              size="40%"
              :before-close="closePwdDialog"
              :show-close="false"
            >
              <template #header>
                <div class="flex justify-between items-center">
                  <span class="text-lg">更新密码</span>
                  <div>
                    <el-button @click="closePwdDialog">取 消</el-button>
                    <el-button type="primary" @click="submitPwdForm">确 定</el-button>
                  </div>
                </div>
              </template>
              <warning-bar title="更新用户" />
              <el-form ref="pwdForm" :model="pwdFormData" :rules="rules" label-width="auto">
                <el-form-item label="旧密码:" prop="password">
                  <el-input type="password" v-model="pwdFormData.password" />
                </el-form-item>
                <el-form-item label="新密码:" prop="newPassword">
                  <el-input type="password" v-model="pwdFormData.newPassword" />
                </el-form-item>
                <el-form-item label="确认密码:" prop="confirmPassword">
                  <el-input type="password" v-model="pwdFormData.confirmPassword" />
                </el-form-item>
              </el-form>
            </el-drawer>
          </div>
        </transition>
        <router-view class="admin-box"></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<style lang="scss" scoped>
.logo-container {
  min-height: 60px;
  text-align: center;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #191a23;
}

.app-title {
  font-size: 24px;
  font-weight: bold;
  color: #ffffff;
}

.main-header {
  width: calc(100% - 220px);
  position: fixed;
  top: 0;
  z-index: 50;
  box-sizing: border-box;
}

.button {
  font-size: 12px;
  color: #666;
  background: rgb(250, 250, 250);
  width: 25px;
  padding: 4px 8px;
  border: 1px solid #eaeaea;
  margin-right: 4px;
  border-radius: 4px;
}

:deep .el-overlay {
  background-color: hsla(0, 0%, 100%, 0.9);
}
</style>
