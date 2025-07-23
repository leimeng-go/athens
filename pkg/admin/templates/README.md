# Go HTML模板嵌套使用指南

本目录展示了如何在Go中使用HTML模板的嵌套功能来构建可维护的Web应用模板系统。

## 📁 文件结构

```
templates/
├── layout.html              # 基础布局模板
├── header.html              # 头部组件
├── nav.html                 # 导航栏组件
├── sidebar.html             # 侧边栏组件
├── footer.html              # 页脚组件
├── styles.html              # 基础样式模板
├── dashboard-styles.html    # 仪表盘样式模板
├── repository-styles.html   # 仓库管理样式模板
├── upload-styles.html       # 上传页面样式模板
├── download-styles.html     # 下载页面样式模板
├── dashboard_nested.html    # 仪表盘页面模板（原版）
├── dashboard_refactored.html # 仪表盘页面模板（重构版）
├── repositories_refactored.html # 仓库管理页面模板
├── upload_refactored.html   # 上传页面模板
├── download_refactored.html # 下载页面模板
├── simple_page.html         # 简单页面模板
├── template_example.go      # Go代码实现示例
└── README.md               # 说明文档
```

## 🏗️ 模板架构

### 1. 基础布局模板 (layout.html)

基础布局定义了页面的整体结构，包含：
- HTML文档结构
- 样式引入区域 (`{{template "base-styles" .}}` 和 `{{block "page-styles" .}}`)
- 头部区域 (`{{block "header" .}}`)
- 导航栏区域 (`{{template "nav" .}}`)
- 主要内容区域 (`{{block "content" .}}`)
- 页脚区域 (`{{template "footer" .}}`)
- 脚本区域 (`{{block "scripts" .}}`)

基础布局使用 `{{block}}` 和 `{{template}}` 指令：

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}默认标题{{end}}</title>
    {{block "head" .}}{{end}}
</head>
<body>
    {{template "header" .}}
    {{template "nav" .}}
    
    <div class="main-container">
        {{block "sidebar" .}}{{end}}
        <main class="content">
            {{block "content" .}}默认内容{{end}}
        </main>
    </div>
    
    {{template "footer" .}}
    {{block "scripts" .}}{{end}}
</body>
</html>
```

### 2. 组件模板

每个组件都定义在独立的文件中，使用 `{{define}}` 指令：
- **header.html**: 网站标题和用户信息
- **nav.html**: 主导航栏，支持当前页面高亮
- **sidebar.html**: 侧边栏，根据当前页面显示不同内容
- **footer.html**: 页脚信息和链接

### 3. 样式模板

- **styles.html**: 基础样式，包含全局样式、布局样式等
- **dashboard-styles.html**: 仪表盘特定样式
- **repository-styles.html**: 仓库管理页面样式
- **upload-styles.html**: 上传页面样式
- **download-styles.html**: 下载页面样式

#### header.html
```html
{{define "header"}}
<header>
    <!-- 头部内容 -->
</header>
{{end}}
```

#### nav.html
```html
{{define "nav"}}
<nav>
    <!-- 导航内容 -->
</nav>
{{end}}
```

### 4. 页面模板

- **dashboard_refactored.html**: 仪表盘页面，展示统计信息和活动
- **repositories_refactored.html**: 仓库管理页面，展示模块列表和操作
- **upload_refactored.html**: 上传页面，支持多种上传方式
- **download_refactored.html**: 下载页面，支持搜索和浏览模块
- **simple_page.html**: 简单页面示例

具体页面继承基础布局并重写特定块：

```html
{{template "layout.html" .}}

{{define "title"}}自定义页面标题{{end}}

{{define "page-styles"}}
{{template "dashboard-styles" .}}
{{end}}

{{define "content"}}
<h1>页面内容</h1>
<p>这里是页面特定的内容</p>
{{end}}
```

## 🔧 Go代码实现

### 模板管理器

```go
type TemplateManager struct {
    templates map[string]*template.Template
}

func NewTemplateManager(templateDir string) *TemplateManager {
    tm := &TemplateManager{
        templates: make(map[string]*template.Template),
    }
    tm.parseTemplates(templateDir)
    return tm
}

