<template>
  <div class="data-analysis-page">
    <el-row :gutter="20">
      <!-- 文件上传区域 -->
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>📊 数据分析 - 上传数据文件</span>
          </template>

          <el-upload
            ref="uploadRef"
            :action="uploadUrl"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            :show-file-list="true"
            :limit="1"
            drag
          >
            <el-icon class="el-icon--upload"><el-icon-upload-filled /></el-icon>
            <div class="el-upload__text">
              拖拽 CSV/Excel 文件到此处或 <em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 .csv, .xlsx, .xlsm, .xls, .pdf 文件，最大 50MB
              </div>
            </template>
          </el-upload>
        </el-card>
      </el-col>

      <!-- 数据预览 -->
      <el-col :span="24" v-if="dataLoaded">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>📋 数据预览 (前 {{ previewRows }} 行)</span>
              <div>
                <el-tag>总行数: {{ totalRows }}</el-tag>
                <el-tag type="success" style="margin-left: 10px">总列数: {{ totalColumns }}</el-tag>
              </div>
            </div>
          </template>

          <el-table
            :data="previewData"
            style="width: 100%"
            max-height="400"
            border
            stripe
          >
            <el-table-column
              v-for="(col, index) in columns"
              :key="index"
              :prop="col"
              :label="col"
              :min-width="120"
            />
          </el-table>
        </el-card>
      </el-col>

      <!-- 统计信息 -->
      <el-col :span="12" v-if="dataLoaded && statistics">
        <el-card>
          <template #header>
            <span>📈 数据统计</span>
          </template>

          <div class="statistics-grid">
            <div class="stat-item" v-for="(stat, column) in statistics" :key="column">
              <h4>{{ column }}</h4>
              <el-descriptions :column="1" border size="small">
                <el-descriptions-item label="类型">{{ stat.type }}</el-descriptions-item>
                <el-descriptions-item label="唯一值" v-if="stat.unique">{{ stat.unique }}</el-descriptions-item>
                <el-descriptions-item label="最小值" v-if="stat.min !== undefined">{{ stat.min }}</el-descriptions-item>
                <el-descriptions-item label="最大值" v-if="stat.max !== undefined">{{ stat.max }}</el-descriptions-item>
                <el-descriptions-item label="平均值" v-if="stat.mean !== undefined">{{ stat.mean?.toFixed(2) }}</el-descriptions-item>
                <el-descriptions-item label="缺失值" v-if="stat.missing">{{ stat.missing }}</el-descriptions-item>
              </el-descriptions>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 图表展示 -->
      <el-col :span="12" v-if="dataLoaded">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 10px">
              <span>📊 数据可视化</span>
              <div style="display: flex; gap: 10px">
                <el-select v-model="chartType" placeholder="图表类型" style="width: 130px" @change="updateChart">
                  <el-option label="柱状图" value="bar" />
                  <el-option label="折线图" value="line" />
                  <el-option label="饼图" value="pie" />
                  <el-option label="散点图" value="scatter" />
                  <el-option label="面积图" value="area" />
                </el-select>
                <el-select v-model="selectedColumn" placeholder="选择列" style="width: 150px" @change="updateChart">
                  <el-option
                    v-for="col in numericColumns"
                    :key="col"
                    :label="col"
                    :value="col"
                  />
                </el-select>
              </div>
            </div>
          </template>

          <div ref="chartRef" style="width: 100%; height: 400px"></div>
        </el-card>
      </el-col>

      <!-- AI 分析建议 -->
      <el-col :span="24" v-if="dataLoaded">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>🤖 AI 数据分析建议</span>
              <el-button type="primary" @click="getAIAnalysis" :loading="aiLoading">
                生成分析报告
              </el-button>
            </div>
          </template>

          <div v-if="aiAnalysis" class="ai-analysis">
            <pre>{{ aiAnalysis }}</pre>
          </div>
          <el-empty v-else description="点击按钮生成 AI 分析报告" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import * as XLSX from 'xlsx'
