<template>
  <div class="dashboard-page">
    <h1>仪表盘</h1>
    
    <!-- 加载状态 -->
    <el-skeleton :loading="loading" animated :count="4" v-if="loading">
      <template #template>
        <div style="padding: 20px;">
          <el-skeleton-item variant="p" style="width: 100%" />
          <el-skeleton-item variant="text" style="margin-right: 16px" />
          <el-skeleton-item variant="text" style="width: 30%" />
        </div>
      </template>
    </el-skeleton>
    
    <!-- 统计卡片 -->
    <div class="stat-cards" v-if="!loading">
      <el-card class="stat-card">
        <template #header>
          <div class="card-header">
            <el-icon><Document /></el-icon>
            <span>模块总数</span>
          </div>
        </template>
        <div class="card-value">{{ dashboardData.stats.totalModules }}</div>
      </el-card>
      
      <el-card class="stat-card">
        <template #header>
          <div class="card-header">
            <el-icon><Download /></el-icon>
            <span>下载总次数</span>
          </div>
        </template>
        <div class="card-value">{{ dashboardData.stats.totalDownloads }}</div>
      </el-card>
      
      <el-card class="stat-card">
        <template #header>
          <div class="card-header">
            <el-icon><FolderOpened /></el-icon>
            <span>仓库总数</span>
          </div>
        </template>
        <div class="card-value">{{ dashboardData.stats.totalRepositories }}</div>
      </el-card>
      
      <el-card class="stat-card">
        <template #header>
          <div class="card-header">
            <el-icon><Cpu /></el-icon>
            <span>存储使用量</span>
          </div>
        </template>
        <div class="card-value">{{ dashboardData.stats.storageUsed }}</div>
      </el-card>
    </div>
    
    <!-- 系统状态 -->
    <el-card class="system-status-card" v-if="!systemStatusLoading">
      <template #header>
        <div class="card-header">
          <el-icon><Monitor /></el-icon>
          <span>系统状态</span>
          <el-button type="primary" size="small" circle @click="fetchSystemStatus" style="margin-left: auto;">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </template>
      <div class="system-status-content">
        <div class="status-item">
          <span class="label">状态:</span>
          <span class="value">
            <el-tag :type="systemStatus.status === 'healthy' ? 'success' : 'danger'">
              {{ systemStatus.status === 'healthy' ? '正常' : '异常' }}
            </el-tag>
          </span>
        </div>
        <div class="status-item">
          <span class="label">运行时间:</span>
          <span class="value">{{ systemStatus.uptime }}</span>
        </div>
        <div class="status-item">
          <span class="label">版本:</span>
          <span class="value">{{ systemStatus.version }}</span>
        </div>
        <div class="status-item">
          <span class="label">Go版本:</span>
          <span class="value">{{ systemStatus.goVersion }}</span>
        </div>
        <div class="status-item">
          <span class="label">内存使用:</span>
          <span class="value">{{ systemStatus.memoryUsage }}</span>
        </div>
        <div class="status-item">
          <span class="label">CPU使用:</span>
          <span class="value">{{ systemStatus.cpuUsage }}</span>
        </div>
      </div>
    </el-card>
    
    <!-- 热门模块 -->
    <el-card class="popular-modules-card" v-if="!loading">
      <template #header>
        <div class="card-header">
          <el-icon><Star /></el-icon>
          <span>热门模块</span>
        </div>
      </template>
      <el-table :data="dashboardData.popularModules" style="width: 100%">
        <el-table-column prop="path" label="模块路径" />
        <el-table-column prop="downloads" label="下载次数" width="120" />
      </el-table>
    </el-card>
    
    <!-- 最近活动 -->
    <el-card class="recent-activities-card" v-if="!loading">
      <template #header>
        <div class="card-header">
          <el-icon><Bell /></el-icon>
          <span>最近活动</span>
        </div>
      </template>
      <el-timeline>
        <el-timeline-item
          v-for="activity in dashboardData.recentActivities"
          :key="activity.id"
          :timestamp="activity.timestamp"
          :type="getActivityType(activity.type)"
        >
          {{ activity.message }}
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { Document, Download, FolderOpened, Cpu, Monitor, Refresh, Star, Bell } from '@element-plus/icons-vue'

// 使用仪表盘store
const dashboardStore = useDashboardStore()
const { 
  dashboardData, 
  loading, 
  systemStatus, 
  systemStatusLoading,
  fetchDashboardData,
  fetchSystemStatus,
  fetchRecentActivities
} = dashboardStore

// 获取活动类型对应的图标类型
function getActivityType(type: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' {
  switch (type) {
    case 'download':
      return 'primary'
    case 'upload':
      return 'success'
    case 'system':
      return 'warning'
    default:
      return 'info'
  }
}

// 组件挂载时加载数据
onMounted(() => {
  fetchDashboardData()
  fetchSystemStatus()
})
</script>

<style lang="scss" scoped>
.dashboard-page {
  padding: 20px;
  
  h1 {
    margin-bottom: 20px;
    font-size: 24px;
    color: #303133;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    
    .el-icon {
      margin-right: 8px;
    }
  }
  
  .stat-cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 20px;
    margin-bottom: 20px;
    
    .stat-card {
      .card-value {
        font-size: 28px;
        font-weight: bold;
        color: #409EFF;
        text-align: center;
        padding: 10px 0;
      }
    }
  }
  
  .system-status-card {
    margin-bottom: 20px;
    
    .system-status-content {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 15px;
      
      .status-item {
        display: flex;
        flex-direction: column;
        
        .label {
          font-size: 14px;
          color: #909399;
          margin-bottom: 5px;
        }
        
        .value {
          font-size: 16px;
          font-weight: 500;
        }
      }
    }
  }
  
  .popular-modules-card,
  .recent-activities-card {
    margin-bottom: 20px;
  }
}
</style>