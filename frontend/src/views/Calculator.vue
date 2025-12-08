<template>
  <div class="calculator-page">
    <el-row :gutter="20">
      <!-- 基础计算器 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>🔢 高级科学计算器</span>
          </template>

          <div class="calculator-display">
            <el-input
              v-model="expression"
              placeholder="输入表达式..."
              size="large"
              @keyup.enter="calculate"
            >
              <template #append>
                <el-button @click="calculate" type="primary">计算</el-button>
              </template>
            </el-input>
          </div>

          <div class="calculator-result" v-if="result !== null">
            <el-alert :title="`结果: ${result}`" type="success" :closable="false" />
          </div>

          <div class="calculator-error" v-if="error">
            <el-alert :title="error" type="error" :closable="false" />
          </div>

          <div class="calculator-buttons">
            <div class="button-group">
              <h4>三角函数</h4>
              <div class="button-row">
                <el-button @click="insertText('sin(')" class="calc-btn">sin</el-button>
                <el-button @click="insertText('cos(')" class="calc-btn">cos</el-button>
                <el-button @click="insertText('tan(')" class="calc-btn">tan</el-button>
                <el-button @click="insertText('sqrt(')" class="calc-btn">√</el-button>
              </div>
            </div>

            <div class="button-group">
              <h4>对数与指数</h4>
              <div class="button-row">
                <el-button @click="insertText('log(')" class="calc-btn">log</el-button>
                <el-button @click="insertText('ln(')" class="calc-btn">ln</el-button>
                <el-button @click="insertText('exp(')" class="calc-btn">exp</el-button>
                <el-button @click="insertText('^')" class="calc-btn">x^y</el-button>
              </div>
            </div>

            <div class="button-group">
              <h4>特殊函数</h4>
              <div class="button-row">
                <el-button @click="insertText('abs(')" class="calc-btn">|x|</el-button>
                <el-button @click="insertText('factorial(')" class="calc-btn">n!</el-button>
                <el-button @click="insertText('pi')" class="calc-btn">π</el-button>
                <el-button @click="insertText('e')" class="calc-btn">e</el-button>
              </div>
            </div>

            <div class="button-group">
              <h4>基本操作</h4>
              <div class="button-row">
                <el-button @click="insertText('(')" class="calc-btn">(</el-button>
                <el-button @click="insertText(')')" class="calc-btn">)</el-button>
                <el-button @click="expression = ''" type="warning" class="calc-btn">清空</el-button>
                <el-button @click="expression = expression.slice(0, -1)" type="danger" class="calc-btn">删除</el-button>
              </div>
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>📐 矩阵计算 (线性代数)</span>
          </template>

          <div class="matrix-input">
            <el-row :gutter="10">
              <el-col :span="12">
                <label>矩阵 A:</label>
                <el-input
                  v-model="matrixA"
                  type="textarea"
                  :rows="4"
                  placeholder="[[1, 2], [3, 4]]"
                />
              </el-col>
              <el-col :span="12">
                <label>矩阵 B:</label>
                <el-input
                  v-model="matrixB"
                  type="textarea"
                  :rows="4"
                  placeholder="[[5, 6], [7, 8]]"
                />
              </el-col>
            </el-row>
          </div>

          <div class="matrix-operations" style="margin-top: 15px">
            <el-button-group>
              <el-button @click="matrixOperation('add')">A + B</el-button>
              <el-button @click="matrixOperation('subtract')">A - B</el-button>
              <el-button @click="matrixOperation('multiply')">A × B</el-button>
              <el-button @click="matrixOperation('transpose')">转置 A</el-button>
              <el-button @click="matrixOperation('det')">det(A)</el-button>
              <el-button @click="matrixOperation('inv')">A⁻¹</el-button>
            </el-button-group>
          </div>

          <div class="matrix-result" v-if="matrixResult">
            <el-divider />
            <pre>{{ matrixResult }}</pre>
          </div>
        </el-card>
      </el-col>

      <!-- 概率统计 & 计算机专业 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>📊 概率统计计算</span>
          </template>

          <div class="stats-input">
            <label>数据集 (用逗号分隔):</label>
            <el-input
              v-model="dataSet"
              placeholder="1, 2, 3, 4, 5, 6, 7, 8, 9, 10"
              @change="calculateStats"
            />
          </div>

          <div class="stats-result" v-if="statsResult">
            <el-descriptions :column="2" border style="margin-top: 15px">
              <el-descriptions-item label="平均值">{{ statsResult.mean }}</el-descriptions-item>
              <el-descriptions-item label="中位数">{{ statsResult.median }}</el-descriptions-item>
              <el-descriptions-item label="标准差">{{ statsResult.std }}</el-descriptions-item>
              <el-descriptions-item label="方差">{{ statsResult.variance }}</el-descriptions-item>
              <el-descriptions-item label="最小值">{{ statsResult.min }}</el-descriptions-item>
              <el-descriptions-item label="最大值">{{ statsResult.max }}</el-descriptions-item>
              <el-descriptions-item label="总和">{{ statsResult.sum }}</el-descriptions-item>
              <el-descriptions-item label="样本数">{{ statsResult.count }}</el-descriptions-item>
            </el-descriptions>
          </div>

          <el-divider />

          <div class="probability-calc">
            <h4>概率计算</h4>
            <el-row :gutter="10">
              <el-col :span="12">
                <label>组合 C(n,r):</label>
                <el-input-group>
                  <el-input v-model="combN" placeholder="n" style="width: 50%" />
                  <el-input v-model="combR" placeholder="r" style="width: 50%" />
                </el-input-group>
                <el-button @click="calcCombination" style="margin-top: 5px; width: 100%">
                  计算 C({{ combN }},{{ combR }})
                </el-button>
              </el-col>
              <el-col :span="12">
                <label>排列 P(n,r):</label>
                <el-input-group>
                  <el-input v-model="permN" placeholder="n" style="width: 50%" />
                  <el-input v-model="permR" placeholder="r" style="width: 50%" />
                </el-input-group>
                <el-button @click="calcPermutation" style="margin-top: 5px; width: 100%">
                  计算 P({{ permN }},{{ permR }})
                </el-button>
              </el-col>
            </el-row>

            <div v-if="probResult" style="margin-top: 15px">
              <el-alert :title="`结果: ${probResult}`" type="info" :closable="false" />
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>💻 计算机专业计算</span>
          </template>

          <el-tabs>
            <el-tab-pane label="进制转换">
              <el-input v-model="numberInput" placeholder="输入数字">
                <template #prepend>
                  <el-select v-model="inputBase" style="width: 100px">
                    <el-option label="二进制" :value="2" />
                    <el-option label="八进制" :value="8" />
                    <el-option label="十进制" :value="10" />
                    <el-option label="十六进制" :value="16" />
                  </el-select>
                </template>
              </el-input>

              <el-button @click="convertBase" style="margin-top: 10px; width: 100%" type="primary">
                转换
              </el-button>

              <div v-if="baseResult" style="margin-top: 15px">
                <el-descriptions :column="1" border>
                  <el-descriptions-item label="二进制">{{ baseResult.bin }}</el-descriptions-item>
                  <el-descriptions-item label="八进制">{{ baseResult.oct }}</el-descriptions-item>
                  <el-descriptions-item label="十进制">{{ baseResult.dec }}</el-descriptions-item>
                  <el-descriptions-item label="十六进制">{{ baseResult.hex }}</el-descriptions-item>
                </el-descriptions>
              </div>
            </el-tab-pane>

            <el-tab-pane label="位运算">
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input v-model="bitA" placeholder="数字 A (十进制)" />
                </el-col>
                <el-col :span="12">
                  <el-input v-model="bitB" placeholder="数字 B (十进制)" />
                </el-col>
              </el-row>

              <el-button-group style="margin-top: 10px; width: 100%">
                <el-button @click="bitOperation('and')">A & B</el-button>
                <el-button @click="bitOperation('or')">A | B</el-button>
                <el-button @click="bitOperation('xor')">A ^ B</el-button>
                <el-button @click="bitOperation('not')">~A</el-button>
                <el-button @click="bitOperation('lshift')">A << 1</el-button>
                <el-button @click="bitOperation('rshift')">A >> 1</el-button>
              </el-button-group>

              <div v-if="bitResult" style="margin-top: 15px">
                <el-alert :title="`结果: ${bitResult}`" type="success" :closable="false" />
              </div>
            </el-tab-pane>

            <el-tab-pane label="数据大小">
              <el-input v-model="dataSize" placeholder="输入字节数">
                <template #append>Bytes</template>
              </el-input>

              <el-button @click="convertDataSize" style="margin-top: 10px; width: 100%">
                转换单位
              </el-button>

              <div v-if="sizeResult" style="margin-top: 15px">
                <el-descriptions :column="1" border>
                  <el-descriptions-item label="Bytes">{{ sizeResult.bytes }}</el-descriptions-item>
                  <el-descriptions-item label="KB">{{ sizeResult.kb }}</el-descriptions-item>
                  <el-descriptions-item label="MB">{{ sizeResult.mb }}</el-descriptions-item>
                  <el-descriptions-item label="GB">{{ sizeResult.gb }}</el-descriptions-item>
                  <el-descriptions-item label="TB">{{ sizeResult.tb }}</el-descriptions-item>
                </el-descriptions>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { create, all } from 'mathjs'

