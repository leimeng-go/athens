<template>
  <nav class="app-nav">
    <el-menu
      mode="horizontal"
      :default-active="activeIndex"
      :ellipsis="false"
      router
      @select="handleSelect"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Odometer /></el-icon>
        <span>仪表盘</span>
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
  </nav>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 计算当前激活的菜单项
const activeIndex = computed(() => {
  const path = route.path
  if (path.startsWith('/repositories')) return '/repositories'
  if (path.startsWith('/upload')) return '/upload'
  if (path.startsWith('/download')) return '/download'
  if (path.startsWith('/settings')) return '/settings'
  return '/dashboard'
})

// 处理菜单选择事件
const handleSelect = (key: string) => {
  console.log('Navigation selected:', key)
}
</script>

<style lang="scss" scoped>
.app-nav {
  flex: 1;
  display: flex;
  justify-content: center;
  
  :deep(.el-menu--horizontal) {
    border-bottom: none;
    
    .el-menu-item {
      height: 60px;
      line-height: 60px;
      
      .el-icon {
        margin-right: 5px;
      }
    }
  }
}
</style>