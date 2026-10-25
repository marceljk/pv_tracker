<template>
  <v-app>
    <default-bar />

    <default-view />

    <!-- Mobile PWA Glass Bottom Navigation Dock -->
    <nav class="mobile-bottom-dock d-flex d-md-none" aria-label="Mobile Navigation">
      <button
        v-for="item in navItems"
        :key="item.id"
        type="button"
        :class="['dock-item', { active: activeSection === item.id }]"
        @click="activeSection = item.id"
      >
        <v-icon size="20" class="dock-icon">{{ item.icon }}</v-icon>
        <span class="dock-label">{{ item.label }}</span>
      </button>
    </nav>
  </v-app>
</template>

<script setup>
import DefaultBar from './AppBar.vue'
import DefaultView from './View.vue'
import { useSettings } from '@/composables/useSettings'

const { activeSection } = useSettings()

const navItems = [
  { id: 'overview', label: 'Übersicht', icon: 'mdi-view-dashboard-outline' },
  { id: 'flow', label: 'Live-Fluss', icon: 'mdi-transit-connection-variant' },
  { id: 'forecast', label: 'Prognose', icon: 'mdi-weather-partly-sunny' },
  { id: 'analytics', label: 'Analyse', icon: 'mdi-chart-box-outline' },
]
</script>

<style scoped>
.mobile-bottom-dock {
  position: fixed;
  bottom: max(12px, env(safe-area-inset-bottom));
  left: 12px;
  right: 12px;
  z-index: 1005;
  background: var(--bg-header);
  backdrop-filter: blur(28px) saturate(190%);
  -webkit-backdrop-filter: blur(28px) saturate(190%);
  border: 1px solid var(--border-strong);
  border-radius: 18px;
  padding: 6px;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.45);
  justify-content: space-around;
  align-items: center;
}

.dock-item {
  appearance: none;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 3px;
  padding: 6px 12px;
  border-radius: 12px;
  flex: 1;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.dock-item.active {
  background: var(--color-pv-soft);
  color: var(--color-pv);
}

.dock-label {
  font-size: 0.68rem;
  font-weight: 600;
  letter-spacing: -0.01em;
}
</style>
