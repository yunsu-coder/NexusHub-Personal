<template>
  <div class="code-arena-page">
    <div class="arena-header glass-panel">
      <div class="header-left">
        <h2>代码竞技场</h2>
        <span class="subtitle">实时编译，多语言支持</span>
      </div>
      <div class="actions">
        <el-select v-model="language" placeholder="选择语言" style="width: 140px">
          <el-option label="Go (1.21)" value="go" />
          <el-option label="Python (3.11)" value="python" />
          <el-option label="JavaScript (Node 18)" value="javascript" />
          <el-option label="C (GCC)" value="c" />
          <el-option label="C++ (G++)" value="cpp" />
        </el-select>
        <el-button type="success" :icon="VideoPlay" :loading="running" @click="runCode">运行代码</el-button>
      </div>
    </div>

    <div class="arena-content">
      <!-- Editor Pane -->
      <div class="pane editor-pane glass-panel">
        <div class="pane-header">
           <el-icon><Edit /></el-icon> 编辑器
        </div>
        <div class="monaco-wrapper">
          <vue-monaco-editor
            v-model:value="code"
            :language="monacoLang"
            theme="vs-dark"
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
import { VideoPlay, Edit, Monitor, Bottom } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const language = ref('go')
const code = ref('')
const input = ref('')
const output = ref('')
const running = ref(false)
const isError = ref(false)

const monacoLang = computed(() => {
  const map = {
    go: 'go',
    python: 'python',
    javascript: 'javascript',
    c: 'c',
    cpp: 'cpp'
  }
  return map[language.value] || 'plaintext'
})

const editorOptions = {
  automaticLayout: true,
  formatOnType: true,
  formatOnPaste: true,
  fontSize: 14,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  padding: { top: 10, bottom: 10 }
}

const templates = {
  go: `package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, NexusHub!")\n}`,
  python: `print("Hello, NexusHub!")`,
  javascript: `console.log("Hello, NexusHub!");`,
  c: `#include <stdio.h>\n\nint main() {\n    printf("Hello, NexusHub!\\n");\n    return 0;\n}`,
  cpp: `#include <iostream>\n\nint main() {\n    std::cout << "Hello, NexusHub!" << std::endl;\n    return 0;\n}`
}

watch(language, (newLang) => {
  const currentCode = code.value.trim()
  // Check if current code matches any template (simple check) to avoid overwriting user code
  const isTemplate = Object.values(templates).some(t => t.trim() === currentCode) || !currentCode
  
  if (isTemplate) {
    code.value = templates[newLang] || ''
  }
}, { immediate: true })

const runCode = async () => {
  if (!code.value.trim()) return
  
  running.value = true
  output.value = 'Compiling and Running...\n'
  isError.value = false
  
  try {
    const res = await api.post('/code/run', {
      language: language.value,
      code: code.value,
      input: input.value
    })
    
    // Simulate slight delay for "realism" if too fast
    if (res.output) {
        output.value = res.output
    } else {
        output.value = "Program exited with no output."
    }
    
    if (res.error) {
      output.value += '\n\n[Runtime Error]:\n' + res.error
      isError.value = true
    }
  } catch (e) {
    output.value = '[System Error]: ' + (e.response?.data?.error || e.message)
    isError.value = true
  } finally {
    running.value = false
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
  background: rgba(30, 30, 30, 0.8); /* Darker for code editor feel */
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
}

.arena-header {
  padding: 15px 25px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #fff;
}
.header-left h2 { margin: 0; font-size: 22px; }
.subtitle { opacity: 0.6; font-size: 13px; }

.actions { display: flex; gap: 15px; }

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
  background: rgba(30, 30, 30, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.pane-header {
  padding: 10px 15px;
  background: rgba(255, 255, 255, 0.05);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  color: #ddd;
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: space-between;
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
  background: transparent;
  border: none;
  color: #ddd;
  padding: 15px;
  font-family: 'Consolas', monospace;
  resize: none;
  outline: none;
}

.terminal-output {
  flex: 1;
  padding: 15px;
  color: #4ade80; /* Matrix Green */
  font-family: 'Consolas', monospace;
  overflow: auto;
  white-space: pre-wrap;
  line-height: 1.5;
}

.terminal-output.error {
  color: #f87171;
}

.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #4ade80;
  animation: blink 1s step-end infinite;
  vertical-align: middle;
}

@keyframes blink {
  50% { opacity: 0; }
}
</style>