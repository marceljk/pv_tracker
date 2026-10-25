<template>
  <div class="dashboard-container">
    <!-- Page Header -->
    <div class="d-flex align-center justify-space-between flex-wrap gap-3 mb-5">
      <div>
        <div class="d-flex align-center mb-1" style="gap: 10px;">
          <h1 class="font-display font-weight-extrabold text-h5 text-sm-h4 text-primary-adaptive mb-0">
            {{ sectionHeader.title }}
          </h1>
          <span class="status-pill d-none d-sm-inline-flex">
            <v-icon size="14" color="emerald">mdi-shield-check-outline</v-icon>
            <span>System Normal</span>
          </span>
        </div>
        <p class="text-body-2 text-secondary-adaptive mb-0">
          {{ sectionHeader.subtitle }}
        </p>
      </div>

      <!-- Date & Quick Tarif Badge -->
      <div class="d-flex align-center flex-wrap" style="gap: 8px;">
        <div class="status-pill py-1.5 px-3">
          <v-icon size="15" color="emerald">mdi-calendar-blank-outline</v-icon>
          <span>{{ formattedToday }}</span>
        </div>
        <div class="status-pill py-1.5 px-3 d-none d-lg-inline-flex">
          <v-icon size="15" color="amber">mdi-flash-outline</v-icon>
          <span>{{ settings.systemCapacityKwp }} kWp • {{ settings.batteryCapacityKwh }} kWh</span>
        </div>
      </div>
    </div>

    <!-- Top Executive Live KPI Ribbon -->
    <v-row v-if="liveData && (activeSection === 'overview' || activeSection === 'flow')" dense class="mb-4">
      <v-col cols="6" md="3">
        <div class="bento-card pa-3.5 h-100 d-flex align-center justify-space-between">
          <div>
            <div class="text-caption text-secondary-adaptive font-weight-medium">Solarleistung</div>
            <div class="stat-value text-h6 text-pv">{{ formatPower(liveData.pvPower, { compact: true }) }}</div>
          </div>
          <div class="kpi-orb bg-pv-soft">
            <v-icon color="emerald" size="20">mdi-white-balance-sunny</v-icon>
          </div>
        </div>
      </v-col>

      <v-col cols="6" md="3">
        <div class="bento-card pa-3.5 h-100 d-flex align-center justify-space-between">
          <div>
            <div class="text-caption text-secondary-adaptive font-weight-medium">Hausverbrauch</div>
            <div class="stat-value text-h6 text-consumption">
              {{ formatPower(liveData.powerConsumption, { compact: true }) }}
            </div>
          </div>
          <div class="kpi-orb bg-consumption-soft">
            <v-icon color="blue" size="20">mdi-home-lightning-bolt-outline</v-icon>
          </div>
        </div>
      </v-col>

      <v-col cols="6" md="3">
        <div class="bento-card pa-3.5 h-100 d-flex align-center justify-space-between">
          <div>
            <div class="text-caption text-secondary-adaptive font-weight-medium">
              Speicher ({{ liveData.batteryPercent || 0 }}%)
            </div>
            <div class="stat-value text-h6 text-battery">
              {{ formatPower(Math.abs(liveData.batteryPower || 0), { compact: true }) }}
            </div>
          </div>
          <div class="kpi-orb bg-battery-soft">
            <v-icon color="amber" size="20">mdi-battery-charging-80</v-icon>
          </div>
        </div>
      </v-col>

      <v-col cols="6" md="3">
        <div class="bento-card pa-3.5 h-100 d-flex align-center justify-space-between">
          <div>
            <div class="text-caption text-secondary-adaptive font-weight-medium">
              {{ (liveData.gridPower || 0) < 0 ? 'Netzbezug' : 'Netzeinspeisung' }}
            </div>
            <div
              :class="[
                'stat-value text-h6',
                (liveData.gridPower || 0) < 0 ? 'text-grid-out' : 'text-grid-in',
              ]"
            >
              {{ formatPower(Math.abs(liveData.gridPower || 0), { compact: true }) }}
            </div>
          </div>
          <div :class="['kpi-orb', (liveData.gridPower || 0) < 0 ? 'bg-grid-out-soft' : 'bg-grid-in-soft']">
            <v-icon :color="(liveData.gridPower || 0) < 0 ? 'rose' : 'cyan'" size="20">
              mdi-transmission-tower
            </v-icon>
          </div>
        </div>
      </v-col>
    </v-row>

    <!-- Primary Bento Row: Live Energy Flow + 48h Yield Forecast -->
    <v-row v-if="activeSection !== 'analytics'" class="mb-2">
      <v-col
        v-if="activeSection === 'overview' || activeSection === 'flow'"
        cols="12"
        :lg="activeSection === 'flow' ? 12 : 6"
      >
        <LiveCard class="h-100" />
      </v-col>

      <v-col
        v-if="activeSection === 'overview' || activeSection === 'forecast'"
        cols="12"
        :lg="activeSection === 'forecast' ? 12 : 6"
      >
        <ForecastCard class="h-100" />
      </v-col>
    </v-row>

    <!-- Period Benchmarks & Historical Analytics Section -->
    <template v-if="activeSection === 'overview' || activeSection === 'analytics'">
      <div class="d-flex align-center justify-space-between mt-6 mb-3">
        <div>
          <h2 class="font-display font-weight-bold text-subtitle-1 text-primary-adaptive mb-0">
            Energiebilanzen &amp; Autarkie-Kennzahlen
          </h2>
          <p class="text-caption text-secondary-adaptive mb-0">
            Kumulierte Erträge, Eigenverbrauchsquoten und Einsparungen
          </p>
        </div>
      </div>

      <!-- Desktop: 3 Period Cards Side-by-Side -->
      <v-row v-if="mdAndUp" class="mb-2">
        <v-col cols="12" md="4">
          <RangeSumCard
            title="Heute"
            subtitle="Tagesbilanz ab 00:00 Uhr"
            :start="today"
            :end="todayEnd"
            class="h-100"
          />
        </v-col>
        <v-col cols="12" md="4">
          <RangeSumCard
            title="Letzte 7 Tage"
            subtitle="Wochenbilanz"
            :start="oneWeekAgo"
            :end="todayEnd"
            class="h-100"
          />
        </v-col>
        <v-col cols="12" md="4">
          <RangeSumCard
            title="Aktuelles Jahr"
            subtitle="Jahresbilanz seit 1. Jan"
            :start="startOfYear"
            :end="todayEnd"
            class="h-100"
          />
        </v-col>
      </v-row>

      <!-- Mobile: Sleek Segmented Card -->
      <v-card v-else class="bento-card pa-3 mb-4">
        <div class="segmented-control w-100 mb-2 d-flex">
          <button
            type="button"
            :class="['segmented-btn flex-grow-1 justify-center', { active: activeTab === 'heute' }]"
            @click="activeTab = 'heute'"
          >
            Heute
          </button>
          <button
            type="button"
            :class="['segmented-btn flex-grow-1 justify-center', { active: activeTab === 'woche' }]"
            @click="activeTab = 'woche'"
          >
            7 Tage
          </button>
          <button
            type="button"
            :class="['segmented-btn flex-grow-1 justify-center', { active: activeTab === 'jahr' }]"
            @click="activeTab = 'jahr'"
          >
            Jahr
          </button>
        </div>

        <v-window v-model="activeTab">
          <v-window-item value="heute">
            <RangeSumCard title="Heute" :start="today" :end="todayEnd" flat />
          </v-window-item>
          <v-window-item value="woche">
            <RangeSumCard title="Letzte 7 Tage" :start="oneWeekAgo" :end="todayEnd" flat />
          </v-window-item>
          <v-window-item value="jahr">
            <RangeSumCard title="Aktuelles Jahr" :start="startOfYear" :end="todayEnd" flat />
          </v-window-item>
        </v-window>
      </v-card>

      <!-- Full-Width Interactive Custom Date Range & Trend Chart Studio -->
      <v-row class="mt-1">
        <v-col cols="12">
          <CustomRangeSumCard title="Historische Zeitraum-Analyse & Trends" />
        </v-col>
      </v-row>
    </template>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import { useDatabase, useDatabaseObject } from 'vuefire'
