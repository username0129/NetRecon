<script setup>
import Aside from '@/views/layout/aside/AsideIndex.vue'
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { fmtTitle } from '@/utils/fmtRouterTitle.js'
import { useRouteStore } from '@/stores/modules/route.js'
import { useUserStore } from '@/stores/modules/user.js'
import { ArrowDown, Setting } from '@element-plus/icons-vue'

defineOptions({
  name: 'LayoutIndex'
})

const route = useRoute()
const router = useRouter()
const routeStore = useRouteStore()
const userStore = useUserStore()

const matched = computed(() => route.matched)

function getRoleName() {
  switch (userStore.userInfo.valueOf().authorityId) {
    case 1:
      return '系统管理员'
    case 2:
      return '普通用户'
    default:
      return '未知角色'
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
                              <span style="margin-left: 5px">{{ userStore.userInfo.nickname }}</span>
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
