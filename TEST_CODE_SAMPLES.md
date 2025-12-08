# 代码编辑器测试示例

## Go 代码示例

```go
package main

import (
    "fmt"
    "time"
)

// User 定义用户结构体
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
    // 模拟数据库查询
    user := &User{
        ID:        id,
        Username:  "admin",
        Email:     "admin@example.com",
        CreatedAt: time.Now(),
    }

    return user, nil
}

func main() {
    // 获取用户
    user, err := GetUserByID(1)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // 打印用户信息
    fmt.Printf("User: %+v\n", user)

    // 遍历切片
    numbers := []int{1, 2, 3, 4, 5}
    for i, num := range numbers {
        fmt.Printf("Index %d: %d\n", i, num)
    }
}
```

## Markdown 代码示例

```markdown
# 标题 1

## 标题 2

### 标题 3

这是一段**加粗文本**和*斜体文本*。

- 列表项 1
- 列表项 2
  - 嵌套列表项 2.1
  - 嵌套列表项 2.2
- 列表项 3

1. 有序列表 1
2. 有序列表 2
3. 有序列表 3

---

> 这是一段引用文本
> 可以有多行

`这是行内代码`

代码块：

​```python
def hello_world():
    print("Hello, World!")
​```

[链接文本](https://example.com)

![图片描述](https://via.placeholder.com/150)

| 表头1 | 表头2 | 表头3 |
|-------|-------|-------|
| 数据1 | 数据2 | 数据3 |
| 数据4 | 数据5 | 数据6 |
```

## 测试步骤

### 1. 测试 Go 代码高亮

1. 进入代码编辑器页面
2. 点击"新建"
3. 标题输入："Go 示例代码"
4. 语言选择："Go"
5. 粘贴上面的 Go 代码
6. 切换到"预览"模式
7. 验证语法高亮效果：
   - ✅ 关键字（package, import, func, type, struct）应该高亮
   - ✅ 字符串应该是绿色
   - ✅ 注释应该是灰色
   - ✅ 函数名应该是黄色

### 2. 测试 Markdown 高亮

1. 新建代码片段
2. 标题输入："Markdown 示例"
3. 语言选择："Markdown"
4. 粘贴上面的 Markdown 代码
5. 切换到"预览"模式
6. 验证语法高亮效果：
   - ✅ 标题（#）应该高亮
   - ✅ 列表标记应该高亮
   - ✅ 代码块应该有背景色
   - ✅ 链接应该高亮

### 3. 测试 Tab 键功能

1. 在编辑模式下输入代码
2. 按 Tab 键
3. 验证：
   - ✅ 插入4个空格
   - ✅ 光标不跳出输入框
   - ✅ 可以正常输入代码

### 4. 测试代码换行

1. 输入长代码行
2. 验证：
   - ✅ 编辑模式支持水平滚动
   - ✅ 预览模式代码不会溢出
   - ✅ 滚动条显示正常

## 预期效果

### 编辑模式
- 等宽字体显示
- 支持 Tab 缩进（4空格）
- 水平滚动条
- 深色背景
- 蓝色聚焦边框

### 预览模式
- Atom One Dark 主题
- 完整的语法高亮
- 深色背景 (#282c34)
- 代码块样式美观
- 支持滚动查看

## 额外功能测试

### 保存和加载
1. ✅ 保存代码片段
2. ✅ 刷新页面
3. ✅ 重新加载代码片段
4. ✅ 切换预览查看高亮

### 多语言支持
测试其他语言：
- ✅ JavaScript
- ✅ Python
- ✅ C++
- ✅ Java
- ✅ HTML
- ✅ CSS

所有语言都应该有正确的语法高亮效果！