import { ref as dbRef, query, orderByChild, limitToLast } from 'firebase/database'
import ForecastCard from '@/components/ForecastCard.vue'
import LiveCard from '@/components/LiveCard.vue'
import RangeSumCard from '@/components/RangeSumCard.vue'
import CustomRangeSumCard from '@/components/CustomRangeSumCard.vue'
import { useSettings, formatPower } from '@/composables/useSettings'

const { mdAndUp } = useDisplay()
const { settings, activeSection } = useSettings()
const activeTab = ref('heute')

const db = useDatabase()
const liveData = useDatabaseObject(
  query(dbRef(db, 'live'), orderByChild('batteryPercent'), limitToLast(10))
)

const today = computed(() => new Date().toISOString().substring(0, 11))
const todayEnd = computed(() => `${new Date().toISOString().substring(0, 11)}\uf8ff`)
const oneWeekAgo = computed(() =>
  new Date(new Date().setDate(new Date().getDate() - 7)).toISOString().substring(0, 11)
)
const startOfYear = computed(() =>
  new Date(new Date().setMonth(0, 1)).toISOString().substring(0, 11)
)

const formattedToday = computed(() =>
  new Date().toLocaleDateString('de-DE', {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
)

const sectionHeader = computed(() => {
  const map = {
    overview: {
      title: 'Energie-Leitstand',
      subtitle: 'Echtzeit-Telemetrie, Ertragsprognosen und historische Autarkie-Analysen',
    },
    flow: {
      title: 'Live-Energiefluss',
      subtitle: 'Interaktive Topologie und Echtzeit-Leistungsverteilung',
    },
    forecast: {
      title: '48h Ertragsprognose',
      subtitle: 'Solcast Strahlungsmodell und erwartete Spitzenleistungen',
    },
    analytics: {
      title: 'Historie & Analyse',
      subtitle: 'Zeitraum-Bilanzen, Autarkiequoten, Finanzertrag und Datenexport',
    },
  }
  return map[activeSection.value] || map.overview
})
</script>

<style scoped>
.pa-3\.5 {
  padding: 14px 16px !important;
}

.kpi-orb {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
</style>
