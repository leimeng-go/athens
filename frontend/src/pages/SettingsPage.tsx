import { configAPI } from '@/lib/api'
import { CheckCircleOutlined, SaveOutlined, SettingOutlined, SyncOutlined } from '@ant-design/icons'
import { Alert, Button, Card, Col, Form, Input, message, Row, Select, Space, Spin, Switch } from 'antd'
import { useEffect, useState } from 'react'

const { Option } = Select

interface ConfigFormData {
  go_env: string
  go_binary: string
  go_get_workers: number
  protocol_workers: number
  log_level: string
  cloud_runtime: string
  enable_pprof: boolean
  pprof_port: string
  storage_type: string
  port: string
  basic_auth_user: string
  basic_auth_pass: string
  force_ssl: boolean
  network_mode: string
  single_flight_type: string
  index_type: string
  shutdown_timeout: number
}

export default function SettingsPage() {
  const [form] = Form.useForm()
  const [config, setConfig] = useState<ConfigFormData | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [resetting, setResetting] = useState(false)

  // 加载配置
  useEffect(() => {
    loadConfig()
  }, [])

  const loadConfig = async () => {
    try {
      setLoading(true)
      const configData = await configAPI.get()
      
      // 将后端配置映射到表单数据
      const formData: ConfigFormData = {
        go_env: configData.GoEnv || 'development',
        go_binary: configData.GoBinary || 'go',
        go_get_workers: configData.GoGetWorkers || 10,
        protocol_workers: configData.ProtocolWorkers || 30,
        log_level: configData.LogLevel || 'debug',
        cloud_runtime: configData.CloudRuntime || 'none',
        enable_pprof: configData.EnablePprof || false,
        pprof_port: configData.PprofPort || ':3001',
        storage_type: configData.StorageType || 'memory',
        port: configData.Port || ':3000',
        basic_auth_user: configData.BasicAuthUser || '',
        basic_auth_pass: configData.BasicAuthPass || '',
        force_ssl: configData.ForceSSL || false,
        network_mode: configData.NetworkMode || 'strict',
        single_flight_type: configData.SingleFlightType || 'memory',
        index_type: configData.IndexType || 'none',
        shutdown_timeout: configData.ShutdownTimeout || 60,
      }
      
      setConfig(formData)
      form.setFieldsValue(formData)
    } catch (error) {
      message.error('加载配置失败: ' + (error as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      setSaving(true)
      const values = await form.validateFields()
      
      // 构造更新请求
      const updateRequest: any = {
        GoEnv: values.go_env,
        GoBinary: values.go_binary,
        GoGetWorkers: values.go_get_workers,
        ProtocolWorkers: values.protocol_workers,
        LogLevel: values.log_level,
        CloudRuntime: values.cloud_runtime,
        EnablePprof: values.enable_pprof,
        PprofPort: values.pprof_port,
        StorageType: values.storage_type,
        Port: values.port,
        BasicAuthUser: values.basic_auth_user || '',
        BasicAuthPass: values.basic_auth_pass || '',
        ForceSSL: values.force_ssl,
        NetworkMode: values.network_mode,
        SingleFlightType: values.single_flight_type,
        IndexType: values.index_type,
        ShutdownTimeout: values.shutdown_timeout,
      }
      
      await configAPI.update(updateRequest)
      message.success('配置保存成功')
      loadConfig() // 重新加载配置
    } catch (error) {
      message.error('保存配置失败: ' + (error as Error).message)
    } finally {
      setSaving(false)
    }
  }

  const handleReset = async () => {
    try {
      setResetting(true)
      await configAPI.reset()
      message.success('配置已重置为默认值')
      loadConfig() // 重新加载配置
    } catch (error) {
      message.error('重置配置失败: ' + (error as Error).message)
    } finally {
      setResetting(false)
    }
  }

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '400px' }}>
        <Spin size="large" tip="加载配置中..." />
      </div>
    )
  }

  if (!config) {
    return (
      <Card>
        <Alert
          message="无法加载配置"
          description="请检查后端服务是否正常运行"
          type="error"
          showIcon
          action={
            <Button onClick={loadConfig} type="primary">
              重新加载
            </Button>
          }
        />
      </Card>
    )
  }

  return (
    <div style={{ padding: '24px' }}>
      {/* 页面标题 */}
      <div style={{ marginBottom: '24px' }}>
        <h1 style={{ fontSize: '24px', fontWeight: 'bold', color: '#1d1d1d', marginBottom: '8px' }}>
          <SettingOutlined style={{ marginRight: '8px' }} />
          系统设置
        </h1>
        <p style={{ color: '#8c8c8c', margin: 0 }}>
          配置 Athens 模块代理的各项参数
        </p>
      </div>

      {/* 操作按钮 */}
      <div style={{ marginBottom: '24px', textAlign: 'right' }}>
        <Space>
          <Button 
            onClick={handleReset} 
            icon={<SyncOutlined spin={resetting} />}
            disabled={resetting}
          >
            重置默认
          </Button>
          <Button 
            type="primary" 
            onClick={handleSave} 
            icon={<SaveOutlined />}
            loading={saving}
          >
            保存配置
          </Button>
        </Space>
      </div>

      <Form
        form={form}
        layout="vertical"
        initialValues={config}
      >
        <Row gutter={[24, 24]}>
          {/* 基础配置 */}
          <Col span={12}>
            <Card title="基础配置" size="small">
              <Form.Item
                name="go_env"
                label="运行环境"
                rules={[{ required: true, message: '请选择运行环境' }]}
              >
                <Select>
                  <Option value="development">Development</Option>
                  <Option value="production">Production</Option>
                </Select>
              </Form.Item>

              <Form.Item
                name="go_binary"
                label="Go 二进制路径"
                rules={[{ required: true, message: '请输入Go二进制路径' }]}
              >
                <Input placeholder="go" />
              </Form.Item>

              <Form.Item
                name="log_level"
                label="日志级别"
                rules={[{ required: true, message: '请选择日志级别' }]}
              >
                <Select>
                  <Option value="debug">Debug</Option>
                  <Option value="info">Info</Option>
                  <Option value="warn">Warning</Option>
                  <Option value="error">Error</Option>
                </Select>
              </Form.Item>

              <Form.Item
                name="port"
                label="服务端口"
                rules={[{ required: true, message: '请输入服务端口' }]}
              >
                <Input placeholder=":3000" />
              </Form.Item>
            </Card>
          </Col>

          {/* 性能配置 */}
          <Col span={12}>
            <Card title="性能配置" size="small">
              <Form.Item
                name="go_get_workers"
                label="Go Get 工作协程数"
                rules={[{ required: true, message: '请输入工作协程数' }]}
              >
                <Input type="number" min="1" />
              </Form.Item>

              <Form.Item
                name="protocol_workers"
                label="协议工作协程数"
                rules={[{ required: true, message: '请输入协议工作协程数' }]}
              >
                <Input type="number" min="1" />
              </Form.Item>

              <Form.Item
                name="shutdown_timeout"
                label="关闭超时(秒)"
                rules={[{ required: true, message: '请输入关闭超时时间' }]}
              >
                <Input type="number" min="0" />
              </Form.Item>
            </Card>
          </Col>

          {/* 网络配置 */}
          <Col span={12}>
            <Card title="网络配置" size="small">
              <Form.Item
                name="network_mode"
                label="网络模式"
                rules={[{ required: true, message: '请选择网络模式' }]}
              >
                <Select>
                  <Option value="strict">Strict</Option>
                  <Option value="offline">Offline</Option>
                  <Option value="fallback">Fallback</Option>
                </Select>
              </Form.Item>

              <Form.Item
                name="force_ssl"
                label="强制 SSL 重定向"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>

              <Form.Item
                name="pprof_port"
                label="性能分析端口"
                rules={[{ required: true, message: '请输入性能分析端口' }]}
              >
                <Input placeholder=":3001" />
              </Form.Item>

              <Form.Item
                name="enable_pprof"
                label="启用性能分析"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </Card>
          </Col>

          {/* 安全配置 */}
          <Col span={12}>
            <Card title="安全配置" size="small">
              <Form.Item
                name="basic_auth_user"
                label="基本认证用户名"
              >
                <Input placeholder="admin" />
              </Form.Item>

              <Form.Item
                name="basic_auth_pass"
                label="基本认证密码"
              >
                <Input.Password placeholder="password" />
              </Form.Item>

              <Form.Item
                name="cloud_runtime"
                label="云运行时"
                rules={[{ required: true, message: '请选择云运行时' }]}
              >
                <Select>
                  <Option value="none">None</Option>
                  <Option value="GCP">GCP</Option>
                  <Option value="AWS">AWS</Option>
                </Select>
              </Form.Item>
            </Card>
          </Col>

          {/* 存储配置 */}
          <Col span={12}>
            <Card title="存储配置" size="small">
              <Form.Item
                name="storage_type"
                label="存储类型"
                rules={[{ required: true, message: '请选择存储类型' }]}
              >
                <Select>
                  <Option value="memory">内存存储</Option>
                  <Option value="mongo">MongoDB</Option>
                  <Option value="disk">磁盘存储</Option>
                </Select>
              </Form.Item>

              <Form.Item
                name="index_type"
                label="索引类型"
                rules={[{ required: true, message: '请选择索引类型' }]}
              >
                <Select>
                  <Option value="none">无索引</Option>
                  <Option value="memory">内存索引</Option>
                  <Option value="mysql">MySQL</Option>
                  <Option value="postgres">PostgreSQL</Option>
                </Select>
              </Form.Item>

              <Form.Item
                name="single_flight_type"
                label="并发控制类型"
                rules={[{ required: true, message: '请选择并发控制类型' }]}
              >
                <Select>
                  <Option value="memory">内存控制</Option>
                  <Option value="redis">Redis</Option>
                  <Option value="etcd">Etcd</Option>
                </Select>
              </Form.Item>
            </Card>
          </Col>

          {/* 状态信息 */}
          <Col span={12}>
            <Card title="配置状态" size="small">
              <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
                <CheckCircleOutlined style={{ color: '#52c41a', marginRight: '8px' }} />
                <span>配置已加载</span>
              </div>
              <div style={{ fontSize: '12px', color: '#8c8c8c' }}>
                最后更新时间: {new Date().toLocaleString()}
              </div>
              <div style={{ fontSize: '12px', color: '#8c8c8c', marginTop: '4px' }}>
                配置文件路径: /config/athens.toml
              </div>
            </Card>
          </Col>
        </Row>
      </Form>
    </div>
  )
}