<template>
  <v-card :title="title" :loading="dataFromDb.pending">
    <template v-slot:text>
      <v-row>
        <v-col cols="6">
          <v-text-field hide-details="auto" v-model="startDate" label="Startdatum" type="date"></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field hide-details="auto" v-model="endDate" label="Enddatum" type="date"></v-text-field>
        </v-col>
      </v-row>
      <v-row v-if="dataFromDb">
        <v-col cols="12">
          <v-row dense v-for="val in Object.entries(dataSum)" :key="val[0]">
            <v-col cols="8">{{ liveText(val[0], val[1]) }}</v-col>
            <v-col cols="4" class="text-right">{{ valueText(val[1]) }} </v-col>
          </v-row>
        </v-col>
      </v-row>
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
  },
  setup(props) {
    const startDate = ref("");
    const endDate = ref("");

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
    })

    const db = useDatabase();
    const dataFromDb = useDatabaseList(
      queryDb,
    );

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

    const liveText = (key, value) => {
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

    return {
      dataFromDb,
      dataSum,
      liveText,
      valueText,
      startDate,
      endDate,
    };
  },
});
</script>
