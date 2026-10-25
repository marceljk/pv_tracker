<template>
  <v-app-bar flat height="68" class="solar-app-bar px-3 px-md-6">
    <div class="d-flex align-center justify-space-between w-100 mx-auto" style="max-width: 1440px;">
      <!-- Left: Brand Identity -->
      <div class="d-flex align-center cursor-pointer" @click="activeSection = 'overview'">
        <div class="brand-badge mr-3 d-flex align-center justify-center">
          <v-icon color="emerald" size="22">mdi-solar-power-variant</v-icon>
        </div>
        <div>
          <div class="d-flex align-center gap-2">
            <span class="font-display font-weight-bold text-subtitle-1 text-primary-adaptive" style="line-height: 1.1;">
              SolarOS
            </span>
            <span class="version-chip">v2.4</span>
          </div>
          <div class="text-caption text-muted-adaptive d-none d-sm-block" style="font-size: 0.7rem !important; line-height: 1.15;">
            PV Telemetry &amp; Grid Intelligence
          </div>
        </div>
      </div>

      <!-- Center: Desktop Segmented View Switcher -->
      <div class="segmented-control d-none d-md-inline-flex">
        <button
          v-for="item in navItems"
          :key="item.id"
          type="button"
          :class="['segmented-btn', { active: activeSection === item.id }]"
          @click="activeSection = item.id"
        >
          <v-icon size="16">{{ item.icon }}</v-icon>
          <span>{{ item.label }}</span>
        </button>
      </div>

      <!-- Right: Live Status, Settings & Theme Toggle -->
      <div class="d-flex align-center" style="gap: 8px;">
        <!-- Live Connection Pill -->
        <div class="status-pill">
          <span :class="['live-dot', { 'offline-dot': !isOnline }]"></span>
          <span class="d-none d-sm-inline">{{ isOnline ? 'Firebase Echtzeit' : 'Offline Modus' }}</span>
          <span class="d-inline d-sm-none">{{ isOnline ? 'Live' : 'Offline' }}</span>
        </div>

        <!-- System & Tariff Parameters Dialog Button -->
        <v-btn
          icon
          variant="text"
          size="small"
          class="action-icon-btn"
          title="Anlagen- & Tarifparameter"
          @click="settingsDialog = true"
        >
          <v-icon size="19">mdi-tune-variant</v-icon>
        </v-btn>

        <!-- Theme Switcher -->
        <v-btn
          icon
          variant="text"
          size="small"
          class="action-icon-btn"
          :title="settings.theme === 'dark' ? 'Helles Design aktivieren' : 'Dunkles Design aktivieren'"
          @click="handleToggleTheme"
        >
          <v-icon size="19">
            {{ settings.theme === 'dark' ? 'mdi-weather-sunny' : 'mdi-weather-night' }}
          </v-icon>
        </v-btn>
      </div>
    </div>

    <!-- System & Tariff Settings Modal -->
    <v-dialog v-model="settingsDialog" max-width="460">
      <v-card class="bento-card pa-2">
        <v-card-title class="d-flex align-center justify-space-between pt-3 px-4">
          <div class="d-flex align-center">
            <v-icon color="emerald" class="mr-2" size="20">mdi-tune-variant</v-icon>
            <span class="font-display font-weight-bold text-subtitle-1">System- &amp; Tarifparameter</span>
          </div>
          <v-btn icon="mdi-close" variant="text" size="small" @click="settingsDialog = false" />
        </v-card-title>

        <v-card-text class="px-4 pt-2 pb-3">
          <p class="text-caption text-secondary-adaptive mb-4">
            Diese Parameter werden lokal gespeichert und für die Echtzeit-Berechnung von Ersparnis,
            Einspeisevergütung und Anlagen-Auslastung genutzt.
          </p>

          <v-row dense>
            <v-col cols="6" class="mb-2">
              <label class="text-caption font-weight-semibold text-secondary-adaptive d-block mb-1">
                Netzbezugspreis (ct/kWh)
              </label>
              <v-text-field
                v-model.number="settings.electricityPriceCent"
                type="number"
                step="0.5"
                density="compact"
                hide-details
                variant="solo-filled"
              />
            </v-col>

            <v-col cols="6" class="mb-2">
              <label class="text-caption font-weight-semibold text-secondary-adaptive d-block mb-1">
                Einspeisevergütung (ct/kWh)
              </label>
              <v-text-field
                v-model.number="settings.feedInTariffCent"
                type="number"
                step="0.1"
                density="compact"
                hide-details
                variant="solo-filled"
              />
            </v-col>

            <v-col cols="6" class="mb-2">
              <label class="text-caption font-weight-semibold text-secondary-adaptive d-block mb-1">
                PV-Nennleistung (kWp)
              </label>
              <v-text-field
                v-model.number="settings.systemCapacityKwp"
                type="number"
                step="0.5"
                density="compact"
                hide-details
                variant="solo-filled"
              />
            </v-col>

            <v-col cols="6" class="mb-2">
              <label class="text-caption font-weight-semibold text-secondary-adaptive d-block mb-1">
                Speicherkapazität (kWh)
              </label>
              <v-text-field
                v-model.number="settings.batteryCapacityKwh"
                type="number"
                step="0.5"
                density="compact"
                hide-details
                variant="solo-filled"
              />
            </v-col>
          </v-row>
        </v-card-text>

        <v-card-actions class="px-4 pb-3 pt-0 d-flex justify-space-between">
          <v-btn
            variant="text"
            size="small"
            color="secondary"
            class="text-none"
            @click="resetSettings"
          >
            Standardwerte
          </v-btn>
          <v-btn
            variant="flat"
            color="emerald"
            size="small"
            class="text-none font-weight-bold px-4 rounded-lg"
            @click="settingsDialog = false"
          >
            Fertig
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app-bar>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useTheme } from 'vuetify'
import { useSettings } from '@/composables/useSettings'

