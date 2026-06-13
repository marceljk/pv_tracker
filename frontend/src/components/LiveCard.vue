<template>
  <v-card class="live-card overflow-hidden h-100">
    <template v-slot:title>
      <div class="d-flex align-center justify-space-between w-100">
        <div class="d-flex align-center">
          <v-icon color="emerald" class="mr-2 pulse-slow">mdi-pulse</v-icon>
          <span class="font-weight-bold">Live-Status</span>
        </div>
        <v-chip size="x-small" color="emerald" variant="tonal" class="font-weight-medium">Echtzeit</v-chip>
      </div>
    </template>

    <template v-slot:text>
      <div v-if="liveData" class="d-flex flex-column justify-center h-100 pt-2">
        <v-row class="align-center">
          <!-- Left: 4 Metric Cards -->
          <v-col cols="12" md="8" class="order-last order-md-first">
            <v-row dense>
              <!-- PV Power -->
              <v-col cols="6" class="pa-1">
                <div class="live-node pv d-flex flex-column align-start">
                  <div class="d-flex align-center w-100 mb-1">
                    <v-icon color="emerald" class="mr-1" size="18">mdi-solar-power-variant</v-icon>
                    <span class="text-caption text-secondary text-truncate">PV Erzeugung</span>
                  </div>
                  <span class="stat-value text-h6 text-emerald">{{ formatVal(liveData.pvPower) }} W</span>
                </div>
              </v-col>

              <!-- Stromverbrauch -->
              <v-col cols="6" class="pa-1">
                <div class="live-node consumption d-flex flex-column align-start">
                  <div class="d-flex align-center w-100 mb-1">
                    <v-icon color="blue" class="mr-1" size="18">mdi-home-lightning-bolt</v-icon>
                    <span class="text-caption text-secondary text-truncate">Stromverbrauch</span>
                  </div>
                  <span class="stat-value text-h6 text-blue">{{ formatVal(liveData.powerConsumption) }} W</span>
                </div>
              </v-col>

              <!-- Netzbezug / Netzeinspeisung -->
              <v-col cols="6" class="pa-1">
                <div :class="['live-node d-flex flex-column align-start', liveData.gridPower < 0 ? 'grid-out' : 'grid']">
                  <div class="d-flex align-center w-100 mb-1">
                    <v-icon :color="liveData.gridPower < 0 ? 'rose' : 'cyan'" class="mr-1" size="18">mdi-transmission-tower</v-icon>
                    <span class="text-caption text-secondary text-truncate">
                      {{ liveData.gridPower < 0 ? 'Netzbezug' : 'Netzeinspeisung' }}
                    </span>
                  </div>
                  <span :class="['stat-value text-h6', liveData.gridPower < 0 ? 'text-rose' : 'text-cyan']">
                    {{ formatVal(Math.abs(liveData.gridPower)) }} W
                  </span>
                </div>
              </v-col>

              <!-- Akku laden / entladen -->
              <v-col cols="6" class="pa-1">
                <div class="live-node battery d-flex flex-column align-start">
                  <div class="d-flex align-center w-100 mb-1">
                    <v-icon color="amber" class="mr-1" size="18">
                      {{ liveData.batteryPower < 0 ? 'mdi-battery-arrow-down' : 'mdi-battery-arrow-up' }}
                    </v-icon>
                    <span class="text-caption text-secondary text-truncate">
                      {{ liveData.batteryPower < 0 ? 'Akku entlädt' : 'Akku lädt' }}
                    </span>
                  </div>
                  <span class="stat-value text-h6 text-amber">{{ formatVal(Math.abs(liveData.batteryPower)) }} W</span>
                </div>
              </v-col>
            </v-row>
          </v-col>

          <!-- Right: Circular Gauge -->
          <v-col cols="12" md="4" class="order-first order-md-last d-flex flex-column align-center justify-center py-4">
            <v-progress-circular 
              :model-value="liveData.batteryPercent || 0" 
              :size="batterySize" 
              :width="batteryWidth"
              color="amber"
              class="pulse-battery"
            >
              <div class="d-flex flex-column align-center">
                <span :class="[xs ? 'text-subtitle-1' : 'text-h5', 'font-weight-bold text-amber stat-value']">{{ liveData.batteryPercent || 0 }}%</span>
                <span class="text-caption text-secondary font-weight-bold" style="font-size: 0.65rem !important;">Akku</span>
              </div>
            </v-progress-circular>
          </v-col>
        </v-row>
      </div>
      <div v-else class="d-flex justify-center align-center h-100 py-8">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
      </div>
    </template>
  </v-card>
</template>

<script>
import { defineComponent, computed } from 'vue';
import { useDatabaseObject, useDatabase } from "vuefire";
import { ref as dbRef, query, orderByChild, limitToLast } from 'firebase/database'
import { useDisplay } from 'vuetify';

export default defineComponent({
  setup() {
    const db = useDatabase();
    const liveData = useDatabaseObject(query(dbRef(db, "live"), orderByChild("batteryPercent"), limitToLast(10)));
    const { xs } = useDisplay();
    
    const batterySize = computed(() => xs.value ? 75 : 105);
    const batteryWidth = computed(() => xs.value ? 7 : 10);

    const formatVal = (val) => {
      if (val === undefined || val === null || isNaN(val)) return '0';
      return Math.round(val).toLocaleString('de-DE');
    };

    return {
      liveData,
      formatVal,
      batterySize,
      batteryWidth,
      xs,
    }
  }
});
</script>