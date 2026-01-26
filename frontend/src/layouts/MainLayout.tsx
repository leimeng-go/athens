import { useState } from 'react'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { Package, Database, Upload, FileText, Settings } from 'lucide-react'

export default function MainLayout() {
  const location = useLocation()
  const [collapsed, setCollapsed] = useState(false)

  const menuItems = [
    { path: '/', icon: Database, label: '模块列表' },
    { path: '/upload', icon: Upload, label: '上传模块' },
    { path: '/stats', icon: FileText, label: '统计分析' },
    { path: '/settings', icon: Settings, label: '系统设置' },
  ]

  return (
    <div className="flex h-screen bg-[#f0f2f5]">
      {/* Sidebar */}
      <aside
        className={`${
          collapsed ? 'w-20' : 'w-56'
        } bg-[#001529] transition-all duration-300 flex flex-col`}
      >
        {/* Logo */}
        <div className="h-16 flex items-center justify-center border-b border-gray-700">
          <div className="flex items-center gap-2">
            <Package className="h-8 w-8 text-[#1890ff]" />
            {!collapsed && (
              <span className="text-white font-semibold text-lg">Athens</span>
            )}
          </div>
        </div>

        {/* Menu */}
        <nav className="flex-1 py-4">
          {menuItems.map((item) => {
            const Icon = item.icon
            const isActive = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`
                  flex items-center gap-3 px-6 py-3 text-sm transition-colors
                  ${
                    isActive
                      ? 'bg-[#1890ff] text-white'
                      : 'text-gray-300 hover:bg-gray-700'
                  }
                `}
              >
                <Icon className="h-5 w-5" />
                {!collapsed && <span>{item.label}</span>}
              </Link>
            )
          })}
        </nav>

        {/* Collapse Button */}
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="h-12 text-gray-400 hover:text-white border-t border-gray-700"
        >
          {collapsed ? '→' : '←'}
        </button>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6">
          <div>
            <h1 className="text-lg font-semibold text-gray-900">
              Athens 依赖管理
            </h1>
            <p className="text-xs text-gray-500">
              Go 模块代理 - 内网依赖管理平台
            </p>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
