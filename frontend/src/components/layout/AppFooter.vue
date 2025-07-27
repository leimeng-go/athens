<template>
  <footer class="admin-layout__footer">
    <div class="footer-content">
      <span>© {{ currentYear }} Athens Go Module Proxy</span>
      <span class="footer-divider">|</span>
      <span>Version: {{ version }}</span>
      <span class="footer-divider">|</span>
      <a href="https://docs.gomods.io" target="_blank">文档</a>
      <span class="footer-divider">|</span>
      <a href="https://api.gomods.io" target="_blank">API</a>
      <span class="footer-divider">|</span>
      <a href="https://github.com/gomods/athens" target="_blank">GitHub</a>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// 当前年份
const currentYear = computed(() => new Date().getFullYear())

// 版本号
const version = ref('v0.11.0') // 这里可以从环境变量或API获取
</script>

<style lang="scss" scoped>
.admin-layout__footer {
  height: 40px;
  background-color: #fff;
  border-top: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #999;
  
  .footer-content {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .footer-divider {
    margin: 0 8px;
  }
  
  a {
    color: #999;
    text-decoration: none;
    
    &:hover {
      color: #409eff;
    }
  }
}

@media (max-width: 768px) {
  .admin-layout__footer {
    height: auto;
    padding: 10px 0;
    
    .footer-content {
      flex-direction: column;
      gap: 5px;
    }
    
    .footer-divider {
      display: none;
    }
  }
}
</style>