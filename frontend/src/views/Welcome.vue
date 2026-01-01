<template>
  <div class="welcome-page">
    <!-- 欢迎内容 -->
    <div class="welcome-content">
      <div class="apple-style">
        <h1 class="hello-title">你好</h1>
        <p class="subtitle">欢迎使用 NexusHub Personal</p>
      </div>

      <div class="action-buttons">
        <el-button class="enter-btn" type="primary" @click="enter">立即进入</el-button>
      </div>

      <div class="quick-links-row">
        <el-button size="large" @click="goTo('/notes')">✏️ 笔记</el-button>
        <el-button size="large" @click="goTo('/todos')">✅ TODO</el-button>
        <el-button size="large" @click="goTo('/files')">📁 文件</el-button>
      </div>

      <div class="welcome-footer">
        <p>v3.0.1</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const enter = () => {
  router.push('/dashboard')
}

const goTo = (path) => {
  router.push(path)
}

onMounted(() => {
  // 如果已登录,直接跳转
  const token = localStorage.getItem('token')
  if (token) {
    router.push('/dashboard')
    return
  }
})
</script>

<style scoped>
.welcome-page {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #fafafa;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #1d1d1f;
  text-align: center;
  padding: 40px;
}

.apple-style {
  margin-bottom: 60px;
  animation: fadeIn 1s ease-out;
}

.hello-title {
  font-size: 96px;
  font-weight: 300;
  margin: 0;
  color: #1d1d1f;
  letter-spacing: -1px;
  animation: helloAnimation 2s ease-out;
}

.subtitle {
  font-size: 24px;
  margin: 20px 0 0 0;
  color: #86868b;
  font-weight: 400;
  animation: subtitleFadeIn 1.5s ease-out 0.5s both;
}

.action-buttons {
  margin: 40px 0;
  animation: fadeIn 1.5s ease-out 1s both;
}

.enter-btn {
  font-size: 18px;
  padding: 15px 40px;
  height: auto;
  border-radius: 10px;
  background: #0071e3;
  color: white;
  border: none;
  transition: all 0.3s ease;
  font-weight: 600;
  min-width: 200px;
}

.enter-btn:hover {
  background: #0077ed;
  transform: scale(1.05);
}

.quick-links-row {
  display: flex;
  gap: 15px;
  margin: 30px 0;
  animation: fadeIn 2s ease-out 1.5s both;
  flex-wrap: wrap;
  justify-content: center;
}

.welcome-footer {
  position: absolute;
  bottom: 25px;
  opacity: 0.6;
  font-size: 14px;
  color: #86868b;
  animation: fadeIn 2.5s ease-out 2s both;
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes helloAnimation {
  0% {
    opacity: 0;
    transform: translateY(-50px);
    font-size: 60px;
  }
  50% {
    opacity: 0.8;
    transform: translateY(10px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
    font-size: 96px;
  }
}

@keyframes subtitleFadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .hello-title {
    font-size: 64px;
  }

  @keyframes helloAnimation {
    0% {
      opacity: 0;
      transform: translateY(-50px);
      font-size: 40px;
    }
    50% {
      opacity: 0.8;
      transform: translateY(10px);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
      font-size: 64px;
    }
  }

  .subtitle {
    font-size: 20px;
  }

  .enter-btn {
    font-size: 16px;
    padding: 12px 30px;
    min-width: 160px;
  }

  .quick-links-row {
    gap: 10px;
  }
}
</style>
