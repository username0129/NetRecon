<template>
  <div class="w-full h-full relative">
    <!--登陆背景-->
    <div class="bg-login-pattern bg-cover w-full h-full bg-center fixed top-0 left-0">
      <div class="w-full h-full flex rounded-lg items-center justify-evenly -mt-[10vh]">
        <div>
          <div class="flex items-center justify-center">

          </div>
          <div class="mb-9">
            <p class="text-center text-4xl font-bold">NetRecon</p>
            <p class="text-center text-1xl font-normal text-gray-500 mt-2.5">A management platform using Golang and
              Vue</p>
          </div>
          <el-form
              ref="loginForm"
              :model="loginFormData"
              :rules="rules"
              :validate-on-rule-change="false"
              @keyup.enter="submitForm"
          >
            <el-form-item prop="username" class="mb-6">  <!-- mb -> margin-bottom，底部距离-->
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

            <el-form-item v-if="captchaImg" prop="captcha" class="mb-6">
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
                  >
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

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import {reactive, Ref, ref} from 'vue';
import {captcha} from '@/api/login';
import {FormInstance} from "element-plus";
import {useUserStore} from "@/store/modules/user";

const captchaImg = ref('') // 验证码图片 Base64

const loginForm: Ref<FormInstance | null> = ref(null);
const loginFormData = reactive({
  username: '',
  password: '',
  answer: '',
  captchaId: '',
})

const rules = reactive({
  username: [
    {required: true, message: '用户名不能为空', trigger: 'blur'},
  ],
  password: [
    {required: true, message: '密码不能为空', trigger: 'blur'},
  ],
  answer: [{required: false, message: '验证码不能为空', trigger: 'blur'}],
});


// 获取验证码
const loginVerify = async () => {
  const response = await captcha();
  if (response.Code == 200) {
    rules.answer.push({required: true, message: '验证码不能为空', trigger: 'blur'});
    captchaImg.value = response.Data.captchaImg
    loginFormData.captchaId = response.Data.captchaId;
    loginFormData.openCaptcha = response.Data.openCaptcha;
  }
}

loginVerify()

const userStore = useUserStore();

// 模拟的 login 函数，返回一个 Promise<boolean>
const login = async (): Promise<void> => {
  console.log("登陆成功")
  return await userStore.login(loginFormData)
};

const submitForm = () => {
  loginForm.value?.validate(async (valid: boolean) => {
    if (valid) {
      const flag = await login();
      // if (!flag) {
      //   await loginVerify();
      // }
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      });
      await loginVerify();
      return false;
    }
  });
};


</script>

<style scoped>

</style>