const math = create(all)

// 基础计算器
const expression = ref('')
const result = ref(null)
const error = ref('')

const insertText = (text) => {
  expression.value += text
}

const calculate = () => {
  try {
    error.value = ''
    // 替换常见数学符号
    let expr = expression.value
      .replace(/ln/g, 'log')
      .replace(/\^/g, '^')

    result.value = math.evaluate(expr)
    if (typeof result.value === 'number') {
      result.value = parseFloat(result.value.toFixed(10))
    }
  } catch (e) {
    error.value = '计算错误: ' + e.message
    result.value = null
  }
}

// 矩阵计算
const matrixA = ref('')
const matrixB = ref('')
const matrixResult = ref('')

const matrixOperation = (op) => {
  try {
    const A = math.matrix(JSON.parse(matrixA.value))

    switch (op) {
      case 'add':
        const B_add = math.matrix(JSON.parse(matrixB.value))
        matrixResult.value = JSON.stringify(math.add(A, B_add).toArray(), null, 2)
        break
      case 'subtract':
        const B_sub = math.matrix(JSON.parse(matrixB.value))
        matrixResult.value = JSON.stringify(math.subtract(A, B_sub).toArray(), null, 2)
        break
      case 'multiply':
        const B_mul = math.matrix(JSON.parse(matrixB.value))
        matrixResult.value = JSON.stringify(math.multiply(A, B_mul).toArray(), null, 2)
        break
      case 'transpose':
        matrixResult.value = JSON.stringify(math.transpose(A).toArray(), null, 2)
        break
      case 'det':
        matrixResult.value = `行列式 = ${math.det(A)}`
        break
      case 'inv':
        matrixResult.value = JSON.stringify(math.inv(A).toArray(), null, 2)
        break
    }
  } catch (e) {
    ElMessage.error('矩阵计算错误: ' + e.message)
    matrixResult.value = ''
  }
}

