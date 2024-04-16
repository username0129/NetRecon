<script setup>
import {captcha, login} from '@/apis/login.js'
import {reactive, ref} from 'vue'
import {ElLoading, ElMessage} from 'element-plus'
import {useUserStore} from '@/stores/modules/user.js'
import {useRouterStore} from '@/stores/modules/route.js'
import router from '@/router/index.js'
import {checkInit} from '@/apis/init.js'

const captchaImg = ref('') // 验证码图片 Base64
const loginForm = ref(null)
const loginFormData = reactive({
  username: '',
  password: '',
  answer: '',
  captchaId: '',
  openCaptcha: ''
})

const rules = reactive({
  username: [{required: true, message: '用户名不能为空', trigger: 'blur'}],
  password: [{required: true, message: '密码不能为空', trigger: 'blur'}],
  answer: [{required: false, message: '验证码不能为空', trigger: 'blur'}]
})

// 获取验证码
async function loginVerify() {
  try {
    const response = await captcha()

    if (response && response.code === 200 && response.data) {
      // 如果验证码获取成功，添加相应的验证规则
      rules.answer.push({
        required: true,
        message: '验证码不能为空',
        trigger: 'blur'
      })

      // 设置验证码图片和相关数据
      captchaImg.value = response.data.captchaImg
      loginFormData.captchaId = response.data.captchaId
      loginFormData.openCaptcha = response.data.openCaptcha
    } else {
      // 如果响应数据不符合预期，则输出错误信息
      console.error('Invalid response or missing data:', response)
    }
  } catch (error) {
    // 捕获并输出任何可能的错误
    console.error('Error fetching captcha:', error)
  }
}

loginVerify()

const userStore = useUserStore()
const routeStore = useRouterStore()

async function submitForm() {
  // 访问 Form 实例
  if (!loginForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  loginForm.value.validate(async (valid) => {
    if (valid) {
      let loadingInstance = ElLoading.service({
        lock: true,
        fullscreen: true,
        text: '正在登陆，请稍候...',
        spinner: 'loading'
      })

      try {
        const response = await login(loginFormData)
        if (response.code == 200) {
          userStore.setToken(response.data.token)
          userStore.setUserInfo(response.data.user)
          await routeStore.setRoutes()

          const routes = routeStore.routes
          routes.forEach(route => {
            router.addRoute(route)
          })

          await router.replace({"name": "dashboard"})
        } else {
          // 登录失败，例如：后端验证失败，密码错误等
          ElMessage({
            type: 'error',
            message: response.msg,
            showClose: true
          })
          await loginVerify()
        }
      } catch (error) {
        // 捕获登录过程中的异常错误
        console.error('登录过程中出现错误:', error)
        ElMessage({
          type: 'error',
          message: '登录过程中出现异常错误',
          showClose: true
        })
        await loginVerify()
      } finally {
        loadingInstance.close()
      }
    } else {
      // 表单验证失败，显示错误消息
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true
      })
      await loginVerify()
    }
  })
}

// 检查数据库初始化状态
async function check() {
  let loadingInstance = ElLoading.service({
    fullscreen: true,
    text: '检查中，请稍候...'
  })
  try {
    const response = await checkInit()
    // 系统初始化成功
    if (response.code == 200) {
      ElMessage({
        type: 'info',
        message: '系统已初始化成功',
        showClose: true
      })
    } else {
      // 系统未初始化，跳转到初始化页面
      userStore.setToken('')
      window.localStorage.removeItem('token')
      await router.push({name: 'Init'})
    }
  } catch (error) {
    // 捕获登录过程中的异常错误
    console.error('检查过程中出现错误:', error)
    ElMessage({
      type: 'error',
      message: '检查过程中出现异常错误',
      showClose: true
    })
  } finally {
    loadingInstance.close()
  }
}
</script>

<template>
  <div class="w-full h-full relative">
    <!--登陆背景-->
    <div class="bg-login-pattern bg-cover w-full h-full bg-center fixed top-0 left-0">
      <div class="w-full h-full flex rounded-lg items-center justify-evenly -mt-[10vh]">
        <div>
          <div class="flex items-center justify-center"></div>
          <div class="mb-9">
            <p class="text-center text-4xl font-bold">NetRecon</p>
            <p class="text-center text-1xl font-normal text-gray-500 mt-2.5">
              A management platform using Golang and Vue
            </p>
          </div>
          <el-form
              ref="loginForm"
              :model="loginFormData"
              :rules="rules"
              :validate-on-rule-change="false"
          >
            <el-form-item prop="username" class="mb-6">
              <!-- mb -> margin-bottom，底部距离-->
              <el-input
                  prefix-icon="user"
                  v-model="loginFormData.username"
                  size="large"
                  placeholder="请输入用户名"
              />
            </el-form-item>

            <el-form-item prop="password" class="mb-6">
              <el-input
                  prefix-icon="lock"
                  show-password
                  v-model="loginFormData.password"
                  size="large"
                  placeholder="请输入密码"
              />
            </el-form-item>

            <el-form-item v-if="loginFormData.openCaptcha" prop="answer" class="mb-6">
              <div class="flex w-full justify-between">
                <el-input
                    prefix-icon="filter"
                    v-model="loginFormData.answer"
                    placeholder="请输入验证码"
                    size="large"
                    class="flex-1 mr-5"
                />

                <div class="w-1/3 h-11 bg-[#c3d4f2] rounded">
                  <img
                      v-if="captchaImg"
                      class="w-full h-full"
                      :src="captchaImg"
                      alt="请输入验证码"
                      @click="loginVerify"
                  />
                </div>
              </div>
            </el-form-item>
          </el-form>

          <el-form-item class="mb-6">
            <el-button
                class="shadow shadow-blue-600 h-11 w-full"
                type="primary"
                size="large"
                @click="submitForm"
            >登 录
            </el-button>
          </el-form-item>

          <el-form-item class="mb-6">
            <el-button
                class="shadow shadow-blue-600 h-11 w-full"
                type="primary"
                size="large"
                @click="check"
            >检查系统初始化
            </el-button>
          </el-form-item>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
