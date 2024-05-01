<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { FetchPortScanResult } from '@/apis/portscan.js'

const route = useRoute()
const taskUUID = ref('')

taskUUID.value = route.query.uuid

console.log(taskUUID.value)

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])


const getTableData = async () => {
  const response = await FetchPortScanResult({ page: page.value, pageSize: pageSize.value, taskUUID: taskUUID.value })
  if (response.code == 200) {
    tableData.value = response.data.data
    total.value = response.data.total
    page.value = response.data.page
    pageSize.value = response.data.pageSize
  } else if (response.code === 404) {
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
}

getTableData()


</script>

<template>
  <div>
    <warning-bar title="注：没有注释" />
  </div>
</template>

<style lang="scss" scoped>

</style>
