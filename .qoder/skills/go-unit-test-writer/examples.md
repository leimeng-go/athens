# Go 单元测试编写示例

## 基础示例

### 简单函数测试
```go
// 被测试函数
func Add(a, b int) int {
    return a + b
}

// 测试函数
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"正数相加", 2, 3, 5},
        {"负数相加", -1, -1, -2},
        {"零值测试", 0, 5, 5},
        {"边界值测试", math.MaxInt32, 1, math.MinInt32},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 结构体方法测试
```go
type Calculator struct {
    precision int
}

func (c *Calculator) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return math.Round((a/b)*math.Pow10(c.precision)) / math.Pow10(c.precision), nil
}

func TestCalculator_Divide(t *testing.T) {
    calc := &Calculator{precision: 2}
    
    tests := []struct {
        name        string
        a, b        float64
        expected    float64
        expectError bool
    }{
        {"正常除法", 10, 3, 3.33, false},
        {"除零错误", 10, 0, 0, true},
        {"负数除法", -10, 2, -5.00, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calc.Divide(tt.a, tt.b)
            
            if tt.expectError {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 高级示例

### 依赖注入测试
```go
// 接口定义
type UserRepository interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}

// 实现类
type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
    return &userService{repo: repo}
}

func (s *userService) GetUserProfile(userID string) (*UserProfile, error) {
    user, err := s.repo.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    return &UserProfile{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

// 测试代码
func TestUserService_GetUserProfile(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // 测试成功情况
    t.Run("successful retrieval", func(t *testing.T) {
        expectedUser := &User{
            ID:    "123",
            Name:  "John Doe",
            Email: "john@example.com",
        }
        
        mockRepo.On("GetUser", "123").Return(expectedUser, nil)
        
        profile, err := service.GetUserProfile("123")
        
        assert.NoError(t, err)
        assert.Equal(t, "123", profile.ID)
        assert.Equal(t, "John Doe", profile.Name)
        assert.Equal(t, "john@example.com", profile.Email)
        
        mockRepo.AssertExpectations(t)
    })
    
    // 测试错误情况
    t.Run("repository error", func(t *testing.T) {
        mockRepo.On("GetUser", "456").Return((*User)(nil), errors.New("database error"))
        
        profile, err := service.GetUserProfile("456")
        
        assert.Error(t, err)
        assert.Nil(t, profile)
        assert.Contains(t, err.Error(), "database error")
        
        mockRepo.AssertExpectations(t)
    })
}
```

### HTTP 处理器测试
```go
type UserHandler struct {
    service UserService
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    if userID == "" {
        http.Error(w, "missing user ID", http.StatusBadRequest)
        return
    }
    
    user, err := h.service.GetUserProfile(userID)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            http.Error(w, "user not found", http.StatusNotFound)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func TestUserHandler_GetUser(t *testing.T) {
    tests := []struct {
        name           string
        userID         string
        mockUser       *UserProfile
        mockError      error
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "successful retrieval",
            userID: "123",
            mockUser: &UserProfile{
                ID:    "123",
                Name:  "John Doe",
                Email: "john@example.com",
            },
            mockError:      nil,
            expectedStatus: http.StatusOK,
            expectedBody:   `{"ID":"123","Name":"John Doe","Email":"john@example.com"}`,
        },
        {
            name:           "missing user ID",
            userID:         "",
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "missing user ID\n",
        },
        {
            name:           "user not found",
            userID:         "456",
            mockUser:       nil,
            mockError:      ErrUserNotFound,
            expectedStatus: http.StatusNotFound,
            expectedBody:   "user not found\n",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := &MockUserService{}
            if tt.userID != "" && tt.mockError == nil {
                mockService.On("GetUserProfile", tt.userID).Return(tt.mockUser, tt.mockError)
            }
            
            handler := &UserHandler{service: mockService}
            
            req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s", tt.userID), nil)
            w := httptest.NewRecorder()
            
            handler.GetUser(w, req)
            
            assert.Equal(t, tt.expectedStatus, w.Code)
            if tt.expectedBody != "" {
                assert.Equal(t, tt.expectedBody, w.Body.String())
            }
            
            if tt.userID != "" && tt.mockError == nil {
                mockService.AssertExpectations(t)
            }
        })
    }
}
```

### 数据库测试
```go
func TestUserRepository_CreateUser(t *testing.T) {
    // 使用测试数据库
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    repo := NewUserRepository(db)
    
    tests := []struct {
        name     string
        user     *User
        wantErr  bool
    }{
        {
            name: "valid user creation",
            user: &User{
                Name:  "John Doe",
                Email: "john@example.com",
                Age:   30,
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            user: &User{
                Name:  "Jane Doe",
                Email: "john@example.com", // 重复邮箱
                Age:   25,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            createdUser, err := repo.CreateUser(tt.user)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, createdUser)
                return
            }
            
            assert.NoError(t, err)
            assert.NotEmpty(t, createdUser.ID)
            assert.Equal(t, tt.user.Name, createdUser.Name)
            assert.Equal(t, tt.user.Email, createdUser.Email)
            assert.WithinDuration(t, time.Now(), createdUser.CreatedAt, time.Second)
            
            // 验证数据确实存储到了数据库
            storedUser, err := repo.GetUser(createdUser.ID)
            assert.NoError(t, err)
            assert.Equal(t, createdUser, storedUser)
        })
    }
}
```

### 并发测试示例
```go
func TestConcurrentCounter(t *testing.T) {
    counter := &SafeCounter{}
    const goroutines = 100
    const increments = 1000
    
    var wg sync.WaitGroup
    wg.Add(goroutines)
    
    // 启动多个 goroutine 并发增加计数器
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < increments; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    
    expected := goroutines * increments
    actual := counter.Value()
    
    assert.Equal(t, expected, actual, 
        "expected %d, got %d after %d goroutines each incrementing %d times", 
        expected, actual, goroutines, increments)
}

// 线程安全的计数器实现
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### 性能测试示例
```go
func BenchmarkStringConcatenation(b *testing.B) {
    strs := make([]string, 1000)
    for i := range strs {
        strs[i] = fmt.Sprintf("string_%d", i)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = strings.Join(strs, "")
    }
}

func BenchmarkStringBuilder(b *testing.B) {
    strs := make([]string, 1000)
    for i := range strs {
        strs[i] = fmt.Sprintf("string_%d", i)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        for _, s := range strs {
            sb.WriteString(s)
        }
        _ = sb.String()
    }
}
```

这些示例展示了不同场景下的 Go 单元测试编写方法，涵盖了从基础函数测试到复杂的集成测试场景。