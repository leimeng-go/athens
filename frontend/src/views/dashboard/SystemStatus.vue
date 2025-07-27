<template>
  <div class="system-status-page">
    <h1>系统状态</h1>
    
    <!-- 加载状态 -->
    <el-skeleton :loading="systemStatusLoading" animated :count="1" v-if="systemStatusLoading">
      <template #template>
        <div style="padding: 20px;">
          <el-skeleton-item variant="p" style="width: 100%" />
          <el-skeleton-item variant="text" style="margin-right: 16px" />
          <el-skeleton-item variant="text" style="width: 30%" />
        </div>
      </template>
    </el-skeleton>
    
    <!-- 系统状态卡片 -->
    <el-card class="system-status-card" v-if="!systemStatusLoading">
      <template #header>
        <div class="card-header">
          <el-icon><Monitor /></el-icon>
          <span>系统状态详情</span>
          <el-button type="primary" size="small" @click="fetchSystemStatus" style="margin-left: auto;">
            <el-icon><Refresh /></el-icon>
            刷新
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
    
    <!-- 系统信息卡片 -->
    <el-card class="system-info-card" v-if="!systemStatusLoading">
      <template #header>
        <div class="card-header">
          <el-icon><InfoFilled /></el-icon>
          <span>系统信息</span>
        </div>
      </template>
      <div class="system-info-content">
        <p>Athens是一个Go模块代理服务器，它可以缓存Go模块，提高构建速度，并提供私有模块托管功能。</p>
        <p>本系统提供了Athens的管理界面，可以管理模块、仓库、上传和下载等功能。</p>
        <p>当前系统版本: {{ systemStatus.version }}</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useDashboardStore } from '@/stores/dashboard'
import { Monitor, Refresh, InfoFilled } from '@element-plus/icons-vue'

// 使用仪表盘store
const dashboardStore = useDashboardStore()
const { systemStatus, systemStatusLoading } = storeToRefs(dashboardStore)
const { fetchSystemStatus } = dashboardStore

// 页面加载时获取系统状态
onMounted(() => {
  fetchSystemStatus()
})
</script>

<style lang="scss" scoped>
.system-status-page {
  h1 {
    margin-bottom: 20px;
  }
  
  .system-status-card,
  .system-info-card {
    margin-bottom: 20px;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    
    .el-icon {
      margin-right: 8px;
    }
  }
  
  .system-status-content {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
    
    .status-item {
      display: flex;
      align-items: center;
      
      .label {
        font-weight: bold;
        margin-right: 10px;
        min-width: 100px;
      }
      
      .value {
        flex: 1;
      }
    }
  }
  
  .system-info-content {
    p {
      margin-bottom: 10px;
      line-height: 1.6;
    }
  }
}
</style>