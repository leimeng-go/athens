import { useState } from 'react'
import { Database, RefreshCw, CheckCircle, ChevronDown, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

interface ModuleVersion {
  version: string
  uploadDate: string
}

interface Module {
  name: string
  versions: ModuleVersion[]
}

export default function ModuleListPage() {
  const [modules] = useState<Module[]>([
    {
      name: 'github.com/gin-gonic/gin',
      versions: [
        { version: 'v1.9.1', uploadDate: '2024-01-20' },
        { version: 'v1.9.0', uploadDate: '2024-01-15' },
        { version: 'v1.8.2', uploadDate: '2024-01-10' },
      ],
    },
    {
      name: 'go.mongodb.org/mongo-driver',
      versions: [
        { version: 'v1.8.3', uploadDate: '2024-01-19' },
        { version: 'v1.8.2', uploadDate: '2024-01-12' },
      ],
    },
    {
      name: 'github.com/redis/go-redis',
      versions: [
        { version: 'v9.5.1', uploadDate: '2024-01-18' },
      ],
    },
  ])
  const [expandedModules, setExpandedModules] = useState<Set<string>>(new Set())

  const toggleModule = (moduleName: string) => {
    const newExpanded = new Set(expandedModules)
    if (newExpanded.has(moduleName)) {
      newExpanded.delete(moduleName)
    } else {
      newExpanded.add(moduleName)
    }
    setExpandedModules(newExpanded)
  }

  const totalVersions = modules.reduce((sum, mod) => sum + mod.versions.length, 0)

  return (
    <div>
      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <Card className="border-gray-200">
          <CardHeader className="pb-3">
            <CardDescription className="text-gray-500 text-sm font-normal">
              总模块数
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-baseline gap-2">
              <Database className="h-5 w-5 text-[#1890ff]" />
              <span className="text-3xl font-semibold text-gray-900">{modules.length}</span>
              <span className="text-sm text-gray-500">个</span>
            </div>
            <p className="text-xs text-gray-400 mt-2">已缓存模块 / 共 {totalVersions} 个版本</p>
          </CardContent>
        </Card>

        <Card className="border-gray-200">
          <CardHeader className="pb-3">
            <CardDescription className="text-gray-500 text-sm font-normal">
              存储类型
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-baseline gap-2">
              <span className="text-3xl font-semibold text-gray-900">MongoDB</span>
            </div>
            <p className="text-xs text-gray-400 mt-2">数据库存储</p>
          </CardContent>
        </Card>

        <Card className="border-gray-200">
          <CardHeader className="pb-3">
            <CardDescription className="text-gray-500 text-sm font-normal">
              系统状态
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-baseline gap-2">
              <CheckCircle className="h-5 w-5 text-[#52c41a]" />
              <span className="text-3xl font-semibold text-[#52c41a]">运行中</span>
            </div>
            <p className="text-xs text-gray-400 mt-2">系统正常</p>
          </CardContent>
        </Card>
      </div>

      {/* Module List Table */}
      <Card className="border-gray-200">
        <CardHeader className="border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold text-gray-900">
                模块列表
              </CardTitle>
              <CardDescription className="text-sm text-gray-500 mt-1">
                当前已缓存的 Go 模块
              </CardDescription>
            </div>
            <Button 
              variant="outline" 
              size="sm"
              className="border-gray-300 text-gray-700 hover:text-[#1890ff] hover:border-[#1890ff]"
            >
              <RefreshCw className="h-4 w-4 mr-1" />
              刷新
            </Button>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          {/* Table Header */}
          <div className="grid grid-cols-12 gap-4 px-6 py-3 bg-gray-50 border-b border-gray-200 text-sm font-medium text-gray-500">
            <div className="col-span-6">模块名称</div>
            <div className="col-span-2">版本数</div>
            <div className="col-span-3">最新版本</div>
            <div className="col-span-1 text-right">操作</div>
          </div>
          
          {/* Table Body */}
          <div>
            {modules.map((mod, index) => {
              const isExpanded = expandedModules.has(mod.name)
              const latestVersion = mod.versions[0]
              
              return (
                <div key={index} className="border-b border-gray-200">
                  {/* Module Row */}
                  <div
                    className="grid grid-cols-12 gap-4 px-6 py-4 hover:bg-gray-50 transition-colors cursor-pointer"
                    onClick={() => toggleModule(mod.name)}
                  >
                    <div className="col-span-6 flex items-center gap-2">
                      {isExpanded ? (
                        <ChevronDown className="h-4 w-4 text-gray-400" />
                      ) : (
                        <ChevronRight className="h-4 w-4 text-gray-400" />
                      )}
                      <code className="text-sm text-gray-900 font-mono">{mod.name}</code>
                    </div>
                    <div className="col-span-2">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-700">
                        {mod.versions.length} 个版本
                      </span>
                    </div>
                    <div className="col-span-3">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-[#e6f7ff] text-[#1890ff] border border-[#91d5ff]">
                        {latestVersion.version}
                      </span>
                    </div>
                    <div className="col-span-1 text-right">
                      <Button 
                        variant="link" 
                        size="sm"
                        className="text-[#1890ff] hover:text-[#40a9ff] p-0 h-auto"
                        onClick={(e) => e.stopPropagation()}
                      >
                        详情
                      </Button>
                    </div>
                  </div>

                  {/* Expanded Versions */}
                  {isExpanded && (
                    <div className="bg-gray-50">
                      {mod.versions.map((ver, vIndex) => (
                        <div
                          key={vIndex}
                          className="grid grid-cols-12 gap-4 px-6 py-3 pl-16 hover:bg-gray-100 transition-colors"
                        >
                          <div className="col-span-6 text-sm text-gray-600">
                            版本详情
                          </div>
                          <div className="col-span-2">
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-[#e6f7ff] text-[#1890ff] border border-[#91d5ff]">
                              {ver.version}
                            </span>
                          </div>
                          <div className="col-span-3 text-sm text-gray-500">
                            {ver.uploadDate}
                          </div>
                          <div className="col-span-1 text-right">
                            <Button 
                              variant="link" 
                              size="sm"
                              className="text-[#1890ff] hover:text-[#40a9ff] p-0 h-auto text-xs"
                            >
                              下载
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )
            })}
          </div>
          
          {/* Table Footer */}
          <div className="flex items-center justify-between px-6 py-3 border-t border-gray-200 bg-white">
            <div className="text-sm text-gray-500">
              共 {modules.length} 个模块，{totalVersions} 个版本
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" disabled className="border-gray-300">
                上一页
              </Button>
              <Button variant="outline" size="sm" className="border-gray-300">
                下一页
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
