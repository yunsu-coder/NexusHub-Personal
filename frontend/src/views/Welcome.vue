<template>
  <div class="welcome-page">
    <!-- 烟花画布 -->
    <canvas ref="fireworksCanvas" class="fireworks-canvas"></canvas>

    <!-- 欢迎内容 -->
    <div class="welcome-content">
      <div class="logo-container">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style="stop-color:#60a5fa;stop-opacity:1" />
                <stop offset="100%" style="stop-color:#3b82f6;stop-opacity:1" />
              </linearGradient>
            </defs>
            <circle cx="100" cy="100" r="80" fill="url(#logoGradient)" opacity="0.2"/>
            <path d="M100,30 L130,70 L170,70 L140,100 L160,140 L100,110 L40,140 L60,100 L30,70 L70,70 Z"
                  fill="url(#logoGradient)" stroke="white" stroke-width="3"/>
          </svg>
        </div>
        <h1 class="title">NexusHub Personal</h1>
        <p class="subtitle">您的个人工作站系统</p>
      </div>

      <div class="features">
        <div class="feature-item" v-for="(feature, index) in features" :key="index">
          <div class="feature-icon">{{ feature.icon }}</div>
          <div class="feature-text">{{ feature.text }}</div>
        </div>
      </div>

      <div class="action-buttons">
        <div class="loading-animation">
          <div class="loader"></div>
          <p>正在加载...</p>
        </div>
        <el-button class="enter-btn" type="primary" @click="enter">立即进入</el-button>
        <div class="quick-links-row">
          <el-button size="large" @click="goTo('/notes')">✏️ 笔记</el-button>
          <el-button size="large" @click="goTo('/todos')">✅ TODO</el-button>
          <el-button size="large" @click="goTo('/files')">📁 文件</el-button>
          <el-button size="large" @click="goTo('/chat')">🤖 AI 聊天</el-button>
        </div>
      </div>

      <div class="welcome-footer">
        <p>v3.0.1 | 功能丰富 | 安全可靠</p>
      </div>
    </div>


  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const fireworksCanvas = ref(null)

const features = [
  { icon: '📁', text: '文件管理' },
  { icon: '📝', text: '笔记编辑' },
  { icon: '💻', text: '代码高亮' },
  { icon: '📊', text: '数据分析' },
  { icon: '🤖', text: 'AI聊天' },
  { icon: '✅', text: 'TODO管理' },
  { icon: '🧮', text: '科学计算' },
  { icon: '🎨', text: '主题定制' }
]



// 不再强制自动跳转，保留欢迎页作为首页 Hub

// 烟花动画
class Firework {
  constructor(canvas) {
    this.canvas = canvas
    this.ctx = canvas.getContext('2d')
    this.particles = []
    this.rockets = []
  }

  createRocket() {
    // 浅色调烟花颜色
    const colors = ['#60a5fa', '#34d399', '#fbbf24', '#f472b6', '#a78bfa']
    return {
      x: Math.random() * this.canvas.width,
      y: this.canvas.height,
      vx: (Math.random() - 0.5) * 2,
      vy: -(Math.random() * 8 + 12),
      color: colors[Math.floor(Math.random() * colors.length)],
      exploded: false
    }
  }

  createParticles(x, y, color) {
    const particleCount = 50
    for (let i = 0; i < particleCount; i++) {
      const angle = (Math.PI * 2 * i) / particleCount
      const velocity = Math.random() * 5 + 2
      this.particles.push({
        x,
        y,
        vx: Math.cos(angle) * velocity,
        vy: Math.sin(angle) * velocity,
        color,
        alpha: 1,
        life: 1
      })
    }
  }

  update() {
    // 更新火箭
    this.rockets.forEach((rocket, index) => {
      rocket.x += rocket.vx
      rocket.y += rocket.vy
      rocket.vy += 0.1

      if (rocket.vy >= 0 && !rocket.exploded) {
        rocket.exploded = true
        this.createParticles(rocket.x, rocket.y, rocket.color)
        this.rockets.splice(index, 1)
      }
    })

    // 更新粒子
    this.particles.forEach((particle, index) => {
      particle.x += particle.vx
      particle.y += particle.vy
      particle.vy += 0.15
      particle.life -= 0.01
      particle.alpha = particle.life

      if (particle.life <= 0) {
        this.particles.splice(index, 1)
      }
    })
  }

  draw() {
    this.ctx.fillStyle = 'rgba(240, 249, 255, 0.1)'
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height)

    // 绘制火箭
    this.rockets.forEach(rocket => {
      this.ctx.beginPath()
      this.ctx.arc(rocket.x, rocket.y, 3, 0, Math.PI * 2)
      this.ctx.fillStyle = rocket.color
      this.ctx.fill()
    })

    // 绘制粒子
    this.particles.forEach(particle => {
      this.ctx.beginPath()
      this.ctx.arc(particle.x, particle.y, 2, 0, Math.PI * 2)
      this.ctx.fillStyle = particle.color
      this.ctx.globalAlpha = particle.alpha
      this.ctx.fill()
    })
    this.ctx.globalAlpha = 1
  }

  animate() {
    this.update()
    this.draw()
  }
}

let fireworkInstance = null
let animationId = null

