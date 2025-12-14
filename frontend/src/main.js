import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor'
import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import Prism from 'prismjs'
import router from './router'
import App from './App.vue'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

// Register Element Plus Icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

VMdEditor.use(vuepressTheme, {
  Prism,
})

app.use(pinia)
app.use(router)
app.use(ElementPlus)
app.use(VMdEditor)
app.use(VueMonacoEditorPlugin, {
  paths: {
    // 使用 CDN 资源加速加载，或者你可以配置本地
    vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.43.0/min/vs'
  }
})

// Global Error Handler
app.config.errorHandler = (err, instance, info) => {
  console.error('Global Error:', err)
  // Prevent infinite loops if error happens in notification
  try {
    // Use a simpler notification to avoid component issues
    console.error(`Error: ${err.message || 'Unknown error'}`)
  } catch (e) {}
}

app.mount('#app')
