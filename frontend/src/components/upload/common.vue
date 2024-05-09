<script setup>
import { ElMessage } from 'element-plus'
import { UploadFile } from '@/apis/upload.js'

defineOptions({
  name: 'UploadCommon'
})

const emit = defineEmits(['on-success'])

// 图片类型判断
const isImageMime = (type) => {
  return ['image/jpeg', 'image/png', 'image/gif'].includes(type)
}

async function uploadFile(file) {
  if (!isImageMime(file.type)) {
    ElMessage.error('仅支持上传 jpg, png, gif 格式')
    return false
  }

  const formData = new FormData()
  formData.append('file', file)

  try {
    const response = await UploadFile(formData)
    if (response.code === 200) {
      emit('on-success', 'true')
      ElMessage.success('上传成功')
    } else {
      ElMessage.error('上传失败')
    }
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

</script>

<template>
  <div>
    <el-upload
      :before-upload="uploadFile"
      :show-file-list="false"
      class="upload-btn"
    >
      <el-button type="primary">上传图片</el-button>
    </el-upload>
  </div>
</template>

<style scoped lang="scss">

</style>
