<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'

interface Metrics {
  cpu: number
  mem: number
  time: number
}

const maxPoints = 60
const cpuTimes = ref<string[]>([])
const cpuValues = ref<number[]>([])
const memValues = ref<number[]>([])
const intervalSec = ref(1)
const connected = ref(false)
const chartRef = ref<HTMLDivElement | null>(null)
let chart: echarts.ECharts | null = null
let eventSource: EventSource | null = null

function pushPoint(m: Metrics) {
  const t = new Date(m.time).toLocaleTimeString('zh-CN', { hour12: false })
  cpuTimes.value.push(t)
  cpuValues.value.push(m.cpu)
  memValues.value.push(m.mem)
  if (cpuTimes.value.length > maxPoints) {
    cpuTimes.value.shift()
    cpuValues.value.shift()
    memValues.value.shift()
  }
}

function updateChart() {
  if (!chart) return
  chart.setOption({
    title: { text: 'CPU / 内存 实时曲线' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['CPU %', '内存 %'] },
    xAxis: { type: 'category', data: cpuTimes.value },
    yAxis: { type: 'value', min: 0, max: 100, axisLabel: { formatter: '{value}%' } },
    series: [
      { name: 'CPU %', type: 'line', smooth: true, data: cpuValues.value },
      { name: '内存 %', type: 'line', smooth: true, data: memValues.value },
    ],
  })
}

function connectSSE() {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  eventSource = new EventSource('/api/stream')
  eventSource.onopen = () => { connected.value = true }
  eventSource.onerror = () => { connected.value = false }
  eventSource.onmessage = (e) => {
    try {
      const m = JSON.parse(e.data) as Metrics
      pushPoint(m)
      updateChart()
    } catch (_) {}
  }
}

async function setPushInterval(sec: number) {
  const res = await fetch(`/api/interval?seconds=${sec}`, { method: 'PUT' })
  if (res.ok) {
    const data = await res.json()
    intervalSec.value = data.seconds
  }
}

onMounted(() => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value)
    updateChart()
  }
  connectSSE()
})

onUnmounted(() => {
  eventSource?.close()
  chart?.dispose()
})

watch(intervalSec, (v) => setPushInterval(v))
</script>

<template>
  <div class="dashboard">
    <header>
      <h1>服务器状态监控看板</h1>
      <div class="controls">
        <label>
          推送间隔（秒）：
          <select v-model.number="intervalSec">
            <option v-for="s in [1, 2, 3, 5, 10]" :key="s" :value="s">{{ s }}s</option>
          </select>
        </label>
        <span class="status" :class="{ connected }">{{ connected ? '已连接' : '未连接' }}</span>
      </div>
    </header>
    <div ref="chartRef" class="chart"></div>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 1rem;
  max-width: 900px;
  margin: 0 auto;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}
h1 {
  font-size: 1.25rem;
  margin: 0;
}
.controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.controls select {
  padding: 0.25rem 0.5rem;
}
.status {
  font-size: 0.875rem;
  color: #666;
}
.status.connected {
  color: #0a0;
}
.chart {
  width: 100%;
  height: 360px;
}
</style>
