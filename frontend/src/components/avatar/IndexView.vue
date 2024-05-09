<script setup>
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { UploadFile } from '@/apis/upload.js'

defineOptions({
  name: 'AvatarIndex'
})


const props = defineProps({
  target: {
    type: Object,
    default: null
  },
  targetKey: {
    type: String,
    default: ''
  }
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
      props.target[props.targetKey] = response.data.url
      ElMessage.success('上传成功')
    } else {
      ElMessage.error('上传失败')
    }
  } catch (error) {
    ElMessage.error('上传过程中出现错误')
  }
}
</script>

<template>
  <div>
    <el-upload
      class="avatar-uploader"
      :show-file-list="false"
      :before-upload="uploadFile"
    >
      <el-icon class="avatar-uploader-icon">
        <Plus />
      </el-icon>
    </el-upload>
  </div>

</template>

<style>
.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}
</style>
