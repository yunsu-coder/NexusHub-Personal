<template>
  <div class="wooden-fish-game">
    <div class="header">
      <h2>🧘 电子木鱼 - 积功德</h2>
      <div class="stats">
        <span class="merit-count">功德: {{ merit }}</span>
      </div>
    </div>
    
    <div class="game-area" @click="knock">
      <div class="wooden-fish" :class="{ 'knocking': isKnocking }">
        <svg viewBox="0 0 100 100" class="fish-icon" width="200" height="200" style="display: block;">
          <circle cx="50" cy="50" r="48" fill="#f0f0f0" stroke="none" />
          <path d="M10,50 Q10,20 50,20 Q90,20 90,50 Q90,80 50,80 Q10,80 10,50 Z" fill="#8B4513" stroke="#5D4037" stroke-width="2"/>
          <path d="M20,50 Q20,30 50,30 Q80,30 80,50 Q80,70 50,70 Q20,70 20,50 Z" fill="#A0522D"/>
          <circle cx="70" cy="50" r="5" fill="#3E2723"/>
          <path d="M30,50 L40,50" stroke="#3E2723" stroke-width="2"/>
        </svg>
      </div>
      <div class="stick" :class="{ 'striking': isKnocking }"></div>
      
      <!-- 浮动文字 -->
      <transition-group name="float-up" tag="div" class="floating-text-container">
        <div v-for="item in floatingTexts" :key="item.id" class="floating-text" :style="{ left: item.x + 'px', top: item.y + 'px' }">
          功德 +1
        </div>
      </transition-group>
    </div>

    <div class="controls">
      <el-switch v-model="autoKnock" active-text="自动敲击" @change="toggleAutoKnock" />
      <el-button type="primary" size="small" @click="reset">重置功德</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const merit = ref(0)
const isKnocking = ref(false)
const autoKnock = ref(false)
const floatingTexts = ref([])
let autoTimer = null
let idCounter = 0

// 音效 (使用简单的Web Audio API振荡器模拟)
const playSound = () => {
  const AudioContext = window.AudioContext || window.webkitAudioContext
  if (!AudioContext) return
  
  const ctx = new AudioContext()
  const osc = ctx.createOscillator()
  const gain = ctx.createGain()
  
  osc.type = 'sine'
  osc.frequency.setValueAtTime(200, ctx.currentTime)
  osc.frequency.exponentialRampToValueAtTime(100, ctx.currentTime + 0.1)
  
  gain.gain.setValueAtTime(1.0, ctx.currentTime)
  gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.1)
  
  osc.connect(gain)
  gain.connect(ctx.destination)
  
  osc.start()
  osc.stop(ctx.currentTime + 0.1)
}

const knock = (e) => {
  if (isKnocking.value) return // 防止过快点击动画未完成
  
  merit.value++
  isKnocking.value = true
  playSound()
  
  // 添加浮动文字
  const x = e ? e.offsetX : 150
  const y = e ? e.offsetY : 100
  const id = idCounter++
  floatingTexts.value.push({ id, x, y })
  
  setTimeout(() => {
    isKnocking.value = false
  }, 100)
  
  setTimeout(() => {
    floatingTexts.value = floatingTexts.value.filter(item => item.id !== id)
  }, 1000)
}

const toggleAutoKnock = (val) => {
  if (val) {
    autoTimer = setInterval(() => {
      knock()
    }, 1000)
  } else {
    clearInterval(autoTimer)
  }
}

const reset = () => {
  merit.value = 0
  localStorage.setItem('merit', 0)
}

onMounted(() => {
  const saved = localStorage.getItem('merit')
  if (saved) merit.value = parseInt(saved)
})

onUnmounted(() => {
  clearInterval(autoTimer)
  localStorage.setItem('merit', merit.value)
})
</script>

<style scoped>
.wooden-fish-game {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: var(--shadow-md);
  height: 100%;
}

.header {
  text-align: center;
  margin-bottom: 20px;
}

.merit-count {
  font-size: 24px;
  font-weight: bold;
  color: var(--primary-color);
}

.game-area {
  position: relative;
  width: 300px;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.wooden-fish {
  width: 200px;
  height: 200px;
  transition: transform 0.1s;
}

.wooden-fish.knocking {
  transform: scale(0.95);
}

.fish-icon {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 5px 10px rgba(0,0,0,0.3));
}

.stick {
  position: absolute;
  top: 50px;
  right: 50px;
  width: 100px;
  height: 20px;
  background: #5D4037;
  border-radius: 10px;
  transform-origin: bottom right;
  transform: rotate(-45deg);
  transition: transform 0.1s;
}

.stick.striking {
  transform: rotate(-65deg);
}

.floating-text-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.floating-text {
  position: absolute;
  font-size: 20px;
  font-weight: bold;
  color: #gold; /* fallback */
  color: #FFD700;
  text-shadow: 1px 1px 2px black;
}

.float-up-enter-active {
  transition: all 1s ease-out;
}

.float-up-enter-from {
  opacity: 1;
  transform: translateY(0);
}

.float-up-leave-to {
  opacity: 0;
  transform: translateY(-50px);
}

.controls {
  margin-top: 20px;
  display: flex;
  gap: 20px;
  align-items: center;
}
</style>