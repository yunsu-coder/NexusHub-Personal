<template>
  <div class="music-player">
    <div class="player-controls">
      <el-button
        :icon="isPlaying ? VideoPause : VideoPlay"
        circle
        size="small"
        @click="togglePlay"
        :disabled="!hasMusic"
      />
      <div class="player-info">
        <div class="track-name">{{ trackName }}</div>
        <el-slider
          v-model="volume"
          :min="0"
          :max="100"
          @input="setVolume"
          style="width: 80px"
          :disabled="!hasMusic"
        />
      </div>
    </div>
    <audio ref="audioRef" @ended="onEnded" @error="onError" @canplay="onCanPlay" />
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { VideoPlay, VideoPause } from '@element-plus/icons-vue'
import { useThemeStore } from '../store/theme'
import { ElMessage } from 'element-plus'

const themeStore = useThemeStore()
const audioRef = ref(null)
const isPlaying = ref(false)
const volume = ref(50)
const hasMusic = ref(false)
const trackName = ref('未设置')

watch(() => themeStore.theme.background_music, (newMusic) => {
  if (newMusic && audioRef.value) {
    audioRef.value.src = newMusic
    hasMusic.value = true
    trackName.value = '背景音乐'
    audioRef.value.load()
  } else {
    hasMusic.value = false
    trackName.value = '未设置'
    isPlaying.value = false
  }
})

watch(() => themeStore.theme.music_volume, (newVolume) => {
  volume.value = Math.round(newVolume * 100)
  if (audioRef.value) {
    audioRef.value.volume = newVolume
  }
})

const togglePlay = () => {
  if (!hasMusic.value || !audioRef.value) return

  if (isPlaying.value) {
    audioRef.value.pause()
    isPlaying.value = false
  } else {
    const playPromise = audioRef.value.play()
    if (playPromise !== undefined) {
      playPromise
        .then(() => {
          isPlaying.value = true
        })
        .catch((error) => {
          console.error('播放失败:', error)
          ElMessage.error('播放失败: ' + error.message)
          isPlaying.value = false
        })
    }
  }
}

const setVolume = (val) => {
  const newVolume = val / 100
  if (audioRef.value) {
    audioRef.value.volume = newVolume
  }
  themeStore.theme.music_volume = newVolume
}

const onEnded = () => {
  isPlaying.value = false
  // 循环播放
  if (hasMusic.value) {
    audioRef.value.play()
    isPlaying.value = true
  }
}

const onError = (event) => {
  console.error('音乐加载失败:', event)
  ElMessage.error('音乐加载失败，请检查文件格式或URL')
  isPlaying.value = false
  hasMusic.value = false
}

const onCanPlay = () => {
  console.log('音乐已加载，可以播放')
  hasMusic.value = true
}

onMounted(() => {
  if (themeStore.theme.background_music) {
    audioRef.value.src = themeStore.theme.background_music
    hasMusic.value = true
    trackName.value = '背景音乐'
  }
  if (audioRef.value) {
    audioRef.value.volume = themeStore.theme.music_volume
    volume.value = Math.round(themeStore.theme.music_volume * 100)
  }
})
</script>

<style scoped>
.music-player {
  padding: 10px;
}

.player-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.player-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.track-name {
  font-size: 12px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
