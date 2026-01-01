<template>
  <div class="code-arena-page">
    <div class="arena-header glass-panel">
      <div class="header-left">
        <h2>代码竞技场</h2>
        <span class="subtitle">实时编译，多语言支持</span>
      </div>
      <div class="actions">
        <el-select v-model="language" placeholder="选择语言" style="width: 180px">
          <el-option label="Go (1.21)" value="go" />
          <el-option label="Python (3.11)" value="python" />
          <el-option label="JavaScript (Node 18)" value="javascript" />
          <el-option label="TypeScript (5.0)" value="typescript" />
          <el-option label="C (GCC 12)" value="c" />
          <el-option label="C++ (G++ 12)" value="cpp" />
          <el-option label="Java (17)" value="java" />
          <el-option label="Rust (1.70)" value="rust" />
          <el-option label="Ruby (3.2)" value="ruby" />
          <el-option label="PHP (8.2)" value="php" />
          <el-option label="Swift (5.8)" value="swift" />
          <el-option label="Kotlin (1.8)" value="kotlin" />
        </el-select>
        <el-select v-model="selectedTemplate" placeholder="代码模板" style="width: 120px">
          <el-option label="Hello World" value="hello" />
          <el-option label="数据结构" value="data" />
          <el-option label="算法" value="algorithm" />
        </el-select>
        <el-tooltip content="新建文件" placement="bottom">
          <el-button type="primary" :icon="DocumentAdd" @click="newFile">新建</el-button>
        </el-tooltip>
        <el-tooltip content="保存代码到本地" placement="bottom">
          <el-button type="warning" :icon="Download" @click="saveCode">保存</el-button>
        </el-tooltip>
        <el-tooltip content="运行代码 (Ctrl+Enter)" placement="bottom">
          <el-button type="success" :icon="VideoPlay" :loading="running" @click="runCode">运行代码</el-button>
        </el-tooltip>
      </div>
    </div>

    <div class="arena-content">
      <!-- Editor Pane -->
      <div class="pane editor-pane glass-panel">
        <div class="pane-header">
           <el-icon><Edit /></el-icon> 编辑器
        </div>
        <div class="monaco-wrapper">
          <VueMonacoEditor
            v-model="code"
            :language="monacoLang"
            :theme="monacoTheme"
            :options="editorOptions"
            class="monaco-editor-container"
          />
        </div>
      </div>
      
      <!-- IO Pane -->
      <div class="pane io-pane glass-panel">
        <div class="io-section input-section">
          <div class="pane-header">
            <span><el-icon><Bottom /></el-icon> 标准输入 (Stdin)</span>
          </div>
          <textarea v-model="input" class="terminal-input" placeholder="输入测试数据..."></textarea>
        </div>
        
        <div class="io-section output-section">
          <div class="pane-header">
            <span><el-icon><Monitor /></el-icon> 运行结果</span>
            <el-button size="small" text bg @click="output = ''">清空</el-button>
          </div>
          <div class="terminal-output" :class="{ error: isError }">
            <pre>{{ output }}</pre>
            <span class="cursor" v-if="!running">_</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { VideoPlay, Edit, Monitor, Bottom, DocumentAdd, Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import api from '../api'
import { useThemeStore } from '../store/theme'

const language = ref('go')
const code = ref('')
const input = ref('')
const output = ref('')
const running = ref(false)
const isError = ref(false)
const selectedTemplate = ref('hello')

// 主题支持
const themeStore = useThemeStore()

const monacoLang = computed(() => {
  const map = {
    go: 'go',
    python: 'python',
    javascript: 'javascript',
    typescript: 'typescript',
    c: 'c',
    cpp: 'cpp',
    java: 'java',
    rust: 'rust',
    ruby: 'ruby',
    php: 'php',
    swift: 'swift',
    kotlin: 'kotlin'
  }
  return map[language.value] || 'plaintext'
})

// 根据全局主题自动切换编辑器主题
const monacoTheme = computed(() => {
  return themeStore.darkMode ? 'vs-dark' : 'vs-light'
})

// 增强的编辑器选项
const editorOptions = {
  automaticLayout: true,
  formatOnType: true,
  formatOnPaste: true,
  fontSize: 14,
  lineNumbers: 'on',
  scrollBeyondLastLine: false,
  minimap: { enabled: true, scale: 0.8 },
  wordWrap: 'on',
  cursorBlinking: 'smooth',
  cursorStyle: 'line',
  scrollbar: {
    verticalScrollbarSize: 10,
    horizontalScrollbarSize: 10
  },
  padding: { top: 10, bottom: 10 },
  tabSize: 2,
  detectIndentation: true,
  autoIndent: 'full',
  quickSuggestions: {
    other: true,
    comments: true,
    strings: true
  },
  parameterHints: {
    enabled: true,
    cycle: true
  },
  inlayHints: {
    enabled: true
  },
  codeLens: {
    enabled: true
  },
  folding: true,
  lineDecorationsWidth: 10,
  lineNumbersMinChars: 3
}

const templates = {
  hello: {
    go: `package main

import "fmt"

func main() {
	fmt.Println("Hello, NexusHub!")
}`,
    python: `print("Hello, NexusHub!")`,
    javascript: `console.log("Hello, NexusHub!");`,
    typescript: `console.log("Hello, NexusHub!");`,
    c: `#include <stdio.h>

int main() {
    printf("Hello, NexusHub!\n");
    return 0;
}`,
    cpp: `#include <iostream>

int main() {
    std::cout << "Hello, NexusHub!" << std::endl;
    return 0;
}`,
    java: `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, NexusHub!");
    }
}`,
    rust: `fn main() {
    println!("Hello, NexusHub!");
}`,
    ruby: `puts "Hello, NexusHub!"`,
    php: `<?php
echo "Hello, NexusHub!";
?>`,
    swift: `print("Hello, NexusHub!")`,
    kotlin: `fun main() {
    println("Hello, NexusHub!")
}`
  },
  data: {
    go: `package main

import "fmt"

type Node struct {
	Value int
	Next  *Node
}

func main() {
	// 创建链表
	head := &Node{Value: 1}
	head.Next = &Node{Value: 2}
	head.Next.Next = &Node{Value: 3}
	
	// 遍历链表
	current := head
	for current != nil {
		fmt.Println(current.Value)
		current = current.Next
	}
}`,
    python: `class Node:
    def __init__(self, value):
        self.value = value
        self.next = None

# 创建链表
head = Node(1)
head.next = Node(2)
head.next.next = Node(3)

# 遍历链表
current = head
while current:
    print(current.value)
    current = current.next`,
    javascript: `class Node {
    constructor(value) {
        this.value = value;
        this.next = null;
    }
}

// 创建链表
const head = new Node(1);
head.next = new Node(2);
head.next.next = new Node(3);

// 遍历链表
let current = head;
while (current) {
    console.log(current.value);
    current = current.next;
}`
  },
  algorithm: {
    go: `package main

import "fmt"

// 二分查找
func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	
	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

func main() {
	arr := []int{1, 3, 5, 7, 9}
	target := 5
	index := binarySearch(arr, target)
	fmt.Printf("目标值 %d 在数组中的索引: %d\n", target, index)
}`,
    python: `# 二分查找
def binary_search(arr, target):
    left, right = 0, len(arr) - 1
    
    while left <= right:
        mid = left + (right - left) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    return -1

# 测试
arr = [1, 3, 5, 7, 9]
target = 5
index = binary_search(arr, target)
print(f"目标值 {target} 在数组中的索引: {index}")`,
    javascript: `// 二分查找
function binarySearch(arr, target) {
    let left = 0;
    let right = arr.length - 1;
    
    while (left <= right) {
        const mid = left + Math.floor((right - left) / 2);
        if (arr[mid] === target) {
            return mid;
        } else if (arr[mid] < target) {
            left = mid + 1;
        } else {
            right = mid - 1;
        }
    }
    return -1;
}`
  }
}

// 更新代码模板
const updateCodeFromTemplate = () => {
  const currentCode = code.value.trim();
  // 检查当前代码是否为空或匹配任何模板，避免覆盖用户代码
  const isEmpty = !currentCode;
  const isMatchingAnyTemplate = Object.values(templates).some(category => 
    Object.values(category).some(template => template.trim() === currentCode)
  );
  
  if (isEmpty || isMatchingAnyTemplate) {
    const category = templates[selectedTemplate.value] || templates.hello;
    code.value = category[language.value] || category.go || '';
  }
}

// 监听语言变化
watch(language, updateCodeFromTemplate, { immediate: true });
// 监听模板变化
watch(selectedTemplate, updateCodeFromTemplate);

// 新建文件
const newFile = async () => {
  if (code.value.trim()) {
    const confirm = await ElMessageBox.confirm(
      '当前代码未保存，是否继续新建？',
      '确认操作',
      {
        confirmButtonText: '继续',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    
    if (confirm !== 'confirm') return;
  }
  
  code.value = '';
  input.value = '';
  output.value = '';
  selectedTemplate.value = 'hello';
  ElMessage.success('已新建文件');
}

// 保存代码
const saveCode = () => {
  if (!code.value.trim()) {
    ElMessage.warning('代码为空，无法保存')
    return
  }
  
  const fileName = `code_${language.value}_${Date.now()}.txt`;
  const blob = new Blob([code.value], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  
  try {
    const a = document.createElement('a');
    a.href = url;
    a.download = fileName;
    
    // 隐藏元素但添加到DOM
    a.style.display = 'none';
    document.body.appendChild(a);
    
    // 执行点击
    a.click();
    
    // 稍后移除元素并释放URL，确保下载完成
    setTimeout(() => {
      if (document.body.contains(a)) {
        document.body.removeChild(a);
      }
      URL.revokeObjectURL(url);
    }, 100);
    
    ElMessage.success('代码已保存');
  } catch (error) {
    // 发生错误时确保释放资源
    URL.revokeObjectURL(url);
    console.error('保存代码失败:', error);
    ElMessage.error('保存代码失败，请重试');
  }
}

const runCode = async () => {
  if (!code.value.trim()) {
    ElMessage.warning('请输入代码');
    return;
  }
  
  running.value = true;
  output.value = '正在编译和运行...\n';
  isError.value = false;
  
  try {
    const res = await api.post('/code/run', {
      language: language.value,
      code: code.value,
      input: input.value
    });
    
    // 处理运行结果
    if (res && typeof res === 'object') {
      if (res.output) {
        output.value = res.output;
      } else {
        output.value = "程序运行结束，未产生输出。";
      }
      
      if (res.error) {
        output.value += '\n\n[运行时错误]:\n' + res.error;
        isError.value = true;
      }
    } else {
      output.value = '无法获取运行结果';
      isError.value = true;
    }
  } catch (e) {
    // 增强错误处理，提供更详细的错误信息
    let errorMsg = '系统错误';
    
    if (e.response) {
      // 服务器响应错误
      const status = e.response.status;
      const errorData = e.response.data;
      
      if (status === 400) {
        errorMsg = '请求参数错误';
      } else if (status === 500) {
        errorMsg = '服务器内部错误';
      } else if (status === 404) {
        errorMsg = '代码执行服务未找到';
      } else {
        errorMsg = `服务器错误 (HTTP ${status})`;
      }
      
      if (errorData && errorData.error) {
        errorMsg += ': ' + errorData.error;
      }
    } else if (e.request) {
      // 网络请求错误
      errorMsg = '网络请求失败，请检查网络连接';
    } else {
      // 其他错误
      errorMsg = '运行代码时发生错误: ' + e.message;
    }
    
    output.value = `[系统错误]: ${errorMsg}`;
    isError.value = true;
  } finally {
    running.value = false;
  }
}
</script>

<style scoped>
.code-arena-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
  gap: 20px;
}

.glass-panel {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.arena-header {
  padding: 20px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.arena-header:hover {
  box-shadow: var(--shadow-md);
}

.header-left h2 { 
  margin: 0; 
  font-size: 24px; 
  font-weight: 700;
  color: var(--text-primary);
}

.subtitle { 
  opacity: 0.8;
  font-size: 13px; 
  color: var(--text-secondary);
  display: block;
  margin-top: 4px;
}

.actions { 
  display: flex; 
  gap: 15px;
  align-items: center;
}

.arena-content {
  flex: 1;
  display: flex;
  gap: 20px;
  overflow: hidden;
}

.pane {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: all 0.3s ease;
}

.pane:hover {
  box-shadow: var(--shadow-md);
}

.editor-pane {
  flex: 3; /* 60% width */
}

.io-pane {
  flex: 2; /* 40% width */
  display: flex;
  flex-direction: column;
  background: transparent;
  border: none;
  gap: 20px;
}

.io-section {
  flex: 1;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: all 0.3s ease;
}

.io-section:hover {
  box-shadow: var(--shadow-md);
}

.pane-header {
  padding: 12px 16px;
  background: var(--bg-light);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: space-between;
  transition: all 0.3s ease;
}

.monaco-wrapper {
  flex: 1;
  position: relative;
}

.monaco-editor-container {
  position: absolute;
  inset: 0;
}

.terminal-input {
  flex: 1;
  background: var(--bg-color);
  border: none;
  color: var(--text-primary);
  padding: 16px;
  font-family: 'Consolas', monospace;
  resize: none;
  outline: none;
  font-size: 14px;
  line-height: 1.5;
  transition: all 0.3s ease;
}

.terminal-input:focus {
  background: var(--card-bg);
}

.terminal-output {
  flex: 1;
  padding: 16px;
  color: var(--text-primary);
  font-family: 'Consolas', monospace;
  overflow: auto;
  white-space: pre-wrap;
  line-height: 1.6;
  font-size: 14px;
  background: var(--bg-color);
  transition: all 0.3s ease;
}

.terminal-output.error {
  color: var(--danger-color);
}

.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: var(--primary-color);
  animation: blink 1s step-end infinite;
  vertical-align: middle;
}

@keyframes blink {
  50% { opacity: 0; }
}

/* 基础页面样式 */
.code-arena-page {
  background: var(--bg-color);
  color: var(--text-primary);
}
</style>