<template>
  <v-card class="bento-card h-100 d-flex flex-column justify-space-between pa-4 pa-sm-6">
    <!-- Card Header -->
    <div class="d-flex align-center justify-space-between flex-wrap gap-2 mb-4">
      <div class="d-flex align-center">
        <div class="header-icon-box bg-pv-soft mr-3">
          <v-icon color="emerald" size="20">mdi-weather-partly-sunny</v-icon>
        </div>
        <div>
          <div class="d-flex align-center" style="gap: 8px;">
            <h2 class="font-display font-weight-bold text-subtitle-1 text-primary-adaptive mb-0">
              Ertragsprognose
            </h2>
            <span class="status-pill py-0 px-2" style="font-size: 0.7rem;">
              Solcast 48h
            </span>
          </div>
          <p class="text-caption text-secondary-adaptive mb-0">
            KI-gestützte Solarstrahlungs- &amp; Leistungsprognose
          </p>
        </div>
      </div>

      <!-- Total 48h Forecast KPI -->
      <div class="text-right">
        <div class="text-caption text-muted-adaptive" style="font-size: 0.68rem !important;">
          Gesamtprognose
        </div>
        <div class="stat-value text-subtitle-1 text-pv" style="line-height: 1.15;">
          {{ formatKwh(totalForecastKwh) }} kWh
        </div>
      </div>
    </div>

    <!-- Day Selector Cards / Pills -->
    <div v-if="dailyForecast && dailyForecast.length" class="d-flex flex-column flex-grow-1 justify-space-between">
      <div class="day-selector-grid mb-3">
        <!-- "All 48h" Option -->
        <button
          type="button"
          :class="['day-pill-card', { active: selectedDayId === 'all' }]"
          @click="selectedDayId = 'all'"
        >
          <div class="d-flex align-center justify-space-between w-100 mb-1">
            <span class="text-caption font-weight-bold text-secondary-adaptive" style="font-size: 0.72rem !important;">
              48h Verlauf
            </span>
            <v-icon size="14" :color="selectedDayId === 'all' ? 'emerald' : 'secondary'">
              mdi-chart-bell-curve-cumulative
            </v-icon>
          </div>
          <div class="stat-value text-body-2 text-primary-adaptive">
            {{ formatKwh(totalForecastKwh) }} kWh
          </div>
          <div class="text-caption text-muted-adaptive" style="font-size: 0.65rem !important;">
            Peak {{ formatKw(overallPeakKw) }} kW
          </div>
        </button>

        <!-- Individual Forecasted Days -->
        <button
          v-for="item in dailyForecast"
          :key="item.id"
          type="button"
          :class="['day-pill-card', { active: selectedDayId === item.id }]"
          @click="selectedDayId = item.id"
        >
          <div class="d-flex align-center justify-space-between w-100 mb-1">
            <span class="text-caption font-weight-bold text-secondary-adaptive text-truncate" style="font-size: 0.72rem !important;">
              {{ formatShortDay(item.id, item.day) }}
            </span>
            <span class="mono-code text-caption text-pv" style="font-size: 0.68rem !important;">
              {{ formatDateLabel(item.id) }}
            </span>
          </div>
          <div class="stat-value text-body-2 text-pv">
            {{ formatKwh(item.estimate) }} kWh
          </div>
          <div class="metric-progress-track mt-1" style="height: 4px;">
            <div
              class="metric-progress-fill"
              :style="{
                width: `${getDayBarPercent(item.estimate)}%`,
                background: 'var(--color-pv)',
              }"
            />
          </div>
        </button>
      </div>

      <!-- Interactive Solar Curve Chart -->
      <div class="chart-stage surface-tile pa-3">
        <div class="d-flex align-center justify-space-between mb-2">
          <div class="d-flex align-center" style="gap: 6px;">
            <span class="text-caption font-weight-semibold text-primary-adaptive">
              {{ activeChartTitle }}
            </span>
            <span v-if="activePeakKw > 0" class="status-pill py-0 px-2" style="font-size: 0.68rem;">
              Spitze: {{ formatKw(activePeakKw) }} kW
            </span>
          </div>
          <span class="text-caption text-muted-adaptive" style="font-size: 0.7rem !important;">
            Leistung (kW)
          </span>
        </div>

        <div style="height: 175px; position: relative;" class="w-100">
          <Line
            v-if="activeChartData.labels.length > 0"
            :key="selectedDayId"
            :options="chartOptions"
            :data="activeChartData"
          />
          <div
            v-else
            class="d-flex flex-column align-center justify-center h-100 text-muted-adaptive"
          >
            <v-icon size="24" class="mb-1 opacity-50">mdi-weather-night</v-icon>
            <span class="text-caption">Keine Sonneneinstrahlung für diesen Zeitraum prognostiziert</span>
          </div>
        </div>
      </div>

      <!-- Bottom Insight Footer -->
      <div class="forecast-footer mt-3 pt-3 d-flex align-center justify-space-between flex-wrap gap-2">
        <div class="d-flex align-center" style="gap: 8px;">
          <v-icon color="emerald" size="16">mdi-lightning-bolt-circle</v-icon>
          <span class="text-caption text-secondary-adaptive">
            Spitzenleistung: <strong class="text-primary-adaptive">{{ formatKw(overallPeakKw) }} kW</strong>
            <span v-if="peakTimeLabel" class="text-muted-adaptive"> ({{ peakTimeLabel }})</span>
          </span>
        </div>

        <div class="text-caption text-secondary-adaptive">
          Geschätzter Ertragswert:
          <strong class="text-pv">{{ estimatedValueEuro }} €</strong>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-else class="d-flex flex-column justify-center align-center py-12 flex-grow-1">
      <v-progress-circular indeterminate color="emerald" size="36" width="3" class="mb-3" />
      <span class="text-caption text-secondary-adaptive">Lade Solcast Prognosedaten...</span>
    </div>
  </v-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useDatabase, useDatabaseList } from 'vuefire'
