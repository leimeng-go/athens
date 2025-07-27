<template>
  <div class="upload-module-page">
    <el-card class="upload-card">
      <template #header>
        <div class="upload-card__title">
          <span>上传模块</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="文件上传" name="file">
          <div class="upload-container">
            <el-upload
              class="upload-area"
              drag
              action="#"
              :auto-upload="false"
              :on-change="handleFileChange"
              :on-remove="handleFileRemove"
              :limit="1"
              :file-list="fileList"
            >
              <el-icon class="upload-area__icon"><Upload /></el-icon>
              <div class="upload-area__text">拖拽文件到此处或 <em>点击上传</em></div>
              <template #tip>
                <div class="upload-area__tip">
                  支持 .zip 格式的 Go 模块文件，最大支持 50MB
                </div>
              </template>
            </el-upload>
            
            <div class="upload-actions" v-if="fileList.length > 0">
              <el-button type="primary" @click="submitUpload" :loading="uploadStore.isUploading">
                <el-icon><Upload /></el-icon>
                开始上传
              </el-button>
              <el-button @click="clearFiles">
                <el-icon><Delete /></el-icon>
                清空
              </el-button>
            </div>
            
            <el-progress
              v-if="uploadStore.isUploading"
              :percentage="uploadStore.uploadProgress"
              :format="percentageFormat"
              status="success"
            />
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="URL导入" name="url">
          <div class="url-import-container">
            <el-form :model="urlForm" label-width="80px" :rules="urlRules" ref="urlFormRef">
              <el-form-item label="模块URL" prop="url">
                <el-input 
                  v-model="urlForm.url" 
                  placeholder="请输入模块URL，例如：https://github.com/gomods/athens"
                  clearable
                />
              </el-form-item>
              
              <el-form-item label="版本" prop="version">
                <el-input 
                  v-model="urlForm.version" 
                  placeholder="请输入版本号，例如：v1.0.0（可选）"
                  clearable
                />
              </el-form-item>
              
              <el-form-item>
                <el-button type="primary" @click="submitUrlImport" :loading="importing">
                  <el-icon><Connection /></el-icon>
                  开始导入
                </el-button>
                <el-button @click="resetUrlForm">
                  <el-icon><RefreshRight /></el-icon>
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    
    <el-card class="upload-tasks-card">
      <template #header>
        <div class="upload-tasks-card__title">
          <span>上传任务</span>
          <div>
            <el-select v-model="statusFilter" placeholder="状态筛选" @change="handleStatusFilterChange">
              <el-option label="全部" value="" />
              <el-option label="等待中" value="pending" />
              <el-option label="处理中" value="processing" />
              <el-option label="成功" value="success" />
              <el-option label="失败" value="failed" />
            </el-select>
            <el-button type="primary" @click="refreshTasks" style="margin-left: 10px">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="uploadStore.loading"
        :data="uploadStore.uploadTasks"
        style="width: 100%"
        border
      >
        <el-table-column prop="modulePath" label="模块路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="version" label="版本" width="120" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                size="small" 
                type="primary" 
                @click="viewTaskDetail(row.id)"
                :disabled="row.status === 'pending'"
              >
                <el-icon><View /></el-icon>
              </el-button>
              <el-button 
                size="small" 
                type="danger" 
                @click="cancelTask(row)"
                :disabled="row.status !== 'pending' && row.status !== 'processing'"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="uploadStore.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUploadStore } from '@/stores/upload'
import { ElMessage, ElMessageBox, FormInstance, UploadFile, UploadUserFile } from 'element-plus'
import { formatDate as formatDateUtil } from '@/utils/format'

const uploadStore = useUploadStore()

// 标签页状态
const activeTab = ref('file')

// 文件上传状态
const fileList = ref<UploadUserFile[]>([])

