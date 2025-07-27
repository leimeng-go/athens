<template>
  <div class="activities-page">
    <h1>最近活动</h1>
    
    <!-- 加载状态 -->
    <el-skeleton :loading="loading" animated :count="5" v-if="loading">
      <template #template>
        <div style="padding: 20px;">
          <el-skeleton-item variant="p" style="width: 100%" />
          <el-skeleton-item variant="text" style="margin-right: 16px" />
          <el-skeleton-item variant="text" style="width: 30%" />
        </div>
      </template>
    </el-skeleton>
    
    <!-- 活动列表 -->
    <el-card class="activities-card" v-if="!loading">
      <template #header>
        <div class="card-header">
          <el-icon><Bell /></el-icon>
          <span>活动列表</span>
          <el-button type="primary" size="small" @click="fetchActivities" style="margin-left: auto;">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      <div class="activities-content">
        <el-timeline>
          <el-timeline-item
            v-for="activity in recentActivities"
            :key="activity.id"
            :type="getActivityType(activity.type)"
            :timestamp="formatTime(activity.timestamp)"
          >
            <div class="activity-item">
              <span class="activity-message">{{ activity.message }}</span>
              <el-tag size="small" :type="getActivityTagType(activity.type)">
                {{ getActivityTypeText(activity.type) }}
              </el-tag>
            </div>
          </el-timeline-item>
        </el-timeline>
        
        <div class="empty-state" v-if="recentActivities.length === 0">
          <el-empty description="暂无活动记录" />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useDashboardStore } from '@/stores/dashboard'
import { Bell, Refresh } from '@element-plus/icons-vue'

// 使用仪表盘store
const dashboardStore = useDashboardStore()
const { dashboardData } = storeToRefs(dashboardStore)
const { fetchDashboardData } = dashboardStore

// 加载状态
const loading = ref(false)

// 最近活动
const recentActivities = computed(() => dashboardData.value.recentActivities)

// 获取活动数据
async function fetchActivities() {
  loading.value = true
  try {
    await fetchDashboardData()
  } finally {
    loading.value = false
  }
}

// 格式化时间
function formatTime(timestamp: string) {
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', { 
      year: 'numeric', 
      month: '2-digit', 
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (e) {
    return timestamp
  }
}

// 获取活动类型对应的图标类型
function getActivityType(type: string) {
  switch (type) {
    case 'download': return 'success'
    case 'upload': return 'primary'
    case 'system': return 'warning'
    default: return 'info'
  }
}

// 获取活动类型对应的标签类型
function getActivityTagType(type: string) {
  switch (type) {
    case 'download': return 'success'
    case 'upload': return 'primary'
    case 'system': return 'warning'
    default: return 'info'
  }
}

// 获取活动类型对应的文本
function getActivityTypeText(type: string) {
  switch (type) {
    case 'download': return '下载'
    case 'upload': return '上传'
    case 'system': return '系统'
    default: return '其他'
  }
}

// 页面加载时获取活动数据
onMounted(() => {
  fetchActivities()
})
</script>

<style lang="scss" scoped>
.activities-page {
  h1 {
    margin-bottom: 20px;
  }
  
  .activities-card {
    margin-bottom: 20px;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    
    .el-icon {
      margin-right: 8px;
    }
  }
  
  .activities-content {
    .activity-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .activity-message {
        flex: 1;
        margin-right: 10px;
      }
    }
    
    .empty-state {
      padding: 20px 0;
    }
  }
}
</style>