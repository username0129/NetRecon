<template>
  <div>
    <div class="my-table-box">
      <warning-bar title="注：没有注释" />
      <div class="my-btn-list">
        <UploadCommon @on-success="getTableData" />
        <el-input v-model="searchInfo.filename" placeholder="请输入文件名" />
        <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
      </div>

      <el-table :data="tableData">
        <el-table-column align="left" label="图片预览" width="150">
          <template #default="scope">
            <ShowImgIndex
              img-type="file"
              :img-src="scope.row.url"
              preview
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="文件名" prop="filename" width="250" />
        <el-table-column align="left" label="链接" prop="url" min-width="300" />
        <el-table-column align="left" label="操作" width="250">
          <template #default="scope">
            <el-button
              icon="Delete"
              type="danger"
              @click="deleteFile(scope.row)"
            >删除
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
  </div>
</template>

<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { DeleteFile, FetchFiles } from '@/apis/upload.js'
import UploadCommon from '@/components/upload/common.vue'
import ShowImgIndex from '@/components/showimg/IndexView.vue'

defineOptions({
  name: 'UploadIndex'
})

const searchInfo = ref({})
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
async function getTableData() {
  try {
    const response = await FetchFiles({
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

const deleteFile = async (row) => {
  ElMessageBox.confirm('此操作将永久删除文件, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      const res = await DeleteFile({ uuid: row.uuid })
      if (res.code === 200) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })
        await getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消删除'
      })
    })
}

</script>

<style scoped>

</style>
