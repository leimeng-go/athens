<template>
  <div class="module-detail-page">
    <div class="page-header">
      <el-page-header @back="goBack">
        <template #content>
          <div class="page-title">
            <span class="module-path">{{ modulePath }}</span>
            <el-tag v-if="moduleDetail.latestVersion" type="success" effect="dark" size="small">
              {{ moduleDetail.latestVersion }}
            </el-tag>
          </div>
        </template>
      </el-page-header>
    </div>
    
    <div class="loading-container" v-if="loading">
      <el-skeleton :rows="10" animated />
    </div>
    
    <template v-else>
      <el-row :gutter="20">
        <el-col :xs="24" :md="16">
          <el-card class="module-info-card">
            <template #header>
              <div class="card-header">
                <span>模块信息</span>
                <el-button-group>
                  <el-button type="primary" size="small" @click="downloadLatestVersion">
                    <el-icon><Download /></el-icon>
                    下载最新版本
                  </el-button>
                  <el-button type="success" size="small" @click="copyGoGetCommand">
                    <el-icon><DocumentCopy /></el-icon>
                    复制 go get 命令
                  </el-button>
                </el-button-group>
              </div>
            </template>
            
            <div class="module-info">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="模块路径" :span="2">
                  {{ moduleDetail.path }}
                </el-descriptions-item>
                <el-descriptions-item label="最新版本">
                  {{ moduleDetail.latestVersion || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="版本数量">
                  {{ moduleDetail.versions || 0 }}
                </el-descriptions-item>
                <el-descriptions-item label="总下载次数">
                  {{ formatNumber(moduleDetail.downloads) }}
                </el-descriptions-item>
                <el-descriptions-item label="最后下载时间">
                  {{ formatDate(moduleDetail.lastDownloaded) }}
                </el-descriptions-item>
                <el-descriptions-item label="总大小">
                  {{ formatFileSize(moduleDetail.size) }}
                </el-descriptions-item>
                <el-descriptions-item label="首次发布">
                  {{ formatDate(moduleDetail.firstPublished) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
            
            <div class="module-description" v-if="moduleDetail.description">
              <h3>模块描述</h3>
              <p>{{ moduleDetail.description }}</p>
            </div>
            
            <div class="module-readme" v-if="moduleDetail.readme">
              <h3>README</h3>
              <div class="readme-content" v-html="renderedReadme"></div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :md="8">
          <el-card class="version-list-card">
            <template #header>
              <div class="card-header">
                <span>版本列表</span>
              </div>
            </template>
            
            <el-table
              :data="versionList"
              style="width: 100%"
              v-loading="versionsLoading"
            >
              <el-table-column prop="version" label="版本" min-width="100">
                <template #default="{ row }">
                  <el-tag :type="getVersionTagType(row)" effect="plain">
                    {{ row.version }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="downloads" label="下载量" width="100" align="center" />
              <el-table-column prop="size" label="大小" width="100" align="center">
                <template #default="{ row }">
                  {{ formatFileSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="published" label="发布时间" min-width="150">
                <template #default="{ row }">
                  {{ formatDate(row.published) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80" align="center" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" type="primary" @click="downloadVersion(row)">
                    <el-icon><Download /></el-icon>
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
          
          <el-card class="dependencies-card" v-if="moduleDetail.dependencies && moduleDetail.dependencies.length > 0">
            <template #header>
              <div class="card-header">
                <span>依赖项</span>
              </div>
            </template>
            
            <el-table
              :data="moduleDetail.dependencies"
              style="width: 100%"
            >
              <el-table-column prop="path" label="依赖模块" min-width="200" show-overflow-tooltip>
                <template #default="{ row }">
                  <el-link type="primary" @click="viewDependency(row.path)">
                    {{ row.path }}
                  </el-link>
                </template>
              </el-table-column>
              <el-table-column prop="version" label="版本" width="120" />
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDownloadStore } from '@/stores/download'
import { ElMessage } from 'element-plus'
import { formatDate, formatFileSize, formatNumber } from '@/utils/format'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const route = useRoute()
const router = useRouter()
const downloadStore = useDownloadStore()

// 模块路径
const modulePath = computed(() => {
  return route.params.path as string
})

// 加载状态
const loading = ref(true)
const versionsLoading = ref(true)

// 模块详情
const moduleDetail = ref<any>({
  path: '',
  latestVersion: '',
  versions: 0,
  downloads: 0,
  lastDownloaded: '',
  size: 0,
  firstPublished: '',
  description: '',
  readme: '',
  dependencies: []
})

// 版本列表
const versionList = ref<any[]>([])

// 渲染README
const renderedReadme = computed(() => {
  if (!moduleDetail.value.readme) return ''
  const html = marked(moduleDetail.value.readme)
  return DOMPurify.sanitize(html)
})

// 初始化
onMounted(async () => {
  if (!modulePath.value) {
    ElMessage.error('模块路径无效')
    goBack()
    return
  }
  
  try {
    // 获取模块详情
    const detail = await downloadStore.fetchModuleDetail(modulePath.value)
    moduleDetail.value = detail
    
    // 获取版本列表
    const versions = await downloadStore.fetchModuleVersions(modulePath.value)
    versionList.value = versions
  } catch (error) {
    ElMessage.error('获取模块信息失败')
    console.error(error)
  } finally {
    loading.value = false
    versionsLoading.value = false
  }
})

// 返回上一页
function goBack() {
  router.push('/download')
}

// 下载最新版本
function downloadLatestVersion() {
  if (!moduleDetail.value.latestVersion) {
    ElMessage.warning('没有可用的版本')
    return
  }
  
  const url = `/download/module/${encodeURIComponent(moduleDetail.value.path)}@${moduleDetail.value.latestVersion}`
  ElMessage.success(`开始下载模块: ${moduleDetail.value.path}@${moduleDetail.value.latestVersion}`)
  
  // 创建一个隐藏的a标签并模拟点击下载
  const link = document.createElement('a')
  link.href = url
  link.target = '_blank'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 下载特定版本
function downloadVersion(version: any) {
  const url = `/download/module/${encodeURIComponent(moduleDetail.value.path)}@${version.version}`
  ElMessage.success(`开始下载模块: ${moduleDetail.value.path}@${version.version}`)
  
  // 创建一个隐藏的a标签并模拟点击下载
  const link = document.createElement('a')
  link.href = url
  link.target = '_blank'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 复制go get命令
function copyGoGetCommand() {
  const command = `go get ${moduleDetail.value.path}@${moduleDetail.value.latestVersion || 'latest'}`
  navigator.clipboard.writeText(command).then(() => {
    ElMessage.success('命令已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 查看依赖项
function viewDependency(path: string) {
  router.push(`/download/modules/${encodeURIComponent(path)}`)
}

// 获取版本标签类型
function getVersionTagType(version: any) {
  if (version.version === moduleDetail.value.latestVersion) {
    return 'success'
  }
  
  if (version.version.includes('beta') || version.version.includes('alpha')) {
    return 'warning'
  }
  
  if (version.version.includes('rc')) {
    return 'info'
  }
  
  return ''
}
</script>

<style lang="scss" scoped>
.module-detail-page {
  .page-header {
    margin-bottom: 20px;
    padding: 15px;
    background-color: #fff;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
  
  .page-title {
    display: flex;
    align-items: center;
    gap: 10px;
    
    .module-path {
      font-size: 16px;
      font-weight: bold;
    }
  }
  
  .module-info-card,
  .version-list-card,
  .dependencies-card {
    margin-bottom: 20px;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .module-description,
  .module-readme {
    margin-top: 20px;
    
    h3 {
      font-size: 16px;
      margin-bottom: 10px;
      padding-bottom: 5px;
      border-bottom: 1px solid #eee;
    }
  }
  
  .readme-content {
    padding: 15px;
    background-color: #f8f9fa;
    border-radius: 4px;
    max-height: 500px;
    overflow-y: auto;
    
    :deep(h1) {
      font-size: 24px;
      margin-bottom: 16px;
    }
    
    :deep(h2) {
      font-size: 20px;
      margin-bottom: 14px;
    }
    
    :deep(h3) {
      font-size: 18px;
      margin-bottom: 12px;
    }
    
    :deep(p) {
      margin-bottom: 10px;
    }
    
    :deep(pre) {
      background-color: #f1f1f1;
      padding: 10px;
      border-radius: 4px;
      overflow-x: auto;
    }
    
    :deep(code) {
      font-family: monospace;
      background-color: #f1f1f1;
      padding: 2px 4px;
      border-radius: 3px;
    }
    
    :deep(ul), :deep(ol) {
      padding-left: 20px;
      margin-bottom: 10px;
    }
    
    :deep(a) {
      color: #409eff;
      text-decoration: none;
      
      &:hover {
        text-decoration: underline;
      }
    }
    
    :deep(table) {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: 10px;
      
      th, td {
        border: 1px solid #ddd;
        padding: 8px;
      }
      
      th {
        background-color: #f2f2f2;
        text-align: left;
      }
      
      tr:nth-child(even) {
        background-color: #f9f9f9;
      }
    }
  }
  
  .loading-container {
    padding: 20px;
    background-color: #fff;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
}
</style>