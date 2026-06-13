<template>
  <v-container fluid class="pa-4 pa-sm-6">
    <!-- Dashboard Page Header -->
    <div class="mb-4 mb-sm-6 d-flex justify-space-between align-center flex-wrap gap-4">
      <div>
        <h1 class="text-h4 font-weight-bold mb-1" style="letter-spacing: -0.04em;">Energie-Zentrale</h1>
        <p class="text-subtitle-2 text-secondary d-none d-sm-block">Echtzeit-Überwachung und PV-Ertragsprognosen</p>
      </div>
      <v-chip variant="flat" color="rgba(255,255,255,0.05)" class="py-4 border px-3 rounded-lg">
        <v-icon start color="primary" class="mr-1">mdi-calendar</v-icon>
        <span class="font-weight-medium text-caption text-secondary">{{ formattedToday }}</span>
      </v-chip>
    </div>

    <!-- Main Grid -->
    <v-row>
      <!-- Live Metrics Visualization (Prominent) -->
      <v-col cols="12" lg="6">
        <LiveCard class="h-100" />
      </v-col>

      <!-- Forecast Card -->
      <v-col cols="12" lg="6">
        <ForecastCard class="h-100" />
      </v-col>
    </v-row>

    <!-- Desktop view: 4 historical cards side-by-side -->
    <v-row v-if="mdAndUp" class="mt-4">
      <v-col cols="12" md="3">
        <RangeSumCard title="Heute" :start="today" :end="today" class="h-100" />
      </v-col>
      <v-col cols="12" md="3">
        <RangeSumCard title="Letzte Woche" :start="oneWeekAgo" :end="today" class="h-100" />
      </v-col>
      <v-col cols="12" md="3">
        <RangeSumCard title="Aktuelles Jahr" :start="startOfYear" :end="today" class="h-100" />
      </v-col>
      <v-col cols="12" md="3">
        <CustomRangeSumCard title="Eigene Zeiträume" class="h-100" />
      </v-col>
    </v-row>

    <!-- Mobile view: Single tabbed card with flat embedded components -->
    <v-card v-else class="mt-4 overflow-hidden">
      <v-tabs v-model="activeTab" grow color="primary" class="border-b" density="compact">
        <v-tab value="heute">Heute</v-tab>
        <v-tab value="woche">Woche</v-tab>
        <v-tab value="jahr">Jahr</v-tab>
        <v-tab value="custom">Eigene</v-tab>
      </v-tabs>
      
      <v-window v-model="activeTab">
        <v-window-item value="heute" class="pa-2">
          <RangeSumCard title="Heute" :start="today" :end="today" flat />
        </v-window-item>
        <v-window-item value="woche" class="pa-2">
          <RangeSumCard title="Letzte Woche" :start="oneWeekAgo" :end="today" flat />
        </v-window-item>
        <v-window-item value="jahr" class="pa-2">
          <RangeSumCard title="Aktuelles Jahr" :start="startOfYear" :end="today" flat />
        </v-window-item>
        <v-window-item value="custom" class="pa-2">
          <CustomRangeSumCard title="Eigene Zeiträume" flat />
        </v-window-item>
      </v-window>
    </v-card>
  </v-container>
</template>

<script>
import { defineComponent, computed, ref } from 'vue';
import { useDisplay } from 'vuetify';
import ForecastCard from '../components/ForecastCard.vue';
import LiveCard from '../components/LiveCard.vue';
import RangeSumCard from '../components/RangeSumCard.vue';
import CustomRangeSumCard from '../components/CustomRangeSumCard.vue';

export default defineComponent({
  components: { ForecastCard, LiveCard, RangeSumCard, CustomRangeSumCard },
  setup() {
    const { mdAndUp } = useDisplay();
    const activeTab = ref('heute');
    
    const today = computed(() => new Date().toISOString().substring(0, 11));
    const oneWeekAgo = computed(() => new Date(new Date().setDate(new Date().getDate() - 7)).toISOString().substring(0, 11));
    const startOfYear = computed(() => new Date(new Date().setMonth(0, 1)).toISOString().substring(0, 11));
    const formattedToday = computed(() => {
      return new Date().toLocaleDateString('de-DE', { 
        weekday: 'long', 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric' 
      });
    });
    
    return {
      today,
      oneWeekAgo,
      startOfYear,
      formattedToday,
      mdAndUp,
      activeTab,
    }
  },
});
</script>