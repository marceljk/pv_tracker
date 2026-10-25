<template>
  <v-card
    :class="[
      flat
        ? 'elevation-0 bg-transparent border-0 pa-2'
        : 'bento-card h-100 d-flex flex-column justify-space-between pa-4 pa-sm-6',
    ]"
  >
    <!-- Top Header + Quick Presets -->
    <div>
      <div class="d-flex align-center justify-space-between flex-wrap gap-2 mb-4">
        <div class="d-flex align-center">
          <div class="header-icon-box bg-consumption-soft mr-3">
            <v-icon color="blue" size="20">mdi-chart-timeline-variant-shimmer</v-icon>
          </div>
          <div>
            <div class="d-flex align-center" style="gap: 8px;">
              <h3 class="font-display font-weight-bold text-subtitle-1 text-primary-adaptive mb-0">
                {{ title || 'Historische Zeitraum-Analyse' }}
              </h3>
              <span v-if="dailyCount > 0" class="status-pill py-0 px-2" style="font-size: 0.7rem;">
                {{ dailyCount }} {{ dailyCount === 1 ? 'Tag' : 'Tage' }}
              </span>
            </div>
            <p class="text-caption text-secondary-adaptive mb-0">
              Tagesgenaue Auswertung, Autarkie-Trends &amp; CSV-Datenexport
            </p>
          </div>
        </div>

        <!-- Preset Range Buttons & CSV Export -->
        <div class="d-flex align-center flex-wrap" style="gap: 8px;">
          <div class="segmented-control">
            <button
              v-for="preset in presets"
              :key="preset.id"
              type="button"
              :class="['segmented-btn', { active: activePreset === preset.id }]"
              @click="applyPreset(preset.id)"
            >
              {{ preset.label }}
            </button>
          </div>

          <v-btn
            v-if="dailyChartPoints.length > 0"
            variant="tonal"
            size="small"
            color="emerald"
            class="text-none font-weight-semibold rounded-lg"
            prepend-icon="mdi-download-outline"
            @click="exportCsv"
          >
            CSV
          </v-btn>
        </div>
      </div>

      <v-progress-linear
        v-if="isLoading"
        color="emerald"
        height="2"
        indeterminate
        class="mb-3 rounded"
      />

      <!-- Date Inputs + View Layout -->
      <v-row dense class="mb-4">
        <v-col cols="6" sm="3">
          <v-text-field
            v-model="startDate"
            label="Startdatum"
            type="date"
            variant="solo-filled"
            density="compact"
            hide-details
            class="date-input"
            @update:model-value="activePreset = 'custom'"
          />
        </v-col>
        <v-col cols="6" sm="3">
          <v-text-field
            v-model="endDate"
            label="Enddatum"
            type="date"
            variant="solo-filled"
            density="compact"
            hide-details
            class="date-input"
            @update:model-value="activePreset = 'custom'"
          />
        </v-col>

        <!-- Summary Efficiency Pills -->
        <v-col cols="6" sm="3">
          <div class="surface-tile py-2 px-3 h-100 d-flex flex-column justify-center">
            <span class="text-caption text-muted-adaptive" style="font-size: 0.68rem !important;">
              Autarkie / Eigenverbrauch
            </span>
            <div class="d-flex align-center" style="gap: 8px;">
              <span class="stat-value text-body-2 text-pv">{{ metrics.autarkiePercent }}%</span>
              <span class="text-muted-adaptive">/</span>
              <span class="stat-value text-body-2 text-consumption">{{ metrics.eigenverbrauchPercent }}%</span>
            </div>
          </div>
        </v-col>

        <v-col cols="6" sm="3">
          <div class="surface-tile py-2 px-3 h-100 d-flex flex-column justify-center">
            <span class="text-caption text-muted-adaptive" style="font-size: 0.68rem !important;">
              Finanzieller Ertrag
            </span>
            <div class="d-flex align-center justify-space-between">
              <span class="stat-value text-body-2 text-pv">{{ formattedSavings }} €</span>
              <span class="text-caption text-secondary-adaptive" style="font-size: 0.68rem !important;">
                {{ formattedCo2 }} kg CO₂
              </span>
            </div>
          </div>
        </v-col>
      </v-row>

      <!-- Main Body: Left 4 Metric Cards + Right Daily Bar Chart -->
      <v-row v-if="hasRange" dense>
        <v-col cols="12" md="5" class="pr-md-3">
          <div class="surface-tile h-100 d-flex flex-column justify-space-between pa-3">
            <div
              v-for="row in metricRows"
              :key="row.key"
              class="d-flex align-center justify-space-between py-2 metric-row"
            >
              <div class="d-flex align-center">
                <div :class="['metric-badge', row.badgeClass]">
                  <v-icon :color="row.iconColor" size="16">{{ row.icon }}</v-icon>
                </div>
                <div>
                  <div class="text-caption font-weight-semibold text-primary-adaptive">
                    {{ row.label }}
                  </div>
                  <div class="text-caption text-muted-adaptive" style="font-size: 0.66rem !important;">
                    Ø {{ formatEnergy(row.value / Math.max(1, dailyCount)) }} / Tag
                  </div>
                </div>
              </div>
              <span :class="['stat-value text-subtitle-2', row.textClass]">
                {{ formatEnergy(row.value) }}
              </span>
            </div>
          </div>
        </v-col>

        <!-- Daily Trend Chart -->
        <v-col cols="12" md="7" class="mt-3 mt-md-0">
          <div class="surface-tile h-100 pa-3 d-flex flex-column justify-space-between">
            <div class="d-flex align-center justify-space-between mb-2">
              <span class="text-caption font-weight-semibold text-primary-adaptive">
                Tagesvergleich (kWh)
              </span>
              <div class="d-flex align-center" style="gap: 12px; font-size: 0.7rem;">
                <span class="d-inline-flex align-center" style="gap: 4px;">
                  <span style="width: 8px; height: 8px; border-radius: 2px; background: #10b981;"></span>
                  <span class="text-secondary-adaptive">PV Erzeugung</span>
                </span>
                <span class="d-inline-flex align-center" style="gap: 4px;">
                  <span style="width: 8px; height: 8px; border-radius: 2px; background: #3b82f6;"></span>
                  <span class="text-secondary-adaptive">Verbrauch</span>
                </span>
              </div>
            </div>

            <div style="height: 185px; position: relative;" class="w-100">
              <Bar
                v-if="dailyChartPoints.length > 0"
                :options="barChartOptions"
                :data="barChartData"
              />
              <div
                v-else
                class="d-flex flex-column align-center justify-center h-100 text-muted-adaptive"
              >
                <v-icon size="24" class="mb-1 opacity-50">mdi-database-search-outline</v-icon>
                <span class="text-caption">Keine Messdaten im gewählten Zeitraum vorhanden</span>
              </div>
            </div>
          </div>
        </v-col>
      </v-row>

      <div v-else class="d-flex flex-column align-center justify-center py-8 text-muted-adaptive">
        <v-icon size="28" class="mb-2 opacity-40">mdi-calendar-range</v-icon>
        <span class="text-caption font-weight-medium">Bitte Start- und Enddatum auswählen</span>
      </div>
    </div>
  </v-card>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useDatabase, useDatabaseList } from 'vuefire'