import * as pdfjsLib from 'pdfjs-dist'
import config from '../config'

// 配置 PDF.js worker
pdfjsLib.GlobalWorkerOptions.workerSrc = `//cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjsLib.version}/pdf.worker.min.js`

const uploadRef = ref()
const chartRef = ref()
const dataLoaded = ref(false)
const previewRows = ref(10)
const totalRows = ref(0)
const totalColumns = ref(0)
const columns = ref([])
const previewData = ref([])
const fullData = ref([])
const statistics = ref(null)
const selectedColumn = ref('')
const chartType = ref('bar')
const chartInstance = ref(null)
const aiAnalysis = ref('')
const aiLoading = ref(false)

// 上传请求头 - 包含JWT token
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return {
    Authorization: `Bearer ${token}`
  }
})

const numericColumns = computed(() => {
  if (!statistics.value) return []
  return Object.keys(statistics.value).filter(col =>
    statistics.value[col].type === 'number'
  )
})

const beforeUpload = (file) => {
  const allowedExtensions = ['.csv', '.xlsx', '.xlsm', '.xls', '.pdf']
  const fileExtension = '.' + file.name.toLowerCase().split('.').pop()
  const isValidType = allowedExtensions.includes(fileExtension)
  const isLt50M = file.size / 1024 / 1024 < 50

  if (!isValidType) {
    ElMessage.error(`文件格式不支持！请上传以下格式的文件：${allowedExtensions.join(', ')}`)
    return false
  }
  if (!isLt50M) {
    ElMessage.error(`文件大小超出限制！最大支持 50MB，当前文件大小：${(file.size / 1024 / 1024).toFixed(2)}MB`)
    return false
  }
  return true
}

const handleUploadSuccess = async (response, file) => {
  ElMessage.success('文件上传成功，正在解析...')

  // 下载并解析文件
  try {
    await parseFile(response.id, file.name)
  } catch (error) {
    console.error('Parse error:', error)
    ElMessage.error(`文件解析失败: ${error.message}`)
  }
}

const handleUploadError = (error) => {
  // 解析错误信息
  let errorMsg = '文件上传失败'
  try {
    const errorData = JSON.parse(error.message)
    if (errorData && errorData.error) {
      errorMsg += `: ${errorData.error}`
    }
  } catch {
    // 如果无法解析JSON，使用原始错误信息
    if (error.response && error.response.status) {
      errorMsg += ` (HTTP ${error.response.status})`
    }
  }
  ElMessage.error(errorMsg)
}

const parseFile = async (fileId, filename) => {
  const fileUrl = `${config.api.baseURL}/files/download/${fileId}`

  try {
    ElMessage.info('正在解析文件，请稍候...')
    
    // 根据文件扩展名判断类型
    const ext = filename.toLowerCase().split('.').pop()

    if (ext === 'csv') {
      await parseCSV(fileUrl)
    } else if (['xlsx', 'xlsm', 'xls'].includes(ext)) {
      await parseExcel(fileUrl)
    } else if (ext === 'pdf') {
      await parsePDF(fileUrl)
    } else {
      throw new Error(`不支持的文件格式: ${ext}`)
    }

    // 检查是否有解析到数据
    if (fullData.value.length === 0) {
      throw new Error('文件中没有找到有效的数据')
    }

    totalRows.value = fullData.value.length
    previewData.value = fullData.value.slice(0, previewRows.value)

    // 计算统计信息
    calculateStatistics()

    dataLoaded.value = true

    // 初始化图表
    await nextTick()
    if (numericColumns.value.length > 0) {
      selectedColumn.value = numericColumns.value[0]
      initChart()
      updateChart()
    }

    // 显示成功消息
    ElMessage.success(`文件解析成功！共解析 ${fullData.value.length} 行数据，${columns.value.length} 个字段`)
  } catch (error) {
    console.error('文件解析失败:', error)
    // 根据错误类型提供更具体的反馈
    let errorMsg = '文件解析失败'
    if (error.message) {
      errorMsg += `: ${error.message}`
    }
    ElMessage.error(errorMsg)
  }
}