// 统计计算
const dataSet = ref('')
const statsResult = ref(null)

const calculateStats = () => {
  try {
    const data = dataSet.value.split(',').map(x => parseFloat(x.trim())).filter(x => !isNaN(x))

    if (data.length === 0) {
      ElMessage.warning('请输入有效的数据')
      return
    }

    statsResult.value = {
      mean: math.mean(data).toFixed(4),
      median: math.median(data).toFixed(4),
      std: math.std(data).toFixed(4),
      variance: math.variance(data).toFixed(4),
      min: math.min(data),
      max: math.max(data),
      sum: math.sum(data),
      count: data.length
    }
  } catch (e) {
    ElMessage.error('统计计算错误: ' + e.message)
  }
}

// 组合排列
const combN = ref('')
const combR = ref('')
const permN = ref('')
const permR = ref('')
const probResult = ref('')

const calcCombination = () => {
  try {
    const n = parseInt(combN.value)
    const r = parseInt(combR.value)
    const result = math.combinations(n, r)
    probResult.value = `C(${n},${r}) = ${result}`
  } catch (e) {
    ElMessage.error('计算错误: ' + e.message)
  }
}

const calcPermutation = () => {
  try {
    const n = parseInt(permN.value)
    const r = parseInt(permR.value)
    const result = math.permutations(n, r)
    probResult.value = `P(${n},${r}) = ${result}`
  } catch (e) {
    ElMessage.error('计算错误: ' + e.message)
  }
}

// 进制转换
const numberInput = ref('')
const inputBase = ref(10)
const baseResult = ref(null)

const convertBase = () => {
  try {
    const num = parseInt(numberInput.value, inputBase.value)
    baseResult.value = {
      bin: num.toString(2),
      oct: num.toString(8),
      dec: num.toString(10),
      hex: num.toString(16).toUpperCase()
    }
  } catch (e) {
    ElMessage.error('转换错误: ' + e.message)
  }
}

// 位运算
const bitA = ref('')
const bitB = ref('')
const bitResult = ref('')

const bitOperation = (op) => {
  try {
    const a = parseInt(bitA.value)
    const b = parseInt(bitB.value)

    let result
    switch (op) {
      case 'and':
        result = a & b
        break
      case 'or':
        result = a | b
        break
      case 'xor':
        result = a ^ b
        break
      case 'not':
        result = ~a
        break
      case 'lshift':
        result = a << 1
        break
      case 'rshift':
        result = a >> 1
        break
    }

    bitResult.value = `${result} (二进制: ${result.toString(2)})`
  } catch (e) {
    ElMessage.error('计算错误: ' + e.message)
  }
}

// 数据大小转换
const dataSize = ref('')
const sizeResult = ref(null)

const convertDataSize = () => {
  try {
    const bytes = parseFloat(dataSize.value)
    sizeResult.value = {
      bytes: bytes.toFixed(2),
      kb: (bytes / 1024).toFixed(2),
      mb: (bytes / (1024 ** 2)).toFixed(2),
      gb: (bytes / (1024 ** 3)).toFixed(6),
      tb: (bytes / (1024 ** 4)).toFixed(9)
    }
  } catch (e) {
    ElMessage.error('转换错误: ' + e.message)
  }
}
</script>

<style scoped>
.calculator-page {
  animation: fadeIn 0.5s;
}

.calculator-display {
  margin-bottom: 20px;
}

.calculator-result,
.calculator-error {
  margin: 15px 0;
}

.calculator-buttons {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.button-group h4 {
  margin: 0 0 10px 0;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
}

.button-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}

.calc-btn {
  height: 50px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.3s;
}

.calc-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.matrix-input label,
.stats-input label,
.probability-calc label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: var(--text-primary);
}

.matrix-result pre {
  background: var(--hover-bg);
  padding: 15px;
  border-radius: 6px;
  overflow-x: auto;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.probability-calc h4 {
  margin: 0 0 15px 0;
  color: var(--text-primary);
  font-size: 18px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
