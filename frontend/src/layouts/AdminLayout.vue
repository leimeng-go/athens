<template>
  <div class="admin-layout">
    <AppHeader />
    
    <div class="admin-layout__main">
      <AppSidebar />
      
      <div class="admin-layout__content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
    
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
</script>

<style lang="scss" scoped>
.admin-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.admin-layout__main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.admin-layout__content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>