// URL导入状态
const urlFormRef = ref<FormInstance>()
const urlForm = ref({
  url: '',
  version: ''
})
const urlRules = {
  url: [
    { required: true, message: '请输入模块URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ]
}
const importing = ref(false)

// 任务列表状态
const currentPage = ref(1)
const pageSize = ref(10)
const statusFilter = ref('')

// 初始化加载数据
onMounted(() => {
  uploadStore.fetchUploadTasks()
})

// 处理文件变更
function handleFileChange(file: UploadFile) {
  // 检查文件类型
  if (file.raw && !file.raw.type.includes('zip')) {
    ElMessage.error('只支持上传 .zip 格式的文件')
    fileList.value = []
    return
  }
  
  // 检查文件大小
  if (file.size && file.size > 50 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过 50MB')
    fileList.value = []
    return
  }
}

// 处理文件移除
function handleFileRemove() {
  fileList.value = []
}

// 清空文件
function clearFiles() {
  fileList.value = []
}

// 提交上传
async function submitUpload() {
  if (fileList.value.length === 0) {
    ElMessage.warning('请先选择要上传的文件')
    return
  }
  
  const file = fileList.value[0].raw
  if (!file) {
    ElMessage.error('文件无效')
    return
  }
  
  const taskId = await uploadStore.uploadModuleFile(file)
  if (taskId) {
    fileList.value = []
    ElMessage.success('上传任务已创建，任务ID: ' + taskId)
  }
}

// 提交URL导入
async function submitUrlImport() {
  if (!urlFormRef.value) return
  
  await urlFormRef.value.validate(async (valid) => {
    if (valid) {
      importing.value = true
      try {
        const taskId = await uploadStore.importFromUrl(
          urlForm.value.url,
          urlForm.value.version || undefined
        )
        if (taskId) {
          ElMessage.success('URL导入任务已创建，任务ID: ' + taskId)
          resetUrlForm()
        }
      } finally {
        importing.value = false
      }
    }
  })
}

// 重置URL表单
function resetUrlForm() {
  if (urlFormRef.value) {
    urlFormRef.value.resetFields()
  }
}

// 刷新任务列表
function refreshTasks() {
  uploadStore.fetchUploadTasks()
}

// 查看任务详情
function viewTaskDetail(id: string) {
  ElMessage.info(`查看任务详情: ${id}，此功能尚未实现`)
  // 可以跳转到详情页面或打开对话框
  // router.push(`/upload/tasks/${id}`)
}

// 取消任务
function cancelTask(row: any) {
  ElMessageBox.confirm(
    `确定要取消上传任务 "${row.modulePath}@${row.version}" 吗？`,
    '取消确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
    .then(() => {
      uploadStore.cancelTask(row.id)
    })
    .catch(() => {
      // 取消操作
    })
}

// 处理状态筛选变更
function handleStatusFilterChange() {
  uploadStore.updateStatusFilter(statusFilter.value as any)
}

// 处理分页变更
function handleSizeChange(val: number) {
  pageSize.value = val
  uploadStore.updatePagination(currentPage.value, pageSize.value)
}

function handleCurrentChange(val: number) {
  currentPage.value = val
  uploadStore.updatePagination(currentPage.value, pageSize.value)
}

// 格式化上传进度
function percentageFormat(percentage: number) {
  return percentage === 100 ? '上传完成' : `${percentage}%`
}

// 格式化日期
function formatDate(date: string) {
  return formatDateUtil(date)
}

// 获取状态类型
function getStatusType(status: string) {
  switch (status) {
    case 'pending': return 'info'
    case 'processing': return 'warning'
    case 'success': return 'success'
    case 'failed': return 'danger'
    default: return 'info'
  }
}

// 获取状态文本
function getStatusText(status: string) {
  switch (status) {
    case 'pending': return '等待中'
    case 'processing': return '处理中'
    case 'success': return '成功'
    case 'failed': return '失败'
    default: return '未知'
  }
}
</script>

<style lang="scss" scoped>
.upload-module-page {
  .upload-card,
  .upload-tasks-card {
    margin-bottom: 20px;
  }
  
  .upload-container {
    padding: 20px 0;
  }
  
  .upload-actions {
    margin-top: 20px;
    display: flex;
    justify-content: center;
    gap: 10px;
  }
  
  .url-import-container {
    max-width: 600px;
    margin: 20px auto;
  }
  
  .upload-tasks-card__title {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>