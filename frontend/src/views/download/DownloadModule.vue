<template>
  <div class="download-module-page">
    <el-row :gutter="20">
      <el-col :xs="24" :md="24">
        <el-card class="stats-card">
          <template #header>
            <div class="stats-card__title">
              <span>下载统计</span>
              <el-button type="primary" size="small" @click="refreshStats" :loading="downloadStore.statsLoading">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          
          <el-row :gutter="20">
            <el-col :xs="24" :sm="6">
              <div class="stats-item">
                <div class="stats-icon">
                  <el-icon><Download /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ downloadStore.stats.totalDownloads.toLocaleString() }}</div>
                  <div class="stats-label">总下载量</div>
                </div>
              </div>
            </el-col>
            <el-col :xs="24" :sm="6">
              <div class="stats-item">
                <div class="stats-icon">
                  <el-icon><Calendar /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ downloadStore.stats.todayDownloads.toLocaleString() }}</div>
                  <div class="stats-label">今日下载</div>
                </div>
              </div>
            </el-col>
            <el-col :xs="24" :sm="6">
              <div class="stats-item">
                <div class="stats-icon">
                  <el-icon><DataAnalysis /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ downloadStore.stats.weeklyDownloads.toLocaleString() }}</div>
                  <div class="stats-label">本周下载</div>
                </div>
              </div>
            </el-col>
            <el-col :xs="24" :sm="6">
              <div class="stats-item">
                <div class="stats-icon">
                  <el-icon><TrendCharts /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ downloadStore.stats.monthlyDownloads.toLocaleString() }}</div>
                  <div class="stats-label">本月下载</div>
                </div>
              </div>
            </el-col>
          </el-row>
          
          <div class="download-trend-chart" v-if="downloadStore.stats.downloadTrend.length > 0">
            <div ref="chartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="demo-form-inline">
        <el-form-item label="模块路径">
          <el-input v-model="searchForm.query" placeholder="输入模块路径搜索" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="data-table">
      <div class="data-table__header">
        <div class="data-table__title">模块列表</div>
        <div class="data-table__actions">
          <el-button type="primary" @click="refreshModules">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
      
      <el-table
        v-loading="downloadStore.loading"
        :data="downloadStore.modules"
        style="width: 100%"
        border
        @sort-change="handleSortChange"
      >
        <el-table-column prop="path" label="模块路径" sortable="custom" min-width="250" show-overflow-tooltip>
          <template #default="{ row }">
            <el-link type="primary" @click="viewModuleDetail(row.path)">
              {{ row.path }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="latestVersion" label="最新版本" width="120" />
        <el-table-column prop="versions" label="版本数" sortable="custom" width="100" align="center" />
        <el-table-column prop="downloads" label="下载次数" sortable="custom" width="120" align="center" />
        <el-table-column prop="lastDownloaded" label="最后下载" sortable="custom" width="180" align="center">
          <template #default="{ row }">
            {{ formatDate(row.lastDownloaded) }}
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" sortable="custom" width="100" align="center" />
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" type="primary" @click="viewModuleDetail(row.path)">
                <el-icon><View /></el-icon>
              </el-button>
              <el-button size="small" type="success" @click="downloadModule(row)">
                <el-icon><Download /></el-icon>
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
          :total="downloadStore.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
    
    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="popular-modules-card">
          <template #header>
            <div class="card-header">
              <span>热门模块</span>
            </div>
          </template>
          <el-table
            :data="downloadStore.popularModules"
            style="width: 100%"
            :show-header="false"
          >
            <el-table-column width="50">
              <template #default="{ $index }">
                <el-tag :type="getPopularRankType($index)" effect="plain" size="large">
                  {{ $index + 1 }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="path">
              <template #default="{ row }">
                <el-link type="primary" @click="viewModuleDetail(row.path)">
                  {{ row.path }}
                </el-link>
              </template>
            </el-table-column>
            <el-table-column prop="downloads" width="100" align="right">
              <template #default="{ row }">
                <span class="download-count">{{ row.downloads.toLocaleString() }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :md="12">
        <el-card class="recent-downloads-card">
          <template #header>
            <div class="card-header">
              <span>最近下载</span>
            </div>
          </template>
          <el-table
            :data="downloadStore.recentDownloads"
            style="width: 100%"
          >
            <el-table-column prop="path" label="模块路径" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">
                <el-link type="primary" @click="viewModuleDetail(row.path)">
                  {{ row.path }}
                </el-link>
              </template>
            </el-table-column>
            <el-table-column prop="version" label="版本" width="100" />
            <el-table-column prop="timestamp" label="时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.timestamp) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted, nextTick } from 'vue'
import { useDownloadStore } from '@/stores/download'
import { ElMessage } from 'element-plus'
import { formatDate as formatDateUtil } from '@/utils/format'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

// 注册ECharts组件
echarts.use([
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  LineChart,
  CanvasRenderer
])

const downloadStore = useDownloadStore()

// 图表引用
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

// 搜索表单
const searchForm = ref({
  query: ''
})

// 分页状态
const currentPage = ref(1)
const pageSize = ref(10)

// 初始化加载数据
onMounted(async () => {
  await Promise.all([
    downloadStore.fetchModules(),
    downloadStore.fetchDownloadStats(),
    downloadStore.fetchPopularModules(),
    downloadStore.fetchRecentDownloads()
  ])
  
  nextTick(() => {
    initChart()
  })
})

// 初始化图表
function initChart() {
  if (!chartRef.value) return
  
  chart = echarts.init(chartRef.value)
  updateChart()
  
  // 监听窗口大小变化，调整图表大小
  window.addEventListener('resize', handleResize)
}

// 更新图表数据
function updateChart() {
  if (!chart) return
  
  const { downloadTrend } = downloadStore.stats
  
  const option = {
    title: {
      text: '下载趋势',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: downloadTrend.map(item => item.date)
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '下载量',
        type: 'line',
        data: downloadTrend.map(item => item.count),
        areaStyle: {},
        smooth: true
      }
    ]
  }
  
  chart.setOption(option)
}

// 处理窗口大小变化
function handleResize() {
  chart?.resize()
}

// 组件卸载时清理
onUnmounted(() => {
  if (chart) {
    chart.dispose()
    chart = null
  }
  window.removeEventListener('resize', handleResize)
})

// 刷新模块列表
function refreshModules() {
  downloadStore.fetchModules()
}

// 刷新统计信息
function refreshStats() {
  downloadStore.fetchDownloadStats().then(() => {
    updateChart()
  })
}

// 处理搜索
function handleSearch() {
  downloadStore.updateSearch(searchForm.value.query)
}

// 重置搜索
function resetSearch() {
  searchForm.value.query = ''
  downloadStore.updateSearch('')
}

// 处理排序变更
function handleSortChange({ prop, order }: { prop: string; order: string }) {
  if (!prop) return
  
  const sortOrder = order === 'ascending' ? 'asc' : 'desc'
  downloadStore.updateSort(prop, sortOrder)
}

// 处理分页变更
function handleSizeChange(val: number) {
  pageSize.value = val
  downloadStore.updatePagination(currentPage.value, pageSize.value)
}

function handleCurrentChange(val: number) {
  currentPage.value = val
  downloadStore.updatePagination(currentPage.value, pageSize.value)
}

// 查看模块详情
function viewModuleDetail(path: string) {
  ElMessage.info(`查看模块详情: ${path}，此功能尚未实现`)
  // 可以跳转到详情页面
  // router.push(`/download/modules/${encodeURIComponent(path)}`)
}

// 下载模块
function downloadModule(row: any) {
  const url = `/download/module/${encodeURIComponent(row.path)}@${row.latestVersion}`
  ElMessage.success(`开始下载模块: ${row.path}@${row.latestVersion}`)
  
  // 创建一个隐藏的a标签并模拟点击下载
  const link = document.createElement('a')
  link.href = url
  link.target = '_blank'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 格式化日期
function formatDate(date: string) {
  return formatDateUtil(date)
}

// 获取热门排名标签类型
function getPopularRankType(index: number) {
  if (index === 0) return 'danger'
  if (index === 1) return 'warning'
  if (index === 2) return 'success'
  return 'info'
}

// 监听统计数据变化，更新图表
watch(
  () => downloadStore.stats,
  () => {
    nextTick(() => {
      updateChart()
    })
  },
  { deep: true }
)
</script>

<style lang="scss" scoped>
.download-module-page {
  .stats-card {
    margin-bottom: 20px;
  }
  
  .stats-item {
    display: flex;
    align-items: center;
    padding: 15px;
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
  
  .download-trend-chart {
    margin-top: 20px;
  }
  
  .search-form {
    margin-bottom: 20px;
    padding: 20px;
    background-color: #fff;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
  
  .data-table {
    margin-bottom: 20px;
  }
  
  .data-table__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
  
  .popular-modules-card,
  .recent-downloads-card {
    margin-bottom: 20px;
    height: 100%;
  }
  
  .download-count {
    font-weight: bold;
    color: #409eff;
  }
}
</style>