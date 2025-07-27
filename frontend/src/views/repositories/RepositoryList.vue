<template>
  <div class="repository-list-page">
    <el-card class="stats-card">
      <template #header>
        <div class="stats-card__title">
          <span>仓库统计</span>
          <el-button type="primary" size="small" @click="refreshStats">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="8">
            <div class="stats-item">
              <div class="stats-icon">
                <el-icon><Files /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ repositoryStore.stats.totalRepositories }}</div>
                <div class="stats-label">仓库总数</div>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="stats-item">
              <div class="stats-icon">
                <el-icon><Collection /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ repositoryStore.stats.totalModules }}</div>
                <div class="stats-label">模块总数</div>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="stats-item">
              <div class="stats-icon">
                <el-icon><Coin /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ repositoryStore.stats.totalSize }}</div>
                <div class="stats-label">存储空间</div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <div class="data-table">
      <div class="data-table__header">
        <div class="data-table__title">仓库列表</div>
        <div class="data-table__actions">
          <el-input
            v-model="searchQuery"
            placeholder="搜索仓库..."
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
          
          <el-button type="primary" @click="refreshList">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
      
      <el-table
        v-loading="repositoryStore.loading"
        :data="repositoryStore.repositories"
        style="width: 100%"
        border
        @sort-change="handleSortChange"
      >
        <el-table-column prop="name" label="仓库名称" sortable="custom" min-width="180">
          <template #default="{ row }">
            <el-link type="primary" @click="viewRepositoryDetail(row.id)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="moduleCount" label="模块数量" sortable="custom" width="120" align="center" />
        <el-table-column prop="size" label="大小" sortable="custom" width="120" align="center" />
        <el-table-column prop="lastUpdated" label="最后更新" sortable="custom" width="180" align="center">
          <template #default="{ row }">
            {{ formatDate(row.lastUpdated) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" type="primary" @click="viewRepositoryDetail(row.id)">
                <el-icon><View /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click="confirmDelete(row)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="repositoryStore.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRepositoryStore } from '@/stores/repository'
import { ElMessageBox, ElMessage } from 'element-plus'
import { formatDate as formatDateUtil } from '@/utils/format'

const repositoryStore = useRepositoryStore()

// 分页和搜索状态
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')

// 初始化加载数据
onMounted(() => {
  loadData()
})

// 加载数据
async function loadData() {
  await Promise.all([
    repositoryStore.fetchRepositories(),
    repositoryStore.fetchRepositoryStats()
  ])
}

// 刷新列表
function refreshList() {
  repositoryStore.fetchRepositories()
}

// 刷新统计信息
function refreshStats() {
  repositoryStore.fetchRepositoryStats()
}

// 处理搜索
function handleSearch() {
  repositoryStore.updateSearch(searchQuery.value)
}

// 处理排序变更
function handleSortChange({ prop, order }: { prop: string; order: string }) {
  if (!prop) return
  
  const sortOrder = order === 'ascending' ? 'asc' : 'desc'
  repositoryStore.updateSort(prop, sortOrder)
}

// 处理分页变更
function handleSizeChange(val: number) {
  pageSize.value = val
  repositoryStore.updatePagination(currentPage.value, pageSize.value)
}

function handleCurrentChange(val: number) {
  currentPage.value = val
  repositoryStore.updatePagination(currentPage.value, pageSize.value)
}

// 查看仓库详情
function viewRepositoryDetail(id: string) {
  ElMessage.info(`查看仓库详情: ${id}，此功能尚未实现`)
  // 可以跳转到详情页面
  // router.push(`/repositories/${id}`)
}

// 确认删除
function confirmDelete(row: any) {
  ElMessageBox.confirm(
    `确定要删除仓库 "${row.name}" 吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
    .then(() => {
      repositoryStore.removeRepository(row.id)
    })
    .catch(() => {
      // 取消删除
    })
}

// 格式化日期
function formatDate(date: string) {
  return formatDateUtil(date)
}

// 监听搜索查询变化
watch(searchQuery, (newVal, oldVal) => {
  if (newVal === '' && oldVal !== '') {
    handleSearch()
  }
})
</script>

<style lang="scss" scoped>
.repository-list-page {
  .stats-overview {
    margin-bottom: 20px;
  }
  
  .stats-item {
    display: flex;
    align-items: center;
    padding: 20px;
    background-color: #f8f9fa;
    border-radius: 4px;
    margin-bottom: 10px;
    
    .stats-icon {
      font-size: 24px;
      color: #409eff;
      margin-right: 15px;
      
      .el-icon {
        font-size: 24px;
      }
    }
    
    .stats-info {
      flex: 1;
    }
    
    .stats-value {
      font-size: 24px;
      font-weight: bold;
      color: #333;
      line-height: 1.2;
    }
    
    .stats-label {
      font-size: 14px;
      color: #909399;
      margin-top: 5px;
    }
  }
  
  .data-table__actions {
    display: flex;
    gap: 10px;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>