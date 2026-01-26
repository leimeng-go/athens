import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Settings } from 'lucide-react'

export default function SettingsPage() {
  return (
    <Card className="border-gray-200">
      <CardHeader className="border-b border-gray-200">
        <CardTitle className="text-base font-semibold text-gray-900 flex items-center gap-2">
          <Settings className="h-5 w-5" />
          系统设置
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        <div className="text-center py-12">
          <Settings className="h-16 w-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">系统配置</h3>
          <p className="text-sm text-gray-500">
            MongoDB 连接、存储配置、同步设置等
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