const parseCSV = async (fileUrl) => {
  const response = await fetch(fileUrl)
  const text = await response.text()

  // 简单的CSV解析
  const lines = text.split('\n').filter(line => line.trim())
  if (lines.length === 0) {
    throw new Error('文件为空')
  }

  // 解析表头
  columns.value = lines[0].split(',').map(col => col.trim())
  totalColumns.value = columns.value.length

  // 解析数据
  fullData.value = []
  for (let i = 1; i < lines.length; i++) {
    const values = lines[i].split(',')
    const row = {}
    columns.value.forEach((col, index) => {
      row[col] = values[index]?.trim() || ''
    })
    fullData.value.push(row)
  }
}

const parseExcel = async (fileUrl) => {
  const response = await fetch(fileUrl)
  const arrayBuffer = await response.arrayBuffer()

  const workbook = XLSX.read(arrayBuffer, { type: 'array' })
  const sheetName = workbook.SheetNames[0]
  const worksheet = workbook.Sheets[sheetName]

  // 转换为 JSON
  const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 })

  if (jsonData.length === 0) {
    throw new Error('文件为空')
  }

  // 解析表头
  columns.value = jsonData[0].map(col => String(col || '').trim())
  totalColumns.value = columns.value.length

  // 解析数据
  fullData.value = []
  for (let i = 1; i < jsonData.length; i++) {
    const row = {}
    columns.value.forEach((col, index) => {
      row[col] = String(jsonData[i][index] || '').trim()
    })
    fullData.value.push(row)
  }
}

const parsePDF = async (fileUrl) => {
  const response = await fetch(fileUrl)
  const arrayBuffer = await response.arrayBuffer()

  const loadingTask = pdfjsLib.getDocument({ data: arrayBuffer })
  const pdf = await loadingTask.promise

  let allText = ''

  // 提取所有页面的文本
  for (let i = 1; i <= pdf.numPages; i++) {
    const page = await pdf.getPage(i)
    const textContent = await page.getTextContent()
    const pageText = textContent.items.map(item => item.str).join(' ')
    allText += pageText + '\n'
  }

  // 尝试从文本中提取表格数据
  // 这是一个简化的实现，假设文本中有逗号或制表符分隔的数据
  const lines = allText.split('\n').filter(line => line.trim())

  if (lines.length === 0) {
    throw new Error('PDF 文件中没有可解析的表格数据')
  }

  // 尝试检测分隔符（逗号或制表符）
  const firstLine = lines[0]
  const delimiter = firstLine.includes(',') ? ',' : '\t'

  // 解析表头
  columns.value = firstLine.split(delimiter).map(col => col.trim())
  totalColumns.value = columns.value.length

  // 解析数据
  fullData.value = []
  for (let i = 1; i < lines.length; i++) {
    const values = lines[i].split(delimiter)
    if (values.length === columns.value.length) {
      const row = {}
      columns.value.forEach((col, index) => {
        row[col] = values[index]?.trim() || ''
      })
      fullData.value.push(row)
    }
  }

  if (fullData.value.length === 0) {
    throw new Error('PDF 文件中没有找到有效的表格数据')
  }
}

const calculateStatistics = () => {
  statistics.value = {}

  columns.value.forEach(col => {
    const values = fullData.value.map(row => row[col]).filter(v => v !== '')
    const stat = {
      unique: new Set(values).size,
      missing: fullData.value.length - values.length
    }

    // 尝试转换为数字
    const numbers = values.map(v => parseFloat(v)).filter(n => !isNaN(n))
    if (numbers.length > values.length * 0.5) {
      // 大部分是数字
      stat.type = 'number'
      stat.min = Math.min(...numbers)
      stat.max = Math.max(...numbers)
      stat.mean = numbers.reduce((a, b) => a + b, 0) / numbers.length
    } else {
      stat.type = 'string'
    }

    statistics.value[col] = stat
  })
}

