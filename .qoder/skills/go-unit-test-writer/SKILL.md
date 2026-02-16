---
name: go-unit-test-writer
description: 编写高质量的 Go 语言单元测试，包括测试结构设计、表格驱动测试、mock 对象创建和测试覆盖率优化。使用于编写 Go 项目的单元测试用例时。
---

# Go 单元测试编写助手

## 核心原则

1. **测试金字塔**：单元测试 > 集成测试 > 端到端测试
2. **FIRST 原则**：Fast, Independent, Repeatable, Self-validating, Timely
3. **覆盖率目标**：核心业务逻辑 80%+，工具函数 100%

## 测试结构模板

### 基础测试结构
```go
func TestFunctionName(t *testing.T) {
    // 准备测试数据
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "正常情况测试",
            input:    validInput,
            expected: expectedResult,
            wantErr:  false,
        },
        {
            name:     "边界条件测试",
            input:    edgeCaseInput,
            expected: edgeCaseResult,
            wantErr:  false,
        },
        {
            name:     "错误情况测试",
            input:    invalidInput,
            expected: zeroValue,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 执行被测试函数
            got, err := FunctionName(tt.input)
            
            // 断言错误处理
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            
            // 断言结果
            assert.Equal(t, tt.expected, got)
        })
    }
}
```

## 常用测试模式

### 1. 表格驱动测试
```go
func TestValidateUser(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {
            name: "valid user",
            user: User{Name: "John", Email: "john@example.com", Age: 25},
            wantErr: false,
        },
        {
            name: "invalid email",
            user: User{Name: "John", Email: "invalid-email", Age: 25},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUser(tt.user)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 2. Mock 对象创建
```go
// 使用 testify/mock
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetUser(id string) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

// 在测试中使用
func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewUserService(mockRepo)
    
    // 设置期望
    expectedUser := &User{ID: "123", Name: "John"}
    mockRepo.On("GetUser", "123").Return(expectedUser, nil)
    
    // 执行测试
    user, err := service.GetUser("123")
    
    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

### 3. HTTP Handler 测试
```go
func TestGetUserHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    
    // 创建 mock 服务
    mockService := new(MockUserService)
    mockService.On("GetUser", "123").Return(&User{ID: "123", Name: "John"}, nil)
    
    handler := GetUserHandler(mockService)
    handler.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    // 验证响应内容
    var response UserResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "John", response.Name)
}
```

## 断言最佳实践

### 基本断言
```go
// 相等性断言
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)

// 布尔断言
assert.True(t, condition)
assert.False(t, condition)

// 错误断言
assert.NoError(t, err)
assert.Error(t, err)
assert.ErrorIs(t, err, expectedErr)

// Nil 断言
assert.Nil(t, value)
assert.NotNil(t, value)

// 切片/映射断言
assert.Len(t, slice, expectedLength)
assert.Contains(t, slice, element)
assert.NotContains(t, slice, element)
```

### 高级断言
```go
// 时间断言
assert.WithinDuration(t, expectedTime, actualTime, time.Second)

// 正则表达式断言
assert.Regexp(t, regexp.MustCompile(`pattern`), actualString)

// HTTP 状态码断言
assert.HTTPStatusCode(t, handler, "GET", "/path", nil, http.StatusOK)

// JSON 断言
assert.JSONEq(t, `{"key": "value"}`, response.Body.String())
```

## 测试工具推荐

### 1. 测试框架
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)
```

### 2. 性能测试
```go
func BenchmarkProcessData(b *testing.B) {
    data := generateTestData()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ProcessData(data)
    }
}
```

### 3. 测试覆盖率
```bash
# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 常见测试场景

### 1. 数据库操作测试
```go
func TestRepository_CreateUser(t *testing.T) {
    // 使用事务回滚避免污染测试数据
    db := setupTestDB()
    defer db.Rollback()
    
    repo := NewRepository(db)
    user := &User{Name: "John", Email: "john@example.com"}
    
    createdUser, err := repo.CreateUser(user)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, createdUser.ID)
    assert.Equal(t, user.Name, createdUser.Name)
}
```

### 2. 并发测试
```go
func TestConcurrentAccess(t *testing.T) {
    counter := &Counter{}
    var wg sync.WaitGroup
    
    // 并发执行
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    assert.Equal(t, 100, counter.Value())
}
```

### 3. 环境变量测试
```go
func TestWithEnvironmentVariables(t *testing.T) {
    // 保存原始值
    originalValue := os.Getenv("TEST_ENV")
    defer os.Setenv("TEST_ENV", originalValue)
    
    // 设置测试值
    os.Setenv("TEST_ENV", "test_value")
    
    // 执行测试
    result := GetConfigValue()
    assert.Equal(t, "test_value", result)
}
```

## 测试质量检查清单

- [ ] 测试函数命名清晰描述测试内容
- [ ] 每个测试用例都有明确的测试目的
- [ ] 覆盖正常路径、边界条件和错误路径
- [ ] 使用表格驱动测试提高可维护性
- [ ] 适当的测试数据隔离
- [ ] 清理测试产生的资源
- [ ] 测试运行速度快（单元测试 < 1秒）
- [ ] 测试结果稳定可重现

## 参考资料

- [Go 官方测试文档](https://golang.org/pkg/testing/)
- [Testify 断言库](https://github.com/stretchr/testify)
- [Go 测试最佳实践](https://github.com/golang/go/wiki/TestComments)