const initFireworks = () => {
  const canvas = fireworksCanvas.value
  if (!canvas) return

  canvas.width = window.innerWidth
  canvas.height = window.innerHeight

  fireworkInstance = new Firework(canvas)

  // 定期发射火箭
  setInterval(() => {
    if (fireworkInstance.rockets.length < 3) {
      fireworkInstance.rockets.push(fireworkInstance.createRocket())
    }
  }, 800)

  // 动画循环
  const animate = () => {
    fireworkInstance.animate()
    animationId = requestAnimationFrame(animate)
  }
  animate()
}

const showAuthDialog = (type) => {
  isLogin.value = type === 'login'
  authDialogVisible.value = true
}

const enter = () => {
  router.push('/dashboard')
}

const goTo = (path) => {
  router.push(path)
}

const handleAuth = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const endpoint = isLogin.value ? '/auth/login' : '/auth/register'
      const payload = isLogin.value
        ? { username: formData.username, password: formData.password }
        : { username: formData.username, email: formData.email, password: formData.password }

      const response = await api.post(endpoint, payload)

      // 保存token和用户信息
      if (response?.token !== undefined) {
        localStorage.setItem('token', response.token || '')
      }
      if (response?.user !== undefined) {
        localStorage.setItem('user', JSON.stringify(response.user))
      }

      ElMessage.success(isLogin.value ? '登录成功' : '注册成功')
      authDialogVisible.value = false

      // 跳转到仪表盘
      router.push('/dashboard')
    } catch (error) {
      console.error('Auth error:', error)
      const message = error.response?.data?.error || (isLogin.value ? '登录失败' : '注册失败')
      ElMessage.error(message)
    } finally {
      loading.value = false
    }
  })
}

const resetForm = () => {
  formData.username = ''
  formData.email = ''
  formData.password = ''
  formData.confirmPassword = ''
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

onMounted(() => {
  // 如果已登录,直接跳转
  const token = localStorage.getItem('token')
  if (token) {
    router.push('/dashboard')
    return
  }

  initFireworks()

  window.addEventListener('resize', () => {
    if (fireworksCanvas.value) {
      fireworksCanvas.value.width = window.innerWidth
      fireworksCanvas.value.height = window.innerHeight
    }
  })
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
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
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 50%, #bae6fd 100%);
  z-index: 9999;
}

.fireworks-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.welcome-content {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #1e40af;
  text-align: center;
  padding: 40px;
}

.logo-container {
  margin-bottom: 40px;
  animation: fadeInDown 1s ease;
}

.logo {
  margin-bottom: 30px;
  animation: float 3s ease-in-out infinite;
}

.logo-icon {
  width: 150px;
  height: 150px;
  filter: drop-shadow(0 10px 20px rgba(96, 165, 250, 0.3));
}

.title {
  font-size: 56px;
  font-weight: bold;
  margin: 0 0 15px 0;
  text-shadow: 2px 2px 8px rgba(96, 165, 250, 0.3);
  letter-spacing: 2px;
}

.subtitle {
  font-size: 22px;
  opacity: 0.95;
  margin: 0;
  text-shadow: 1px 1px 4px rgba(96, 165, 250, 0.2);
}

.features {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 25px;
  margin: 40px 0;
  animation: fadeInUp 1s ease 0.3s both;
  max-width: 900px;
}

.feature-item {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  padding: 25px 20px;
  border-radius: 20px;
  transition: all 0.3s ease;
  border: 1px solid rgba(96, 165, 250, 0.2);
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.95);
  transform: translateY(-8px);
  box-shadow: 0 15px 35px rgba(96, 165, 250, 0.3);
}

.feature-icon {
  font-size: 42px;
  margin-bottom: 12px;
  color: white;
}

.feature-icon i {
  font-size: 42px;
}

.feature-text {
  font-size: 16px;
  font-weight: 500;
  color: #1e40af;
}

.action-buttons {
  margin: 50px 0 20px;
  animation: fadeInUp 1s ease 0.6s both;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 25px;
}

.enter-btn {
  font-size: 22px;
  padding: 25px 70px;
  height: auto;
  border-radius: 50px;
  background: white;
  color: #667eea;
  border: none;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  font-weight: 600;
}

.enter-btn:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
  background: #f8f9ff;
}

.auth-links {
  display: flex;
  align-items: center;
  gap: 15px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  padding: 12px 30px;
  border-radius: 50px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.auth-link-btn {
  color: white !important;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.auth-link-btn:hover {
  transform: scale(1.1);
  color: #f0f0f0 !important;
}

.loading-animation {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.loader {
  width: 60px;
  height: 60px;
  border: 5px solid rgba(96, 165, 250, 0.3);
  border-top: 5px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-animation p {
  font-size: 20px;
  font-weight: 500;
  color: #1e40af;
  margin: 0;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.welcome-footer {
  position: absolute;
  bottom: 25px;
  opacity: 0.8;
  font-size: 14px;
  text-shadow: 1px 1px 3px rgba(96, 165, 250, 0.2);
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-15px);
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .features {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }

  .title {
    font-size: 42px;
  }
}

@media (max-width: 768px) {
  .features {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }

  .title {
    font-size: 36px;
  }

  .enter-btn {
    font-size: 18px;
    padding: 20px 50px;
  }
}
</style>