const initChart = () => {
  if (chartInstance.value) {
    chartInstance.value.dispose()
  }
  chartInstance.value = echarts.init(chartRef.value)
}

const updateChart = () => {
  if (!selectedColumn.value || !chartInstance.value) return

  const values = fullData.value
    .map(row => parseFloat(row[selectedColumn.value]))
    .filter(v => !isNaN(v))

  let option = {}

  switch (chartType.value) {
    case 'bar':
      option = createBarChart(values)
      break
    case 'line':
      option = createLineChart(values)
      break
    case 'pie':
      option = createPieChart(values)
      break
    case 'scatter':
      option = createScatterChart(values)
      break
    case 'area':
      option = createAreaChart(values)
      break
    default:
      option = createBarChart(values)
  }

  chartInstance.value.setOption(option, true)
}

const createBarChart = (values) => {
  const min = Math.min(...values)
  const max = Math.max(...values)
  const bins = 20
  const binSize = (max - min) / bins
  const histogram = new Array(bins).fill(0)

  values.forEach(v => {
    const binIndex = Math.min(Math.floor((v - min) / binSize), bins - 1)
    histogram[binIndex]++
  })

  return {
    title: {
      text: `${selectedColumn.value} 分布 (柱状图)`,
      textStyle: { color: '#fff' }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    xAxis: {
      type: 'category',
      data: histogram.map((_, i) => (min + i * binSize).toFixed(2)),
      axisLabel: { color: '#ccc', rotate: 45 }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#ccc' }
    },
    series: [{
      data: histogram,
      type: 'bar',
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#83bff6' },
          { offset: 0.5, color: '#188df0' },
          { offset: 1, color: '#188df0' }
        ])
      }
    }],
    backgroundColor: 'transparent',
    grid: { left: '10%', right: '5%', bottom: '15%', containLabel: true }
  }
}

const createLineChart = (values) => {
  return {
    title: {
      text: `${selectedColumn.value} 趋势 (折线图)`,
      textStyle: { color: '#fff' }
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: values.map((_, i) => i + 1),
      axisLabel: { color: '#ccc' }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#ccc' }
    },
    series: [{
      data: values,
      type: 'line',
      smooth: true,
      itemStyle: { color: '#67c23a' },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
          { offset: 1, color: 'rgba(103, 194, 58, 0.05)' }
        ])
      }
    }],
    backgroundColor: 'transparent',
    grid: { left: '10%', right: '5%', bottom: '10%', containLabel: true }
  }
}

const createPieChart = (values) => {
  // 分组数据
  const ranges = [
    { min: Math.min(...values), max: statistics.value[selectedColumn.value].mean, name: '低于平均' },
    { min: statistics.value[selectedColumn.value].mean, max: Math.max(...values), name: '高于平均' }
  ]

  const pieData = ranges.map(range => ({
    name: range.name,
    value: values.filter(v => v >= range.min && v < range.max).length
  }))

  return {
    title: {
      text: `${selectedColumn.value} 分布 (饼图)`,
      textStyle: { color: '#fff' }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      textStyle: { color: '#ccc' }
    },
    series: [{
      name: selectedColumn.value,
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: true,
        color: '#fff'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 20,
          fontWeight: 'bold'
        }
      },
      data: pieData
    }],
    backgroundColor: 'transparent'
  }
}