import { ref as dbRef, query, orderByKey, startAt, endAt } from 'firebase/database'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
} from 'chart.js'
import { useSettings, formatEnergy, computePeriodMetrics } from '@/composables/useSettings'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

defineProps({
  title: {
    type: String,
    default: 'Historische Zeitraum-Analyse',
  },
  flat: {
    type: Boolean,
    default: false,
  },
})

const { settings } = useSettings()
const db = useDatabase()

const toIsoDate = (d) => d.toISOString().substring(0, 10)

const activePreset = ref('14d')
const endDate = ref(toIsoDate(new Date()))
const startDate = ref(toIsoDate(new Date(new Date().setDate(new Date().getDate() - 13))))
const isLoading = ref(false)

const presets = [
  { id: '7d', label: '7T' },
  { id: '14d', label: '14T' },
  { id: '30d', label: '30T' },
  { id: 'month', label: 'Monat' },
  { id: 'year', label: 'Jahr' },
]

const applyPreset = (id) => {
  activePreset.value = id
  const now = new Date()
  endDate.value = toIsoDate(now)
  if (id === '7d') {
    const start = new Date()
    start.setDate(now.getDate() - 6)
    startDate.value = toIsoDate(start)
  } else if (id === '14d') {
    const start = new Date()
    start.setDate(now.getDate() - 13)
    startDate.value = toIsoDate(start)
  } else if (id === '30d') {
    const start = new Date()
    start.setDate(now.getDate() - 29)
    startDate.value = toIsoDate(start)
  } else if (id === 'month') {
    const start = new Date(now.getFullYear(), now.getMonth(), 1)
    startDate.value = toIsoDate(start)
  } else if (id === 'year') {
    const start = new Date(now.getFullYear(), 0, 1)
    startDate.value = toIsoDate(start)
  }
}

const hasRange = computed(() => Boolean(startDate.value && endDate.value))

const queryDb = computed(() => {
  if (startDate.value && endDate.value) {
    return query(
      dbRef(db, 'dailySum'),
      orderByKey(),
      startAt(startDate.value),
      endAt(`${endDate.value}T\uf8ff`)
    )
  }
  return null
})

const dataFromDb = useDatabaseList(queryDb)

watch(queryDb, (newQuery) => {
  isLoading.value = Boolean(newQuery)
})

watch(dataFromDb, () => {
  isLoading.value = false
})

const dailyCount = computed(() => (dataFromDb.value ? dataFromDb.value.length : 0))
const metrics = computed(() => computePeriodMetrics(dataFromDb.value, settings.value))

