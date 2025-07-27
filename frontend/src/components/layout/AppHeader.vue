<template>
  <header class="admin-layout__header">
    <div class="header-left">
      <el-icon class="menu-toggle" :size="24"><Menu /></el-icon>
    </div>
    
    <div class="header-right">
      <el-dropdown trigger="click">
        <div class="user-info">
          <el-avatar :size="32" icon="User" />
          <span class="user-name">管理员</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item>
              <el-icon><Setting /></el-icon>
              <span>个人设置</span>
            </el-dropdown-item>
            <el-dropdown-item divided>
              <el-icon><SwitchButton /></el-icon>
              <span>退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ArrowDown, Menu, Setting, SwitchButton } from '@element-plus/icons-vue'
</script>

<style lang="scss" scoped>
.admin-layout__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  position: relative;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  
  .menu-toggle {
    cursor: pointer;
    padding: 8px;
    border-radius: 4px;
    transition: background-color 0.3s;
    
    &:hover {
      background-color: #f5f7fa;
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 8px;
  height: 40px;
  border-radius: 4px;
  transition: background-color 0.3s;
  
  &:hover {
    background-color: #f5f7fa;
  }
  
  .user-name {
    margin: 0 8px;
    font-size: 14px;
  }
}
</style>