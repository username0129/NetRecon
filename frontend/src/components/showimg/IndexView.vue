<script setup>
import defaultAvatar from '@/assets/defaut.png'
import { useUserStore } from '@/stores/modules/user.js'
import { computed } from 'vue'

defineOptions({
  name: 'ShowImgIndex'
})

const props = defineProps({
  imgType: {
    type: String,
    default: 'avatar'
  },
  imgSrc: {
    type: String,
    required: false,
    default: ''
  },
  preview: {
    type: Boolean,
    default: false
  }
})

const path = 'http://103.228.64.175:8081/'
const userStore = useUserStore()

const avatar = computed(() => {
  return props.imgSrc === '' ? `${path}${userStore.userInfo.avatar || defaultAvatar}` : `${path}${props.imgSrc}`
})

const file = computed(() => `${path}${props.imgSrc}`)

const previewSrcList = computed(() => props.preview ? [file.value] : [])

</script>

<template>
  <span class="headerAvatar">
    <template v-if="imgType === 'avatar'">
      <el-avatar
        v-if="imgSrc"
        :size="30"
        :src="avatar"
      />
      <el-avatar
        v-else
        :size="30"
        :src="defaultAvatar"
      />
    </template>
    <template v-if="imgType === 'file'">
      <el-image
        :src="file"
        class="file"
        :preview-src-list="previewSrcList"
        :preview-teleported="true"
      />
    </template>
  </span>
</template>


<style scoped>
.headerAvatar {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 8px;
}

.file {
  width: 80px;
  height: 80px;
  position: relative;
}
</style>
