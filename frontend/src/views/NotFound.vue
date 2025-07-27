<template>
  <div class="not-found-page">
    <div class="not-found-container">
      <el-result
        icon="error"
        title="404"
        sub-title="抱歉，您访问的页面不存在"
      >
        <template #extra>
          <el-button type="primary" @click="goHome">返回首页</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}
</script>

<style lang="scss" scoped>
.not-found-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
  
  .not-found-container {
    max-width: 500px;
    width: 100%;
    padding: 40px;
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
}
</style>