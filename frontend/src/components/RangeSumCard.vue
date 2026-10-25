<template>
  <v-card
    :class="[
      flat
        ? 'elevation-0 bg-transparent border-0 pa-2'
        : 'bento-card bento-card-interactive h-100 d-flex flex-column justify-space-between pa-4 pa-sm-5',
    ]"
  >
    <!-- Header -->
    <div>
      <div v-if="!flat" class="d-flex align-center justify-space-between mb-3">
        <div class="d-flex align-center">
          <div class="period-icon-box mr-2.5">
            <v-icon color="secondary" size="18">{{ headerIcon }}</v-icon>
          </div>
          <div>
            <h3 class="font-display font-weight-bold text-subtitle-2 text-primary-adaptive mb-0">
              {{ title }}
            </h3>
            <span v-if="subtitle" class="text-caption text-muted-adaptive" style="font-size: 0.68rem !important;">
              {{ subtitle }}
            </span>
          </div>
        </div>

        <!-- Autarkie Pill -->
        <span
          v-if="dailySumList"
          class="status-pill py-0 px-2"
          style="font-size: 0.7rem;"
          title="Autarkiegrad im Zeitraum"
        >
          <span class="text-pv font-weight-bold">{{ metrics.autarkiePercent }}%</span>
          <span>Autark</span>
        </span>
      </div>

      <!-- Autarkie / Eigenverbrauch Mini Progress Bar -->
      <div v-if="dailySumList" class="surface-tile pa-2.5 mb-3">
        <div class="d-flex align-center justify-space-between mb-1">
          <span class="text-caption text-secondary-adaptive" style="font-size: 0.7rem !important;">
            Autarkie <strong class="text-pv">{{ metrics.autarkiePercent }}%</strong>
          </span>
          <span class="text-caption text-secondary-adaptive" style="font-size: 0.7rem !important;">
            Eigenverbrauch <strong class="text-consumption">{{ metrics.eigenverbrauchPercent }}%</strong>
          </span>
        </div>
        <div class="metric-progress-track" style="height: 5px;">
          <div
            class="metric-progress-fill"
            :style="{ width: `${metrics.autarkiePercent}%`, background: 'var(--color-pv)' }"
          />
        </div>
      </div>

      <!-- 4 Energy Metrics List -->
      <div v-if="dailySumList">
        <div
          v-for="row in metricRows"
          :key="row.key"
          class="d-flex align-center justify-space-between py-2 metric-row"
        >
          <div class="d-flex align-center">
            <div :class="['metric-badge mr-2.5', row.badgeClass]">
              <v-icon :color="row.iconColor" size="15">{{ row.icon }}</v-icon>
            </div>
            <span class="text-caption font-weight-medium text-secondary-adaptive">
              {{ row.label }}
            </span>
          </div>
          <span :class="['stat-value text-body-2', row.textClass]">
            {{ formatEnergy(row.value) }}
          </span>
        </div>
      </div>

      <!-- Loading State -->
      <div v-else class="d-flex justify-center align-center py-8">
        <v-progress-circular indeterminate size="24" width="2" color="emerald" />
      </div>
    </div>

    <!-- Bottom Financial & CO2 Strip -->
    <div v-if="dailySumList" class="period-footer mt-3 pt-2.5 d-flex align-center justify-space-between">
      <div class="d-flex align-center" style="gap: 5px;" title="Ersparter Netzbezug + Einspeisevergütung">
        <v-icon size="14" color="emerald">mdi-cash-multiple</v-icon>
        <span class="text-caption text-muted-adaptive" style="font-size: 0.7rem !important;">Ertrag:</span>
        <span class="stat-value text-caption text-pv">{{ formattedSavings }} €</span>
      </div>

      <div class="d-flex align-center" style="gap: 4px;" title="Eingesparte CO2-Emissionen">
        <v-icon size="14" color="secondary">mdi-leaf</v-icon>
        <span class="stat-value text-caption text-secondary-adaptive">{{ formattedCo2 }} kg CO₂</span>
      </div>
    </div>
  </v-card>
</template>

<script setup>
import { computed } from 'vue'
import { useDatabase, useDatabaseList } from 'vuefire'
import { ref as dbRef, query, orderByKey, startAt, endAt } from 'firebase/database'
import { useSettings, formatEnergy, computePeriodMetrics } from '@/composables/useSettings'

const props = defineProps({
  start: {
    type: String,
    required: true,
  },
  end: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    default: '',
  },
  subtitle: {
    type: String,
    default: '',
  },
  flat: {
    type: Boolean,
    default: false,
  },
})

const { settings } = useSettings()
const db = useDatabase()

const queryRef = computed(() =>
  query(dbRef(db, 'dailySum'), orderByKey(), startAt(props.start), endAt(props.end))
)

const dailySumList = useDatabaseList(queryRef)

const metrics = computed(() => computePeriodMetrics(dailySumList.value, settings.value))

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

const headerIcon = computed(() => {
  if (props.title === 'Heute') return 'mdi-calendar-today'
  if (props.title.includes('Woche') || props.title.includes('7')) return 'mdi-calendar-week'
  if (props.title.includes('Monat') || props.title.includes('30')) return 'mdi-calendar-month'
  if (props.title.includes('Jahr')) return 'mdi-calendar-range'
  return 'mdi-calendar-clock'
})
</script>

<style scoped>
.period-icon-box {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  background: var(--bg-subtle);
  border: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
}

.metric-badge {
  width: 28px;
  height: 28px;
  border-radius: 8px;
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

.period-footer {
  border-top: 1px solid var(--border-subtle);
}
</style>