import { ref as dbRef, query, orderByKey } from 'firebase/database'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Filler,
} from 'chart.js'
import { useSettings } from '@/composables/useSettings'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  PointElement,
  LineElement,
  CategoryScale,
  LinearScale,
  Filler
)

const { settings } = useSettings()
const selectedDayId = ref('all')

const db = useDatabase()
const dailyForecast = useDatabaseList(query(dbRef(db, 'dailyForecast'), orderByKey()))
const hourlyForecast = useDatabaseList(query(dbRef(db, 'hourlyForecast'), orderByKey()))

const totalForecastKwh = computed(() => {
  if (!dailyForecast.value) return 0
  return dailyForecast.value.reduce((acc, item) => acc + (Number(item.estimate) || 0), 0)
})

const maxDayKwh = computed(() => {
  if (!dailyForecast.value || !dailyForecast.value.length) return 1
  return Math.max(1, ...dailyForecast.value.map((d) => Number(d.estimate) || 0))
})

const getDayBarPercent = (estimate) => {
  return Math.min(100, Math.max(6, Math.round(((Number(estimate) || 0) / maxDayKwh.value) * 100)))
}

const overallPeakEntry = computed(() => {
  if (!hourlyForecast.value || !hourlyForecast.value.length) return null
  return hourlyForecast.value.reduce((best, curr) => {
    return (Number(curr.pv_estimate) || 0) > (Number(best?.pv_estimate) || 0) ? curr : best
  }, null)
})

const overallPeakKw = computed(() => Number(overallPeakEntry.value?.pv_estimate) || 0)