const metricRows = computed(() => [
  {
    key: 'pvPower',
    label: 'PV Erzeugung',
    value: metrics.value.pvWh,
    icon: 'mdi-solar-power-variant',
    iconColor: 'emerald',
    badgeClass: 'bg-pv-soft',
    textClass: 'text-pv',
  },
  {
    key: 'powerConsumption',
    label: 'Stromverbrauch',
    value: metrics.value.consumptionWh,
    icon: 'mdi-home-lightning-bolt-outline',
    iconColor: 'blue',
    badgeClass: 'bg-consumption-soft',
    textClass: 'text-consumption',
  },
  {
    key: 'gridPowerIn',
    label: 'Netzeinspeisung',
    value: metrics.value.gridExportWh,
    icon: 'mdi-transmission-tower-export',
    iconColor: 'cyan',
    badgeClass: 'bg-grid-in-soft',
    textClass: 'text-grid-in',
  },
  {
    key: 'gridPowerOut',
    label: 'Netzbezug',
    value: metrics.value.gridImportWh,
    icon: 'mdi-transmission-tower-import',
    iconColor: 'rose',
    badgeClass: 'bg-grid-out-soft',
    textClass: 'text-grid-out',
  },
])

const formattedSavings = computed(() =>
  metrics.value.savingsEuro.toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
)

const formattedCo2 = computed(() =>
  metrics.value.co2SavedKg.toLocaleString('de-DE', {
    minimumFractionDigits: 1,
    maximumFractionDigits: 1,
  })
)

const dailyChartPoints = computed(() => {
  if (!dataFromDb.value || !dataFromDb.value.length) return []
  return dataFromDb.value.map((day) => {
    const count = day.count || 1
    const measureTime = count / 120
    const rawKey = String(day.id || '').replace(/T$/, '')
    const pvKwh = (((day.pvPower || 0) / count) * measureTime) / 1000
    const consKwh = (((day.powerConsumption || 0) / count) * measureTime) / 1000
    const exportKwh = ((Math.abs(day.gridPowerIn || 0) / count) * measureTime) / 1000
    const importKwh = ((Math.abs(day.gridPowerOut || 0) / count) * measureTime) / 1000

    let label = rawKey
    const parsed = new Date(rawKey)
    if (!isNaN(parsed.getTime())) {
      label = parsed.toLocaleDateString('de-DE', { day: '2-digit', month: '2-digit' })
    }

    return {
      date: rawKey,
      label,
      pvKwh: Math.max(0, Number(pvKwh.toFixed(2))),
      consKwh: Math.max(0, Number(consKwh.toFixed(2))),
      exportKwh: Math.max(0, Number(exportKwh.toFixed(2))),
      importKwh: Math.max(0, Number(importKwh.toFixed(2))),
    }
  })
})

const barChartData = computed(() => ({
  labels: dailyChartPoints.value.map((p) => p.label),
  datasets: [
    {
      label: 'PV Erzeugung (kWh)',
      data: dailyChartPoints.value.map((p) => p.pvKwh),
      backgroundColor: '#10b981',
      borderRadius: 5,
      maxBarThickness: 18,
    },
    {
      label: 'Stromverbrauch (kWh)',
      data: dailyChartPoints.value.map((p) => p.consKwh),
      backgroundColor: '#3b82f6',
      borderRadius: 5,
      maxBarThickness: 18,
    },
  ],
}))

const barChartOptions = {
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
      bodyColor: '#cbd5e1',
      borderColor: 'rgba(255, 255, 255, 0.12)',
      borderWidth: 1,
      padding: 10,
      cornerRadius: 10,
    },
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: {
        color: '#64748b',
        maxTicksLimit: 10,
        font: { family: 'Inter', size: 10 },
      },
    },
    y: {
      beginAtZero: true,
      grid: { color: 'rgba(148, 163, 184, 0.08)' },
      ticks: {
        color: '#64748b',
        maxTicksLimit: 5,
        font: { family: 'Inter', size: 10 },
      },
    },
  },
}

const exportCsv = () => {
  if (!dailyChartPoints.value.length) return
  const headers = ['Datum', 'PV_Erzeugung_kWh', 'Stromverbrauch_kWh', 'Netzeinspeisung_kWh', 'Netzbezug_kWh']
  const rows = dailyChartPoints.value.map((r) =>
    [r.date, r.pvKwh, r.consKwh, r.exportKwh, r.importKwh].join(',')
  )
  const csvContent = [headers.join(','), ...rows].join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', `solaros-export-${startDate.value}-bis-${endDate.value}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
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

.metric-badge {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  flex-shrink: 0;
}

.metric-row {
  border-bottom: 1px solid var(--border-subtle);
}

.metric-row:last-child {
  border-bottom: none;
}

.date-input :deep(.v-field__input) {
  padding-top: 6px !important;
  padding-bottom: 6px !important;
  font-size: 0.82rem;
}
</style>