const createScatterChart = (values) => {
  const scatterData = values.map((v, i) => [i, v])

  return {
    title: {
      text: `${selectedColumn.value} 散点图`,
      textStyle: { color: '#fff' }
    },
    tooltip: {
      trigger: 'item',
      formatter: 'Index: {c0}<br/>Value: {c1}'
    },
    xAxis: {
      type: 'value',
      name: '索引',
      nameTextStyle: { color: '#ccc' },
      axisLabel: { color: '#ccc' }
    },
    yAxis: {
      type: 'value',
      name: '值',
      nameTextStyle: { color: '#ccc' },
      axisLabel: { color: '#ccc' }
    },
    series: [{
      type: 'scatter',
      data: scatterData,
      symbolSize: 8,
      itemStyle: {
        color: '#e6a23c',
        opacity: 0.8
      }
    }],
    backgroundColor: 'transparent',
    grid: { left: '10%', right: '5%', bottom: '10%', containLabel: true }
  }
}

const createAreaChart = (values) => {
  return {
    title: {
      text: `${selectedColumn.value} 面积图`,
      textStyle: { color: '#fff' }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' }
    },
    xAxis: {
      type: 'category',
      data: values.map((_, i) => i + 1),
      boundaryGap: false,
      axisLabel: { color: '#ccc' }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#ccc' }
    },
    series: [{
      data: values,
      type: 'line',
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(245, 108, 108, 0.5)' },
          { offset: 1, color: 'rgba(245, 108, 108, 0.1)' }
        ])
      },
      smooth: true,
      itemStyle: { color: '#f56c6c' }
    }],
    backgroundColor: 'transparent',
    grid: { left: '10%', right: '5%', bottom: '10%', containLabel: true }
  }
}

const getAIAnalysis = async () => {
  const aiSettings = localStorage.getItem('ai_settings')
  if (!aiSettings) {
    ElMessage.warning('请先在设置中配置 AI API')
    return
  }

  const settings = JSON.parse(aiSettings)
  if (!settings.apiKey) {
    ElMessage.warning('请先配置 AI API Key')
    return
  }

  aiLoading.value = true

  try {
    // 准备数据摘要
    const summary = {
      rows: totalRows.value,
      columns: totalColumns.value,
      columnNames: columns.value,
      statistics: statistics.value
    }

    const prompt = `请分析以下数据集并提供洞察：

数据摘要：
- 总行数：${summary.rows}
- 总列数：${summary.columns}
- 列名：${summary.columnNames.join(', ')}

各列统计信息：
${JSON.stringify(summary.statistics, null, 2)}

请提供：
1. 数据质量评估
2. 关键发现和趋势
3. 可能的异常值
4. 建议的分析方向
5. 数据改进建议`

    const response = await fetch(settings.apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${settings.apiKey}`
      },
      body: JSON.stringify({
        model: settings.model,
        messages: [{ role: 'user', content: prompt }],
        max_tokens: 1500,
        temperature: 0.7
      })
    })

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`)
    }

    const data = await response.json()
    aiAnalysis.value = data.choices?.[0]?.message?.content || '分析失败'
    ElMessage.success('AI 分析完成！')
  } catch (error) {
    ElMessage.error(`AI 分析失败: ${error.message}`)
  } finally {
    aiLoading.value = false
  }
}

onMounted(() => {
  // 初始化
})
</script>

<style scoped>
.data-analysis-page {
  min-height: 70vh;
}

.el-icon--upload {
  font-size: 67px;
  color: #409EFF;
  margin-bottom: 16px;
}

.statistics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 15px;
  max-height: 500px;
  overflow-y: auto;
}

.stat-item h4 {
  margin: 0 0 10px 0;
  color: var(--text-primary);
  font-weight: 600;
}

.ai-analysis {
  padding: 20px;
  background: var(--hover-bg);
  border-radius: 6px;
  max-height: 400px;
  overflow-y: auto;
}

.ai-analysis pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--text-primary);
  line-height: 1.6;
}
</style>
const uploadUrl = computed(() => `${config.api.baseURL}/files/upload`)
