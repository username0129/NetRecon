<script setup>
import { reactive, ref } from 'vue'
import { ElLoading, ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { init } from '@/apis/init.js'
import router from '@/router/index.js'

const initForm = ref(null)

const initFormData = reactive({
  dbType: '',
  host: '',
  port: '',
  username: '',
  password: '',
  dbName: ''
})

const rules = reactive({
  dbType: [{ required: true, message: '请选择数据库类型', trigger: 'blur' }],
  host: [
    { required: true, message: '请输入数据库地址', trigger: 'blur' },
    {
      pattern:
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
      message: '数据库地址必须是有效的 IPv4 地址',
      trigger: 'blur'
    }
  ],
  port: [
    { required: true, message: '请输入数据库端口', trigger: 'blur' },
    { pattern: /^[0-9]+$/, message: '数据库端口必须为数字', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        const portNumber = parseInt(value, 10)
        if (portNumber >= 1 && portNumber <= 65535) {
          callback() // 如果通过校验，则不传递任何参数给callback
        } else {
          callback(new Error('数据库端口必须在1到65535之间')) // 如果校验失败，传递一个 Error 对象给 callback
        }
      },
      trigger: 'blur'
    }
  ],
  username: [{ required: true, message: '请输入数据库用户名', trigger: 'blur' }],
  dbName: [{ required: true, message: '请输入数据库名', trigger: 'blur' }]
})

const onChange = (val) => {
  switch (val) {
    case 'mysql':
      Object.assign(initFormData, {
        dbType: 'mysql',
        host: '127.0.0.1',
        port: '3306',
        username: 'root',
        password: '',
        dbName: 'nr'
      })
      break
    case 'pgsql':
      Object.assign(initFormData, {
        dbType: 'pgsql',
        host: '127.0.0.1',
        port: '5432',
        username: 'postgres',
        password: '',
        dbName: 'nr'
      })
      break
    default:
      Object.assign(initFormData, {
        dbType: '',
        host: '',
        port: '',
        username: '',
        password: '',
        dbName: '',
        dbPath: ''
      })
  }

  if (initForm.value) {
    initForm.value.validate()
  }
}

async function submitForm() {
  // 访问 Form 实例
  if (!initForm.value) {
    console.error('Form 实例未生效。')
    return
  }

  // 使用 Element Plus 的 validate 方法进行表单验证
  initForm.value.validate(async (valid) => {
    if (valid) {
      let loadingInstance = ElLoading.service({
        lock: true,
        fullscreen: true,
        text: '正在初始化数据库，请稍候...',
        spinner: 'loading'
      })

      try {
        const response = await init(initFormData)
        if (response.code === 200) {
          ElMessage({
            type: 'success',
            message: '数据库初始化成功'
          })
          await router.push({ name: 'Login' })
        } else {
          ElMessage({
            type: 'error',
            message: response.msg,
            showClose: true
          })
        }
      } catch (error) {
        // 捕获登录过程中的异常错误
        console.error('初始化过程中出现错误:', error)
        ElMessage({
          type: 'error',
          message: '初始化过程中出现异常错误',
          showClose: true
        })
      } finally {
        loadingInstance.close()
      }
    } else {
      // 表单验证失败，显示错误消息
      ElMessage({
        type: 'error',
        message: '请正确填写数据配置信息',
        showClose: true
      })
    }
  })
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
            ref="initForm"
            :model="initFormData"
            :rules="rules"
            :validate-on-rule-change="false"
          >
            <el-form-item prop="dbType" class="mb-6">
              <!-- mb -> margin-bottom，底部距离-->
              <el-select
                size="large"
                clearable
                v-model="initFormData.dbType"
                placeholder="请选择数据库类型"
                @change="onChange"
              >
                <template #prefix>
                  <el-icon>
                    <Search />
                  </el-icon>
                </template>

                <el-option key="mysql" label="mysql" value="mysql" />
                <el-option key="pgsql" label="pgsql" value="pgsql" />
              </el-select>
            </el-form-item>

            <el-form-item prop="host" class="mb-6">
              <el-input
                prefix-icon="lock"
                v-model="initFormData.host"
                size="large"
                placeholder="请输入数据库地址"
              />
            </el-form-item>

            <el-form-item prop="port" class="mb-6">
              <el-input
                prefix-icon="lock"
                v-model="initFormData.port"
                size="large"
                placeholder="请输入数据库端口"
              />
            </el-form-item>

            <el-form-item prop="username" class="mb-6">
              <el-input
                prefix-icon="user"
                v-model="initFormData.username"
                size="large"
                placeholder="请输入数据库用户名"
              />
            </el-form-item>

            <el-form-item prop="password" class="mb-6">
              <el-input
                prefix-icon="lock"
                show-password
                v-model="initFormData.password"
                size="large"
                placeholder="请输入数据库密码"
              />
            </el-form-item>

            <el-form-item prop="dbName" class="mb-6">
              <el-input
                prefix-icon="lock"
                v-model="initFormData.dbName"
                size="large"
                placeholder="请输入数据库名"
              />
            </el-form-item>
          </el-form>
          <el-form-item class="mb-6">
            <el-button
              class="shadow shadow-blue-600 h-11 w-full"
              type="primary"
              size="large"
              @click="submitForm"
              >初始化数据库
            </el-button>
          </el-form-item>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
