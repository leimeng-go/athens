# Go 测试最佳实践参考

## 测试命名规范

### 测试函数命名
```go
// 好的命名 - 清晰描述测试内容
func TestUserService_CreateUser(t *testing.T)
func TestValidateEmail_ValidFormat(t *testing.T)
func TestCalculateTax_BoundaryConditions(t *testing.T)

// 避免的命名 - 过于宽泛或不清晰
func TestCreate(t *testing.T)
func TestValidation(t *testing.T)
func Test1(t *testing.T)
```

### 子测试命名
```go
func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name string
        // ...
    }{
        {
            name: "should_return_user_when_user_exists",
            // ...
        },
        {
            name: "should_return_error_when_user_not_found", 
            // ...
        },
        {
            name: "should_handle_database_connection_failure",
            // ...
        },
    }
}
```

## 测试数据组织

### 测试用例结构
```go
tests := []struct {
    name        string        // 测试用例名称
    setup       func()        // 测试前置条件（可选）
    input       interface{}   // 输入数据
    expected    interface{}   // 期望输出
    wantErr     bool          // 是否期望错误
    assertFunc  func()        // 自定义断言函数（可选）
}{
    // 测试用例...
}
```

### 测试数据工厂函数
```go
func newUser(name, email string) *User {
    return &User{
        ID:        uuid.New().String(),
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
}

func newValidUser() *User {
    return newUser("John Doe", "john@example.com")
}

func newInvalidUser() *User {
    return newUser("", "invalid-email")
}

// 在测试中使用
func TestUserService_CreateUser(t *testing.T) {
    validUser := newValidUser()
    invalidUser := newInvalidUser()
    // ...
}
```

## Mock 对象最佳实践

### 接口设计便于测试
```go
// 好的设计 - 明确的接口
type EmailService interface {
    SendWelcomeEmail(user *User) error
    SendPasswordResetEmail(email string, token string) error
}

// 避免的设计 - 过于宽泛的接口
type Service interface {
    DoEverything() (interface{}, error)
}
```

### Mock 实现技巧
```go
// 使用 testify/mock
type MockEmailService struct {
    mock.Mock
}

func (m *MockEmailService) SendWelcomeEmail(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(email, token string) error {
    args := m.Called(email, token)
    return args.Error(0)
}

// 手动实现 mock（适用于简单场景）
type fakeEmailService struct {
    sentEmails []EmailRecord
    shouldFail bool
}

type EmailRecord struct {
    Type  string
    Email string
    Data  map[string]interface{}
}

func (f *fakeEmailService) SendWelcomeEmail(user *User) error {
    if f.shouldFail {
        return errors.New("email service unavailable")
    }
    f.sentEmails = append(f.sentEmails, EmailRecord{
        Type:  "welcome",
        Email: user.Email,
        Data:  map[string]interface{}{"name": user.Name},
    })
    return nil
}
```

## 测试环境管理

### 测试数据库设置
```go
func setupTestDatabase(t *testing.T) *sql.DB {
    // 使用测试数据库
    db, err := sql.Open("postgres", "postgres://test:test@localhost/test_db?sslmode=disable")
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }
    
    // 创建测试表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        t.Fatalf("failed to create test tables: %v", err)
    }
    
    return db
}

func cleanupTestDatabase(t *testing.T, db *sql.DB) {
    // 清理测试数据
    _, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
    if err != nil {
        t.Errorf("failed to cleanup test data: %v", err)
    }
    
    db.Close()
}

// 在测试中使用
func TestUserRepository(t *testing.T) {
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    repo := NewUserRepository(db)
    // ... 测试代码
}
```

### 临时文件处理
```go
func createTempFile(t *testing.T, content string) string {
    tmpFile, err := ioutil.TempFile("", "test-*.txt")
    if err != nil {
        t.Fatalf("failed to create temp file: %v", err)
    }
    defer tmpFile.Close()
    
    _, err = tmpFile.WriteString(content)
    if err != nil {
        t.Fatalf("failed to write to temp file: %v", err)
    }
    
    return tmpFile.Name()
}

func TestFileProcessor(t *testing.T) {
    tempFileName := createTempFile(t, "test content")
    defer os.Remove(tempFileName) // 清理临时文件
    
    processor := NewFileProcessor()
    result, err := processor.ProcessFile(tempFileName)
    
    assert.NoError(t, err)
    assert.Equal(t, "processed: test content", result)
}
```

## 测试覆盖率优化

### 覆盖率测量
```bash
# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看详细覆盖率信息
go tool cover -func=coverage.out

# 生成 HTML 覆盖率报告
go tool cover -html=coverage.out -o coverage.html

# 设置覆盖率阈值
go test -coverprofile=coverage.out ./...
echo "Coverage: $(go tool cover -func=coverage.out | grep total | awk '{print $3}')"
```

### 提高覆盖率的策略
```go
// 1. 测试边界条件
func TestDivide(t *testing.T) {
    tests := []struct {
        name string
        a, b float64
        want float64
        wantErr bool
    }{
        {"normal division", 10, 2, 5, false},
        {"division by zero", 10, 0, 0, true},
        {"negative numbers", -10, 2, -5, false},
        {"zero dividend", 0, 5, 0, false},
        {"both negative", -10, -2, 5, false},
    }
    // ...
}

// 2. 测试错误路径
func TestRepository_GetUser(t *testing.T) {
    tests := []struct {
        name string
        userID string
        mockError error
        expectError bool
    }{
        {"success case", "123", nil, false},
        {"not found", "456", ErrUserNotFound, true},
        {"database error", "789", errors.New("connection failed"), true},
    }
    // ...
}

// 3. 测试并发场景
func TestConcurrentAccess(t *testing.T) {
    const goroutines = 100
    const operations = 1000
    
    service := NewConcurrentService()
    var wg sync.WaitGroup
    results := make(chan int, goroutines*operations)
    
    for i := 0; i < goroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < operations; j++ {
                result := service.Process(id*j)
                results <- result
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
    
    // 验证所有操作都正确执行
    processed := 0
    for range results {
        processed++
    }
    assert.Equal(t, goroutines*operations, processed)
}
```

## 测试调试技巧

### 测试失败时的信息输出
```go
func TestComplexCalculation(t *testing.T) {
    input := generateComplexInput()
    expected := calculateExpectedResult(input)
    actual := performCalculation(input)
    
    // 详细的错误信息
    if actual != expected {
        t.Errorf("calculation failed for input %+v\nexpected: %v\ngot: %v\ndifference: %v", 
            input, expected, actual, math.Abs(actual-expected))
    }
}
```

### 条件测试跳过
```go
func TestExternalServiceIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    if os.Getenv("INTEGRATION_TESTS") != "true" {
        t.Skip("integration tests not enabled")
    }
    
    // 实际的集成测试代码
}

func TestDatabaseConnection(t *testing.T) {
    dbURL := os.Getenv("TEST_DATABASE_URL")
    if dbURL == "" {
        t.Skip("TEST_DATABASE_URL not set, skipping database tests")
    }
    
    // 数据库测试代码
}
```

这些最佳实践可以帮助你编写更高质量、更易维护的 Go 单元测试。