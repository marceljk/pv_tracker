<template>
  <v-card
    :class="[flat ? 'elevation-0 bg-transparent border-0' : 'summary-card custom-range-card overflow-hidden h-100']">
    <template v-slot:title v-if="!flat">
      <div class="d-flex align-center justify-space-between w-100">
        <div class="d-flex align-center">
          <v-icon color="secondary" class="mr-2" size="20">mdi-calendar-edit</v-icon>
          <span class="font-weight-bold">{{ title }}</span>
        </div>
      </div>
    </template>

    <template v-slot:loader v-if="isLoading">
      <v-progress-linear color="primary" height="2" indeterminate></v-progress-linear>
    </template>

    <template v-slot:text>
      <v-row class="mb-3" dense>
        <v-col cols="6" class="pr-1">
          <v-text-field hide-details="auto" v-model="startDate" label="Start" type="date" variant="solo-filled"
            density="compact" class="rounded-lg date-input"></v-text-field>
        </v-col>
        <v-col cols="6" class="pl-1">
          <v-text-field hide-details="auto" v-model="endDate" label="Ende" type="date" variant="solo-filled"
            density="compact" class="rounded-lg date-input"></v-text-field>
        </v-col>
      </v-row>

      <div v-if="hasRange" class="py-1">
        <div v-for="val in Object.entries(dataSum)" :key="val[0]"
          class="d-flex align-center justify-space-between py-2 border-b-dimmed">
          <div class="d-flex align-center">
            <v-avatar size="28" :color="getMetricBgColor(val[0])" class="mr-3" rounded="lg">
              <v-icon :color="getMetricColor(val[0])" size="16">{{ getMetricIcon(val[0]) }}</v-icon>
            </v-avatar>
            <span class="text-caption text-secondary font-weight-medium">{{ liveText(val[0]) }}</span>
          </div>
          <span :class="['font-weight-bold stat-value text-body-1', getMetricTextClass(val[0])]">
            {{ valueText(val[1]) }}
          </span>
        </div>
      </div>
      <div v-else class="d-flex flex-column align-center justify-center py-6 text-muted">
        <v-icon size="28" class="mb-2 opacity-40">mdi-calendar-alert</v-icon>
        <span class="text-caption font-weight-medium">Zeitraum auswählen</span>
      </div>
    </template>
  </v-card>
</template>

<script>
import { defineComponent, ref, computed, watch } from "vue";
import { useDatabase, useDatabaseList } from "vuefire";
import {
  ref as dbRef,
  query,
  orderByKey,
  startAt,
  endAt,
} from "firebase/database";

export default defineComponent({
  props: {
    title: {
      type: String,
    },
    flat: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const startDate = ref("");
    const endDate = ref("");
    const db = useDatabase();
    const isLoading = ref(false);

    const hasRange = computed(() => !!(startDate.value && endDate.value));

    const queryDb = computed(() => {
      if (startDate.value && endDate.value) {
        return query(
          dbRef(db, "dailySum"),
          orderByKey(),
          startAt(startDate.value),
          endAt(endDate.value)
        );
      }
      return null;
    });

    const dataFromDb = useDatabaseList(queryDb);

    watch(queryDb, (newQuery) => {
      if (newQuery) {
        isLoading.value = true;
      } else {
        isLoading.value = false;
      }
    });

    watch(dataFromDb, () => {
      isLoading.value = false;
    });

    const dataSum = computed(() => {
      let sum = {
        gridPowerIn: 0,
        gridPowerOut: 0,
        powerConsumption: 0,
        pvPower: 0,
      };
      if (dataFromDb.value) {
        dataFromDb.value.forEach((day) => {
          const measureTime = day.count / 120;
          sum = {
            gridPowerIn:
              sum.gridPowerIn +
              (day.gridPowerIn / day.count) * measureTime,
            gridPowerOut:
              sum.gridPowerOut +
              (day.gridPowerOut / day.count) * measureTime,
            powerConsumption:
              sum.powerConsumption +
              (day.powerConsumption / day.count) * measureTime,
            pvPower: sum.pvPower + (day.pvPower / day.count) * measureTime,
          };
        });
      }
      return sum;
    });

    const liveText = (key) => {
      const x = {
        gridPowerIn: "Netzeinspeisung",
        gridPowerOut: "Netzbezug",
        powerConsumption: "Stromverbrauch",
        pvPower: "PV Erzeugung",
      };
      return x[key];
    };

    const valueText = (value) => {
      if (value < 10000 && value > -10000) {
        return `${Math.round(Math.abs(value)).toLocaleString()} W`;
      }
      const kWh = Math.round(value / 100) / 10; // one digit
      return `${Math.abs(kWh).toLocaleString()} kW`;
    };

    const getMetricIcon = (key) => {
      const icons = {
        pvPower: "mdi-solar-power-variant",
        powerConsumption: "mdi-home-lightning-bolt",
        gridPowerIn: "mdi-transmission-tower-export",
        gridPowerOut: "mdi-transmission-tower-import",
      };
      return icons[key] || "mdi-flash";
    };

    const getMetricColor = (key) => {
      const colors = {
        pvPower: "emerald",
        powerConsumption: "blue",
        gridPowerIn: "cyan",
        gridPowerOut: "rose",
      };
      return colors[key] || "primary";
    };

    const getMetricBgColor = (key) => {
      const colors = {
        pvPower: "rgba(16, 185, 129, 0.08)",
        powerConsumption: "rgba(59, 130, 246, 0.08)",
        gridPowerIn: "rgba(6, 182, 212, 0.08)",
        gridPowerOut: "rgba(244, 63, 94, 0.08)",
      };
      return colors[key] || "rgba(255,255,255,0.05)";
    };

    const getMetricTextClass = (key) => {
      const classes = {
        pvPower: "text-emerald",
        powerConsumption: "text-blue",
        gridPowerIn: "text-cyan",
        gridPowerOut: "text-rose",
      };
      return classes[key] || "";
    };

    return {
      dataFromDb,
      dataSum,
      liveText,
      valueText,
      startDate,
      endDate,
      hasRange,
      getMetricIcon,
      getMetricColor,
      getMetricBgColor,
      getMetricTextClass,
      isLoading,
    };
  },
});
</script>

<style scoped>
.border-b-dimmed {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.border-b-dimmed:last-child {
  border-bottom: none;
}

.date-input :deep(.v-field__input) {
  padding-top: 6px !important;
  padding-bottom: 6px !important;
  font-size: 0.85rem;
}
</style>
