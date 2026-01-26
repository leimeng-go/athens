import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { FileText } from 'lucide-react'

export default function StatsPage() {
  return (
    <Card className="border-gray-200">
      <CardHeader className="border-b border-gray-200">
        <CardTitle className="text-base font-semibold text-gray-900 flex items-center gap-2">
          <FileText className="h-5 w-5" />
          统计分析
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        <div className="text-center py-12">
          <FileText className="h-16 w-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">数据统计</h3>
          <p className="text-sm text-gray-500">
            模块下载趋势、存储使用情况等统计信息
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