func (tm *TemplateManager) parseTemplates(templateDir string) {
    // 基础模板文件
    baseFiles := []string{
        filepath.Join(templateDir, "layout.html"),
        filepath.Join(templateDir, "header.html"),
        filepath.Join(templateDir, "nav.html"),
        filepath.Join(templateDir, "sidebar.html"),
        filepath.Join(templateDir, "footer.html"),
    }
    
    // 解析具体页面模板
    dashboardFiles := append(baseFiles, filepath.Join(templateDir, "dashboard_nested.html"))
    tm.templates["dashboard"] = template.Must(template.ParseFiles(dashboardFiles...))
}
```

### 渲染模板

```go
func (tm *TemplateManager) RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
    tmpl, exists := tm.templates[templateName]
    if !exists {
        return fmt.Errorf("template %s not found", templateName)
    }
    
    return tmpl.ExecuteTemplate(w, "layout.html", data)
}
```

## 📋 关键概念

### 1. `{{template}}` vs `{{block}}`

- **`{{template "name" .}}`**: 包含一个已定义的模板
- **`{{block "name" .}}默认内容{{end}}`**: 定义一个可被重写的块，如果没有重写则显示默认内容

### 2. 模板继承

```html
<!-- 子模板 -->
{{template "layout.html" .}}  <!-- 继承基础布局 -->

{{define "title"}}自定义标题{{end}}  <!-- 重写title块 -->
{{define "page-styles"}}
{{template "dashboard-styles" .}}
{{end}}  <!-- 引入页面特定样式 -->
{{define "content"}}自定义内容{{end}}  <!-- 重写content块 -->
```

### 3. CSS样式分离

使用 `{{define}}` 定义样式模板，实现CSS的模块化管理：

```html
<!-- 在样式文件中定义 -->
{{define "dashboard-styles"}}
<style>
/* 仪表盘特定样式 */
</style>
{{end}}

<!-- 在页面模板中引用 -->
{{define "page-styles"}}
{{template "dashboard-styles" .}}
{{end}}
```

### 4. 条件渲染

```html
{{if eq .CurrentPage "dashboard"}}
    <!-- 仪表盘特定内容 -->
{{else if eq .CurrentPage "settings"}}
    <!-- 设置页面特定内容 -->
{{end}}
```

## 🎯 最佳实践

### 1. 文件组织
- 将基础布局和组件分离
- 每个页面一个模板文件
- 样式模板按功能模块分离
- 使用清晰的命名约定

### 2. CSS管理
- 基础样式放在 `styles.html` 中
- 页面特定样式使用独立的样式模板
- 通过 `{{template}}` 引用样式，避免重复
- 保持样式的模块化和可维护性

### 3. 数据结构
```go
type TemplateData struct {
    Title       string
    CurrentPage string
    User        User
    // 页面特定数据
}
```
- 定义清晰的数据结构
- 使用指针避免大结构体的复制
- 提供默认值和空值检查
- 支持分页、搜索等常用功能

### 4. 错误处理
```go
if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
    log.Printf("Template execution error: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
```
- 检查模板解析错误
- 处理模板执行错误
- 提供友好的错误页面
- 记录详细的错误日志

### 5. 性能优化
- 在应用启动时预解析所有模板
- 使用模板缓存避免重复解析
- 避免在请求处理中解析模板
- 合理使用模板缓存
- 考虑使用 `html/template` 的自动转义功能

## 🚀 优势

1. **代码复用**: 避免重复的HTML代码
2. **样式分离**: CSS模块化管理，提高可维护性
3. **易于维护**: 修改组件或样式只需要改一个文件
4. **结构清晰**: 模板职责分离，结构更清晰
5. **灵活性**: 可以根据需要选择性包含组件和样式
6. **类型安全**: Go的模板系统提供编译时检查
7. **性能优化**: 模板预编译和缓存机制
8. **响应式设计**: 支持现代Web开发需求

## 📝 使用示例

查看 `template_example.go` 文件了解完整的实现示例，包括：
- 模板管理器的实现
- 路由处理器的编写
- 数据传递的方法

## 🔗 相关资源

- [Go官方模板文档](https://pkg.go.dev/html/template)
- [模板语法参考](https://pkg.go.dev/text/template)
- [最佳实践指南](https://golang.org/doc/articles/wiki/)