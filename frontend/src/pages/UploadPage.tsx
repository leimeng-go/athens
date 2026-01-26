import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Upload } from 'lucide-react'

export default function UploadPage() {
  return (
    <Card className="border-gray-200">
      <CardHeader className="border-b border-gray-200">
        <CardTitle className="text-base font-semibold text-gray-900 flex items-center gap-2">
          <Upload className="h-5 w-5" />
          上传模块
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        <div className="text-center py-12">
          <Upload className="h-16 w-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">上传 Go 模块</h3>
          <p className="text-sm text-gray-500">
            将外网下载的模块包上传到内网 MongoDB
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
