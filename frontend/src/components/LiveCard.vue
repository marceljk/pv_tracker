<template>
  <v-card class="bento-card h-100 d-flex flex-column justify-space-between pa-4 pa-sm-6">
    <!-- Top Header -->
    <div class="d-flex align-center justify-space-between flex-wrap gap-2 mb-4">
      <div class="d-flex align-center">
        <div class="header-icon-box bg-pv-soft mr-3">
          <v-icon color="emerald" size="20">mdi-transit-connection-variant</v-icon>
        </div>
        <div>
          <div class="d-flex align-center" style="gap: 8px;">
            <h2 class="font-display font-weight-bold text-subtitle-1 text-primary-adaptive mb-0">
              Echtzeit-Energiefluss
            </h2>
            <span class="status-pill py-0 px-2" style="font-size: 0.7rem;">
              <span class="live-dot" style="width: 6px; height: 6px;"></span>
              {{ systemStateLabel }}
            </span>
          </div>
          <p class="text-caption text-secondary-adaptive mb-0">
            Live-Telemetrie von PV, Varta-Speicher, Hausnetz &amp; Stromnetz
          </p>
        </div>
      </div>

      <!-- View Mode Switcher -->
      <div class="segmented-control">
        <button
          type="button"
          :class="['segmented-btn', { active: viewMode === 'topology' }]"
          @click="viewMode = 'topology'"
        >
          <v-icon size="15">mdi-vector-bezier</v-icon>
          <span>Topologie</span>
        </button>
        <button
          type="button"
          :class="['segmented-btn', { active: viewMode === 'cards' }]"
          @click="viewMode = 'cards'"
        >
          <v-icon size="15">mdi-view-grid-outline</v-icon>
          <span>Kacheln</span>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-if="liveData" class="flex-grow-1 d-flex flex-column justify-space-between">
      <!-- MODE 1: Interactive FANG Energy Flow Topology -->
      <div v-if="viewMode === 'topology'" class="topology-wrapper my-1">
        <div class="topology-stage">
          <!-- Animated SVG Flow Lines -->
          <svg
            class="topology-svg"
            viewBox="0 0 400 260"
            preserveAspectRatio="xMidYMid meet"
            aria-hidden="true"
          >
            <defs>
              <radialGradient id="hubGlow" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stop-color="#10b981" stop-opacity="0.28" />
                <stop offset="100%" stop-color="#10b981" stop-opacity="0" />
              </radialGradient>
            </defs>

            <!-- Base Track Paths -->
            <!-- PV (Top: 200, 42) to Hub (200, 130) -->
            <path id="path-pv-hub" d="M 200 48 L 200 130" class="flow-track" />
            <!-- Hub (200, 130) to Home (Right: 332, 130) -->
            <path id="path-hub-home" d="M 200 130 L 332 130" class="flow-track" />
            <!-- Grid (Left: 68, 130) to Hub (200, 130) -->
            <path id="path-grid-hub" d="M 68 130 L 200 130" class="flow-track" />
            <!-- Hub (200, 130) to Battery (Bottom: 200, 212) -->
            <path id="path-hub-battery" d="M 200 130 L 200 212" class="flow-track" />

            <!-- Active Coloured Paths -->
            <path
              v-if="pvWatts > 10"
              d="M 200 48 L 200 130"
              class="flow-active stroke-pv"
              :style="{ strokeDasharray: '6 8', animationDuration: flowDuration(pvWatts) }"
            />
            <path
              v-if="homeWatts > 10"
              d="M 200 130 L 332 130"
              class="flow-active stroke-consumption"
              :style="{ strokeDasharray: '6 8', animationDuration: flowDuration(homeWatts) }"
            />
            <path
              v-if="Math.abs(gridWatts) > 10"
              d="M 68 130 L 200 130"
              :class="[
                'flow-active',
                gridWatts < 0 ? 'stroke-grid-out' : 'stroke-grid-in flow-reverse',
              ]"
              :style="{ strokeDasharray: '6 8', animationDuration: flowDuration(gridWatts) }"
            />
            <path
              v-if="Math.abs(batteryWatts) > 10"
              d="M 200 130 L 200 212"
              :class="[
                'flow-active',
                batteryWatts >= 0 ? 'stroke-battery' : 'stroke-battery flow-reverse',
              ]"
              :style="{ strokeDasharray: '6 8', animationDuration: flowDuration(batteryWatts) }"
            />

            <!-- Moving Energy Particles -->
            <circle v-if="pvWatts > 10" r="4" fill="#10b981">
              <animateMotion :dur="flowDuration(pvWatts)" repeatCount="indefinite" path="M 200 48 L 200 130" />
            </circle>
            <circle v-if="homeWatts > 10" r="4" fill="#3b82f6">
              <animateMotion :dur="flowDuration(homeWatts)" repeatCount="indefinite" path="M 200 130 L 332 130" />
            </circle>
            <circle v-if="gridWatts > 10" r="4" fill="#06b6d4">
              <animateMotion :dur="flowDuration(gridWatts)" repeatCount="indefinite" path="M 200 130 L 68 130" />
            </circle>
            <circle v-if="gridWatts < -10" r="4" fill="#f43f5e">
              <animateMotion :dur="flowDuration(gridWatts)" repeatCount="indefinite" path="M 68 130 L 200 130" />
            </circle>
            <circle v-if="batteryWatts > 10" r="4" fill="#f59e0b">
              <animateMotion :dur="flowDuration(batteryWatts)" repeatCount="indefinite" path="M 200 130 L 200 212" />
            </circle>
            <circle v-if="batteryWatts < -10" r="4" fill="#f59e0b">
              <animateMotion :dur="flowDuration(batteryWatts)" repeatCount="indefinite" path="M 200 212 L 200 130" />
            </circle>

            <!-- Central Inverter Hub -->
            <circle cx="200" cy="130" r="34" fill="url(#hubGlow)" />
            <circle cx="200" cy="130" r="20" class="hub-core" />
            <circle cx="200" cy="130" r="6" fill="#10b981" />
          </svg>

          <!-- TOP NODE: PV Generation -->
          <div class="topo-node node-top">
            <div class="node-card border-pv">
              <div class="d-flex align-center justify-space-between mb-1" style="gap: 8px;">
                <span class="d-flex align-center text-caption font-weight-semibold text-secondary-adaptive">
                  <v-icon color="emerald" size="16" class="mr-1">mdi-solar-power-variant</v-icon>
                  PV-Anlage
                </span>
                <span class="mono-code text-caption text-pv font-weight-bold">{{ pvCapacityPercent }}%</span>
              </div>
              <div class="stat-value text-subtitle-1 text-sm-h6 text-pv" style="line-height: 1.15;">
                {{ formatPower(pvWatts) }}
              </div>
            </div>
          </div>

          <!-- LEFT NODE: Grid -->
          <div class="topo-node node-left">
            <div :class="['node-card', gridWatts < 0 ? 'border-grid-out' : 'border-grid-in']">
              <div class="d-flex align-center mb-1">
                <v-icon :color="gridWatts < 0 ? 'rose' : 'cyan'" size="16" class="mr-1">
                  mdi-transmission-tower
                </v-icon>
                <span class="text-caption font-weight-semibold text-secondary-adaptive text-truncate">
                  {{ gridWatts < 0 ? 'Netzbezug' : 'Einspeisung' }}
                </span>
              </div>
              <div
                :class="[
                  'stat-value text-subtitle-1 text-sm-h6',
                  gridWatts < 0 ? 'text-grid-out' : 'text-grid-in',
                ]"
                style="line-height: 1.15;"
              >
                {{ formatPower(Math.abs(gridWatts)) }}
              </div>
            </div>
          </div>

          <!-- RIGHT NODE: Home Consumption -->
          <div class="topo-node node-right">
            <div class="node-card border-consumption">
              <div class="d-flex align-center mb-1">
                <v-icon color="blue" size="16" class="mr-1">mdi-home-lightning-bolt-outline</v-icon>
                <span class="text-caption font-weight-semibold text-secondary-adaptive text-truncate">
                  Hausverbrauch
                </span>
              </div>
              <div class="stat-value text-subtitle-1 text-sm-h6 text-consumption" style="line-height: 1.15;">
                {{ formatPower(homeWatts) }}
              </div>
            </div>
          </div>

          <!-- BOTTOM NODE: Varta Battery -->
          <div class="topo-node node-bottom">
            <div class="node-card border-battery">
              <div class="d-flex align-center justify-space-between mb-1" style="gap: 8px;">
                <span class="d-flex align-center text-caption font-weight-semibold text-secondary-adaptive">
                  <v-icon color="amber" size="16" class="mr-1">
                    {{ batteryWatts < 0 ? 'mdi-battery-arrow-down-outline' : 'mdi-battery-charging-80' }}
                  </v-icon>
                  {{ batteryStatusLabel }}
                </span>
                <span class="mono-code text-caption font-weight-bold text-battery">
                  {{ batterySoc }}%
                </span>
              </div>
              <div class="stat-value text-subtitle-1 text-sm-h6 text-battery" style="line-height: 1.15;">
                {{ formatPower(Math.abs(batteryWatts)) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- MODE 2: High-Density Engineering Telemetry Cards -->
      <div v-else class="my-2">
        <v-row align="center" dense>
          <v-col cols="12" sm="8">
            <v-row dense>
              <!-- PV Power -->
              <v-col cols="6" class="pa-1">
                <div class="surface-tile h-100">
                  <div class="d-flex align-center justify-space-between mb-2">
                    <span class="text-caption font-weight-medium text-secondary-adaptive">PV Erzeugung</span>
                    <v-icon color="emerald" size="18">mdi-solar-power-variant</v-icon>
                  </div>
                  <div class="stat-value text-h6 text-pv mb-1">{{ formatPower(pvWatts) }}</div>
                  <div class="metric-progress-track">
                    <div
                      class="metric-progress-fill"
                      :style="{ width: `${pvCapacityPercent}%`, background: 'var(--color-pv)' }"
                    />
                  </div>
                  <div class="text-caption text-muted-adaptive mt-1" style="font-size: 0.68rem !important;">
                    {{ pvCapacityPercent }}% von {{ settings.systemCapacityKwp }} kWp
                  </div>
                </div>
              </v-col>

              <!-- Home Consumption -->
              <v-col cols="6" class="pa-1">
                <div class="surface-tile h-100">
                  <div class="d-flex align-center justify-space-between mb-2">
                    <span class="text-caption font-weight-medium text-secondary-adaptive">Hausverbrauch</span>
                    <v-icon color="blue" size="18">mdi-home-lightning-bolt-outline</v-icon>
                  </div>
                  <div class="stat-value text-h6 text-consumption mb-1">{{ formatPower(homeWatts) }}</div>
                  <div class="metric-progress-track">
                    <div
                      class="metric-progress-fill"
                      :style="{ width: `${liveAutarkie}%`, background: 'var(--color-consumption)' }"
                    />
                  </div>
                  <div class="text-caption text-muted-adaptive mt-1" style="font-size: 0.68rem !important;">
                    {{ liveAutarkie }}% lokal gedeckt
                  </div>
                </div>
              </v-col>

              <!-- Grid Power -->
              <v-col cols="6" class="pa-1">
                <div class="surface-tile h-100">
                  <div class="d-flex align-center justify-space-between mb-2">
                    <span class="text-caption font-weight-medium text-secondary-adaptive">
                      {{ gridWatts < 0 ? 'Netzbezug' : 'Netzeinspeisung' }}
                    </span>
                    <v-icon :color="gridWatts < 0 ? 'rose' : 'cyan'" size="18">mdi-transmission-tower</v-icon>
                  </div>
                  <div :class="['stat-value text-h6 mb-1', gridWatts < 0 ? 'text-grid-out' : 'text-grid-in']">
                    {{ formatPower(Math.abs(gridWatts)) }}
                  </div>
                  <div class="text-caption text-muted-adaptive" style="font-size: 0.68rem !important;">
                    {{ gridWatts < 0 ? 'Import aus öffentlichem Netz' : 'Überschuss ins Stromnetz' }}
                  </div>
                </div>
              </v-col>

              <!-- Battery Power -->
              <v-col cols="6" class="pa-1">
                <div class="surface-tile h-100">
                  <div class="d-flex align-center justify-space-between mb-2">
                    <span class="text-caption font-weight-medium text-secondary-adaptive">
                      {{ batteryStatusLabel }}
                    </span>
                    <v-icon color="amber" size="18">
                      {{ batteryWatts < 0 ? 'mdi-battery-arrow-down-outline' : 'mdi-battery-arrow-up' }}
                    </v-icon>
                  </div>
                  <div class="stat-value text-h6 text-battery mb-1">
                    {{ formatPower(Math.abs(batteryWatts)) }}
                  </div>
                  <div class="text-caption text-muted-adaptive" style="font-size: 0.68rem !important;">
                    {{ batteryRuntimeEstimate }}
                  </div>
                </div>
              </v-col>
            </v-row>
          </v-col>

          <!-- Right Circular Battery Gauge -->
          <v-col cols="12" sm="4" class="d-flex flex-column align-center justify-center py-3">
            <v-progress-circular
              :model-value="batterySoc"
              :size="122"
              :width="11"
              color="amber"
              bg-color="rgba(245, 158, 11, 0.14)"
            >
              <div class="d-flex flex-column align-center">
                <span class="stat-value text-h5 text-battery">{{ batterySoc }}%</span>
                <span class="text-caption text-secondary-adaptive font-weight-semibold" style="font-size: 0.7rem !important;">
                  Speicher SoC
                </span>
              </div>
            </v-progress-circular>
            <span class="text-caption text-muted-adaptive mt-2 text-center">
              {{ batteryRuntimeEstimate }}
            </span>
          </v-col>
        </v-row>
      </div>

      <!-- Bottom Derived Telemetry Strip -->
      <div class="telemetry-footer mt-3 pt-3">
        <v-row dense>
          <v-col cols="4">
            <div class="d-flex flex-column">
              <span class="text-caption text-muted-adaptive" style="font-size: 0.7rem !important;">
                Live-Autarkie
              </span>
              <div class="d-flex align-center" style="gap: 6px;">
                <span class="stat-value text-body-2 text-pv">{{ liveAutarkie }}%</span>
                <span class="text-caption text-secondary-adaptive d-none d-sm-inline" style="font-size: 0.7rem !important;">
                  Unabhängig
                </span>
              </div>
            </div>
          </v-col>

          <v-col cols="4">
            <div class="d-flex flex-column">
              <span class="text-caption text-muted-adaptive" style="font-size: 0.7rem !important;">
                Speicher-Ladestand
              </span>
              <div class="d-flex align-center" style="gap: 6px;">
                <span class="stat-value text-body-2 text-battery">{{ batterySoc }}%</span>
                <span class="text-caption text-secondary-adaptive text-truncate" style="font-size: 0.7rem !important;">
                  ({{ batteryStoredKwh }} kWh)
                </span>
              </div>
            </div>
          </v-col>

          <v-col cols="4">
            <div class="d-flex flex-column align-end text-right">
              <span class="text-caption text-muted-adaptive" style="font-size: 0.7rem !important;">
                Speicher-Prognose
              </span>
              <span class="stat-value text-body-2 text-primary-adaptive text-truncate">
                {{ batteryRuntimeEstimate }}
              </span>
            </div>
          </v-col>
        </v-row>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-else class="d-flex flex-column justify-center align-center py-12 flex-grow-1">
      <v-progress-circular indeterminate color="emerald" size="36" width="3" class="mb-3" />
      <span class="text-caption text-secondary-adaptive">Verbinde mit Echtzeit-Telemetrie...</span>
    </div>
  </v-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useDatabaseObject, useDatabase } from 'vuefire'
import { ref as dbRef, query, orderByChild, limitToLast } from 'firebase/database'
import { useSettings, formatPower } from '@/composables/useSettings'

const { settings } = useSettings()
const viewMode = ref('topology')

const db = useDatabase()
const liveData = useDatabaseObject(
  query(dbRef(db, 'live'), orderByChild('batteryPercent'), limitToLast(10))
)

const pvWatts = computed(() => Number(liveData.value?.pvPower) || 0)
const homeWatts = computed(() => Number(liveData.value?.powerConsumption) || 0)
const gridWatts = computed(() => Number(liveData.value?.gridPower) || 0)
const batteryWatts = computed(() => Number(liveData.value?.batteryPower) || 0)
const batterySoc = computed(() => Math.min(100, Math.max(0, Number(liveData.value?.batteryPercent) || 0)))

const pvCapacityPercent = computed(() => {
  const maxWatts = (Number(settings.value.systemCapacityKwp) || 10) * 1000
  if (maxWatts <= 0) return 0
  return Math.min(100, Math.max(0, Math.round((pvWatts.value / maxWatts) * 100)))
})

const liveAutarkie = computed(() => {
  if (homeWatts.value <= 0) return 100
  if (gridWatts.value >= 0) return 100
  const gridImport = Math.abs(gridWatts.value)
  return Math.min(100, Math.max(0, Math.round(((homeWatts.value - gridImport) / homeWatts.value) * 100)))
})

const batteryStoredKwh = computed(() => {
  const cap = Number(settings.value.batteryCapacityKwh) || 10
  return ((batterySoc.value / 100) * cap).toLocaleString('de-DE', {
    minimumFractionDigits: 1,
    maximumFractionDigits: 1,
  })
})

const batteryStatusLabel = computed(() => {
  if (batteryWatts.value > 15) return 'Akku lädt'
  if (batteryWatts.value < -15) return 'Akku entlädt'
  return 'Akku Standby'
})

const batteryRuntimeEstimate = computed(() => {
  const capWh = (Number(settings.value.batteryCapacityKwh) || 10) * 1000
  const power = batteryWatts.value
  if (power > 25 && batterySoc.value < 100) {
    const remainingWh = ((100 - batterySoc.value) / 100) * capWh
    const hours = remainingWh / power
    if (hours > 24) return 'Lädt (>24h)'
    const h = Math.floor(hours)
    const m = Math.round((hours - h) * 60)
    return `Voll in ~${h}h ${m}m`
  }
  if (power < -25 && batterySoc.value > 0) {
    const availableWh = (batterySoc.value / 100) * capWh
    const hours = availableWh / Math.abs(power)
    if (hours > 24) return 'Reicht >24h'
    const h = Math.floor(hours)
    const m = Math.round((hours - h) * 60)
    return `Reicht ~${h}h ${m}m`
  }
  return `Bereit (${batterySoc.value}%)`
})

const systemStateLabel = computed(() => {
  if (pvWatts.value > homeWatts.value && gridWatts.value > 20) return 'Überschuss-Einspeisung'
  if (pvWatts.value > 50) return 'Solarbetrieb'
  if (batteryWatts.value < -20) return 'Speicherbetrieb'
  if (gridWatts.value < -20) return 'Netzbezug'
  return 'Standby'
})

const flowDuration = (watts) => {
  const abs = Math.abs(Number(watts) || 0)
  if (abs > 4000) return '1.1s'
  if (abs > 1500) return '1.6s'
  if (abs > 400) return '2.2s'
  return '3.0s'
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

.topology-wrapper {
  position: relative;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.topology-stage {
  position: relative;
  width: 100%;
  height: 268px;
  max-width: 520px;
  margin: 0 auto;
}

.topology-svg {
  width: 100%;
  height: 100%;
  display: block;
}

.flow-track {
  fill: none;
  stroke: var(--border-strong);
  stroke-width: 2;
  stroke-linecap: round;
}

.flow-active {
  fill: none;
  stroke-width: 2.5;
  stroke-linecap: round;
  animation: dash-flow 2s linear infinite;
}

.flow-reverse {
  animation-direction: reverse;
}

.stroke-pv { stroke: var(--color-pv); }
.stroke-consumption { stroke: var(--color-consumption); }
.stroke-battery { stroke: var(--color-battery); }
.stroke-grid-in { stroke: var(--color-grid-in); }
.stroke-grid-out { stroke: var(--color-grid-out); }

.hub-core {
  fill: var(--bg-elevated-solid);
  stroke: var(--border-highlight);
  stroke-width: 2;
}

@keyframes dash-flow {
  from {
    stroke-dashoffset: 28;
  }
  to {
    stroke-dashoffset: 0;
  }
}

.topo-node {
  position: absolute;
  z-index: 2;
}

.node-top {
  top: 0;
  left: 50%;
  transform: translateX(-50%);
}

.node-bottom {
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
}

.node-left {
  left: 0;
  top: 50%;
  transform: translateY(-50%);
}

.node-right {
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

.node-card {
  background: var(--bg-elevated);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid var(--border-strong);
  border-radius: 14px;
  padding: 8px 13px;
  min-width: 132px;
  box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.35);
  transition: transform 0.2s ease, border-color 0.2s ease;
}

.node-card:hover {
  transform: scale(1.03);
}

.border-pv { border-left: 3px solid var(--color-pv); }
.border-consumption { border-right: 3px solid var(--color-consumption); }
.border-battery { border-left: 3px solid var(--color-battery); }
.border-grid-in { border-left: 3px solid var(--color-grid-in); }
.border-grid-out { border-left: 3px solid var(--color-grid-out); }

.telemetry-footer {
  border-top: 1px solid var(--border-subtle);
}

@media (max-width: 480px) {
  .topology-stage {
    height: 250px;
  }
  .node-card {
    min-width: 112px;
    padding: 6px 10px;
  }
}
</style>
