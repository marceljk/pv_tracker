<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" md="3">
        <LiveCard />
      </v-col>
      <v-col cols="12" md="3">
        <ForecastCard />
      </v-col>
      <v-col cols="12" md="3">
        <RangeSumCard title="Heute" :start="today" :end="today" />
      </v-col>
      <v-col cols="12" md="3">
        <RangeSumCard title="Letzte Woche" :start="oneWeekAgo" :end="today" />
      </v-col>
      <v-col cols="12" md="3">
        <RangeSumCard title="Aktuelles Jahr" :start="startOfYear" :end="today" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { defineComponent, computed } from 'vue';
import ForecastCard from '../components/ForecastCard.vue';
import LiveCard from '../components/LiveCard.vue';
import RangeSumCard from '../components/RangeSumCard.vue';

export default defineComponent({
  components: { ForecastCard, LiveCard, RangeSumCard },
  setup() {
    const today = computed(() => new Date().toISOString().substring(0, 11));
    const oneWeekAgo = computed(() => new Date(new Date().setDate(new Date().getDate() - 7)).toISOString().substring(0, 11));
    const startOfYear = computed(() => new Date(new Date().setMonth(1, 1)).toISOString().substring(0, 11));
    return {
      today,
      oneWeekAgo,
      startOfYear,
    }
  },
});
</script>