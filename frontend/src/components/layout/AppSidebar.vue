<template>
  <aside class="admin-layout__sidebar">
    <div class="sidebar-header">
      <img src="@/assets/images/logo.png" alt="Athens Logo" class="sidebar-logo">
      <h2 class="sidebar-title">Athens Admin</h2>
    </div>
    
    <!-- 主导航菜单 -->
    <el-menu
      :default-active="activeIndex"
      class="sidebar-menu"
      :router="true"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Odometer /></el-icon>
        <span>仪表盘</span>
      </el-menu-item>
      
      <el-menu-item index="/dashboard/system">
        <el-icon><Monitor /></el-icon>
        <span>系统状态</span>
      </el-menu-item>
      
      <el-menu-item index="/dashboard/activities">
        <el-icon><Bell /></el-icon>
        <span>最近活动</span>
      </el-menu-item>
      
      <el-menu-item index="/repositories">
        <el-icon><Files /></el-icon>
        <span>仓库管理</span>
      </el-menu-item>
      
      <el-menu-item index="/upload">
        <el-icon><Upload /></el-icon>
        <span>上传模块</span>
      </el-menu-item>
      
      <el-menu-item index="/download">
        <el-icon><Download /></el-icon>
        <span>下载模块</span>
      </el-menu-item>
      
      <el-menu-item index="/settings">
        <el-icon><Setting /></el-icon>
        <span>系统设置</span>
      </el-menu-item>
    </el-menu>
    
    <!-- 子菜单 -->
    <div class="sidebar-submenu" v-if="showSubmenu">
      <h3 class="submenu-title">{{ currentModuleTitle }}</h3>
      
      <el-menu
        :default-active="activeMenuItem"
        class="submenu"
        :router="true"
      >
        <template v-if="currentModule === 'repositories'">
          <el-menu-item index="/repositories">
            <el-icon><List /></el-icon>
            <span>仓库列表</span>
          </el-menu-item>
          <el-menu-item index="/repositories/stats">
            <el-icon><DataAnalysis /></el-icon>
            <span>仓库统计</span>
          </el-menu-item>
        </template>
        
        <template v-else-if="currentModule === 'upload'">
          <el-menu-item index="/upload">
            <el-icon><Upload /></el-icon>
            <span>上传模块</span>
          </el-menu-item>
          <el-menu-item index="/upload/tasks">
            <el-icon><List /></el-icon>
            <span>上传任务</span>
          </el-menu-item>
          <el-menu-item index="/upload/import">
            <el-icon><Link /></el-icon>
            <span>URL导入</span>
          </el-menu-item>
        </template>
        
        <template v-else-if="currentModule === 'download'">
          <el-menu-item index="/download">
            <el-icon><Download /></el-icon>
            <span>模块列表</span>
          </el-menu-item>
          <el-menu-item index="/download/stats">
            <el-icon><DataAnalysis /></el-icon>
            <span>下载统计</span>
          </el-menu-item>
          <el-menu-item index="/download/popular">
            <el-icon><Star /></el-icon>
            <span>热门模块</span>
          </el-menu-item>
        </template>
        
        <template v-else-if="currentModule === 'settings'">
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>基本设置</span>
          </el-menu-item>
          <el-menu-item index="/settings/storage">
            <el-icon><FolderOpened /></el-icon>
            <span>存储设置</span>
          </el-menu-item>
          <el-menu-item index="/settings/proxy">
            <el-icon><Connection /></el-icon>
            <span>代理设置</span>
          </el-menu-item>
          <el-menu-item index="/settings/logs">
            <el-icon><Document /></el-icon>
            <span>日志管理</span>
          </el-menu-item>
        </template>
        
        <template v-else-if="currentModule === 'dashboard'">
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>概览</span>
          </el-menu-item>
        </template>
      </el-menu>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Files, Odometer, Upload, Download, Setting, List, DataAnalysis, Link, Star, FolderOpened, Connection, Document, Monitor, Bell } from '@element-plus/icons-vue'

const route = useRoute()

// 计算当前模块
const currentModule = computed(() => {
  const path = route.path
  if (path.startsWith('/repositories')) return 'repositories'
  if (path.startsWith('/upload')) return 'upload'
  if (path.startsWith('/download')) return 'download'
  if (path.startsWith('/settings')) return 'settings'
  return 'dashboard'
})

// 计算当前模块标题
const currentModuleTitle = computed(() => {
  switch (currentModule.value) {
    case 'repositories': return '仓库管理'
    case 'upload': return '上传模块'
    case 'download': return '下载模块'
    case 'settings': return '系统设置'
    default: return '仪表盘'
  }
})

// 计算当前激活的菜单项
const activeMenuItem = computed(() => route.path)

// 计算主导航激活的菜单项
const activeIndex = computed(() => {
  const path = route.path
  if (path.startsWith('/repositories')) return '/repositories'
  if (path.startsWith('/upload')) return '/upload'
  if (path.startsWith('/download')) return '/download'
  if (path.startsWith('/settings')) return '/settings'
  return '/dashboard'
})

// 是否显示子菜单
const showSubmenu = computed(() => {
  return route.path !== '/' && currentModule.value !== ''
})
</script>

<style lang="scss" scoped>
.admin-layout__sidebar {
  width: 220px;
  background-color: #fff;
  border-right: 1px solid #e6e6e6;
  height: 100%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  
  .sidebar-header {
    padding: 16px;
    border-bottom: 1px solid #e6e6e6;
    display: flex;
    align-items: center;
    flex-direction: column;
    margin-bottom: 10px;
  }
  
  .sidebar-logo {
    height: 40px;
    margin-bottom: 10px;
  }
  
  .sidebar-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin: 0;
  }
  
  .sidebar-menu {
    border-right: none;
    
    :deep(.el-menu-item) {
      height: 50px;
      line-height: 50px;
      
      .el-icon {
        margin-right: 5px;
      }
    }
  }
  
  .sidebar-submenu {
    margin-top: 10px;
    border-top: 1px solid #e6e6e6;
    padding-top: 10px;
    
    .submenu-title {
      font-size: 14px;
      font-weight: 600;
      color: #606266;
      margin: 0 0 10px 16px;
    }
    
    .submenu {
      border-right: none;
      
      :deep(.el-menu-item) {
        height: 40px;
        line-height: 40px;
        padding-left: 20px !important;
        
        .el-icon {
          margin-right: 5px;
        }
      }
    }
  }
}
</style>