const { settings, activeSection, isOnline, toggleTheme, resetSettings } = useSettings()
const vuetifyTheme = useTheme()
const settingsDialog = ref(false)

// Sync initial Vuetify theme with persisted setting
vuetifyTheme.global.name.value = settings.value.theme

watch(
  () => settings.value.theme,
  (newTheme) => {
    vuetifyTheme.global.name.value = newTheme
  }
)

const handleToggleTheme = () => {
  toggleTheme()
}

const navItems = [
  { id: 'overview', label: 'Dashboard', icon: 'mdi-view-dashboard-outline' },
  { id: 'flow', label: 'Energiefluss', icon: 'mdi-transit-connection-variant' },
  { id: 'forecast', label: '48h Prognose', icon: 'mdi-weather-partly-sunny' },
  { id: 'analytics', label: 'Historie & Analyse', icon: 'mdi-chart-box-outline' },
]
</script>

<style scoped>
.solar-app-bar {
  background: var(--bg-header) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  border-bottom: 1px solid var(--border-subtle) !important;
}

.brand-badge {
  width: 38px;
  height: 38px;
  border-radius: 11px;
  background: var(--color-pv-soft);
  border: 1px solid var(--color-pv-border);
  box-shadow: 0 0 20px rgba(16, 185, 129, 0.15);
}

.version-chip {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 6px;
  background: var(--bg-pill);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
}

.action-icon-btn {
  border: 1px solid var(--border-subtle);
  background: var(--bg-subtle);
  border-radius: 10px !important;
  color: var(--text-secondary) !important;
  width: 36px !important;
  height: 36px !important;
}

.action-icon-btn:hover {
  border-color: var(--border-strong);
  color: var(--text-primary) !important;
  background: var(--bg-subtle-hover);
}

.offline-dot {
  background: var(--color-grid-out) !important;
}
.offline-dot::after {
  background: var(--color-grid-out) !important;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
