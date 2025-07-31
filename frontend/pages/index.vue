<template>
    <div class="p-6">
      <h1 class="text-3xl font-bold mb-4">TradeLens Dashboard</h1>
      <div v-if="loading">Loading...</div>
      <table v-else class="border-collapse border border-gray-400 w-1/2">
        <thead>
          <tr>
            <th class="border border-gray-300 p-2">Symbol</th>
            <th class="border border-gray-300 p-2">Price</th>
            <th class="border border-gray-300 p-2">RSI</th>
            <th class="border border-gray-300 p-2">EMA12</th>
            <th class="border border-gray-300 p-2">EMA26</th>
            <th class="border border-gray-300 p-2">Signal</th>
            <th class="border border-gray-300 p-2">Confidence</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>{{ summary.symbol }}</td>
            <td>{{ summary.price.price.toFixed(2) }}</td>
            <td>{{ summary.indicators.rsi.toFixed(2) }}</td>
            <td>{{ summary.indicators.ema_short.toFixed(2) }}</td>
            <td>{{ summary.indicators.ema_long.toFixed(2) }}</td>
            <td>{{ summary.signal.signal_type }}</td>
            <td>{{ summary.signal.confidence.toFixed(0) }}%</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script setup lang="ts">
  const config = useRuntimeConfig()
  const summary = ref<any>(null)
  const loading = ref(true)
  
  async function fetchSummary() {
    loading.value = true
    try {
      const { data } = await useFetch(`${config.public.apiUrl}/summary?symbol=BTCUSDT`)
      summary.value = data.value
    } catch (error) {
      console.error(error)
    }
    loading.value = false
  }
  
  onMounted(() => {
    fetchSummary()
    setInterval(fetchSummary, 10000) // auto-refresh every 10s
  })
  </script>
  