const peakTimeLabel = computed(() => {
  if (!overallPeakEntry.value?.period_end) return ''
  const d = new Date(overallPeakEntry.value.period_end)
  return d.toLocaleString('de-DE', {
    weekday: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
})

const estimatedValueEuro = computed(() => {
  // Assume ~60% self-consumption + 40% feed-in as blended value estimate
  const blendedRate =
    ((Number(settings.value.electricityPriceCent) || 32) * 0.6 +
      (Number(settings.value.feedInTariffCent) || 8.2) * 0.4) /
    100
  return (totalForecastKwh.value * blendedRate).toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
})

const formatKwh = (val) => {
  if (val === undefined || val === null || isNaN(val)) return '0,0'
  return Number(val).toLocaleString('de-DE', {
    minimumFractionDigits: 1,
    maximumFractionDigits: 1,
  })
}

const formatKw = (val) => {
  if (val === undefined || val === null || isNaN(val)) return '0,0'
  return Number(val).toLocaleString('de-DE', {
    minimumFractionDigits: 1,
    maximumFractionDigits: 2,
  })
}

const formatShortDay = (id, fallbackDay) => {
  const num = Number(id)
  if (!isNaN(num) && num > 1000000000000) {
    const target = new Date(num)
    const now = new Date()
    if (target.toDateString() === now.toDateString()) return 'Heute'
    const tomorrow = new Date()
    tomorrow.setDate(now.getDate() + 1)
    if (target.toDateString() === tomorrow.toDateString()) return 'Morgen'
    return target.toLocaleDateString('de-DE', { weekday: 'short' })
  }
  return fallbackDay || ''
}

const formatDateLabel = (id) => {
  const num = Number(id)
  if (!isNaN(num) && num > 1000000000000) {
    return new Date(num).toLocaleDateString('de-DE', {
      day: '2-digit',
      month: '2-digit',
    })
  }
  return ''
}

const filteredHourlyPoints = computed(() => {
  if (!hourlyForecast.value) return []
  if (selectedDayId.value === 'all') {
    return hourlyForecast.value.filter((item) => Number(item.pv_estimate) > 0)
  }
  const targetDateStr = new Date(Number(selectedDayId.value)).toLocaleDateString()
  return hourlyForecast.value.filter(
    (item) =>
      new Date(item.period_end).toLocaleDateString() === targetDateStr &&
      Number(item.pv_estimate) > 0
  )
})

const activePeakKw = computed(() => {
  if (!filteredHourlyPoints.value.length) return 0
  return Math.max(...filteredHourlyPoints.value.map((p) => Number(p.pv_estimate) || 0))
})

const activeChartTitle = computed(() => {
  if (selectedDayId.value === 'all') return '48-Stunden Leistungsverlauf'
  const num = Number(selectedDayId.value)
  if (!isNaN(num) && num > 1000000000000) {
    return new Date(num).toLocaleDateString('de-DE', {
      weekday: 'long',
      day: '2-digit',
      month: 'long',
    })
  }
  return 'Tagesverlauf'
})

const activeChartData = computed(() => {
  const points = filteredHourlyPoints.value
  const labels = points.map((item) => {
    const d = new Date(item.period_end)
    const time = d.toLocaleTimeString('de-DE', { hour: '2-digit', minute: '2-digit' })
    if (selectedDayId.value === 'all') {
      const wd = d.toLocaleDateString('de-DE', { weekday: 'short' })
      return `${wd} ${time}`
    }
    return time
  })
  const values = points.map((item) => Number(item.pv_estimate) || 0)

  return {
    labels,
    datasets: [
      {
        data: values,
        label: 'Prognose (kW)',
        backgroundColor: (ctx) => {
          const canvas = ctx.chart?.ctx
          if (!canvas) return 'rgba(16, 185, 129, 0.16)'
          const gradient = canvas.createLinearGradient(0, 0, 0, 170)
          gradient.addColorStop(0, 'rgba(16, 185, 129, 0.36)')
          gradient.addColorStop(1, 'rgba(16, 185, 129, 0.0)')
          return gradient
        },
        borderColor: '#10b981',
        borderWidth: 2.2,
        tension: 0.4,
        fill: true,
        pointBackgroundColor: '#10b981',
        pointBorderColor: '#090a0f',
        pointBorderWidth: 1.5,
        pointRadius: points.length > 28 ? 0 : 2,
        pointHoverRadius: 5,
      },
    ],
  }
})

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index',
    intersect: false,
  },
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      backgroundColor: 'rgba(15, 23, 42, 0.95)',
      titleColor: '#f8fafc',
      bodyColor: '#10b981',
      borderColor: 'rgba(255, 255, 255, 0.12)',
      borderWidth: 1,
      padding: 10,
      cornerRadius: 10,
      displayColors: false,
      callbacks: {
        label: (ctx) => `Erwartete Leistung: ${formatKw(ctx.parsed.y)} kW`,
      },
    },
  },
  scales: {
    x: {
      grid: {
        display: false,
      },
      ticks: {
        color: '#64748b',
        maxTicksLimit: 7,
        font: {
          family: 'Inter',
          size: 10,
        },
      },
    },
    y: {
      beginAtZero: true,
      grid: {
        color: 'rgba(148, 163, 184, 0.08)',
      },
      ticks: {
        color: '#64748b',
        maxTicksLimit: 5,
        font: {
          family: 'Inter',
          size: 10,
        },
      },
    },
  },
}))
</script>

<style scoped>
.header-icon-box {
  width: 38px;
  height: 38px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.day-selector-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(105px, 1fr));
  gap: 8px;
}

.day-pill-card {
  appearance: none;
  background: var(--bg-subtle);
  border: 1px solid var(--border-subtle);
  border-radius: 13px;
  padding: 10px 12px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.day-pill-card:hover {
  background: var(--bg-subtle-hover);
  border-color: var(--border-strong);
}

.day-pill-card.active {
  background: var(--color-pv-soft);
  border-color: var(--color-pv-border);
  box-shadow: 0 4px 16px -4px rgba(16, 185, 129, 0.2);
}

.forecast-footer {
  border-top: 1px solid var(--border-subtle);
}
</style>
