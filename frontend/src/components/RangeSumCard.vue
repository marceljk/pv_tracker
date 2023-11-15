<template>
  <v-card :title="title">
    <template v-slot:text>
      <v-row v-if="dailySumList">
        <v-col cols="12">
          <v-row dense v-for="val in Object.entries(todaySum)" :key="val[0]">
            <v-col cols="8">{{ liveText(val[0], val[1]) }}</v-col>
            <v-col cols="4" class="text-right">{{ Math.round(Math.abs(val[1])).toLocaleString() }} W</v-col>
          </v-row>
        </v-col>
      </v-row>
    </template>
  </v-card>
</template>

<script>
import { defineComponent, computed } from "vue";
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
    start: {
      type: String,
    },
    end: {
      type: String,
    },
    title: {
      type: String,
    }
  },
  setup(props) {
    const db = useDatabase();
    const dailySumList = useDatabaseList(
      query(dbRef(db, "dailySum"), orderByKey(), startAt(props.start), endAt(props.end))
    );

    const todaySum = computed(() => {
      let sum = {
        batteryPower: 0,
        gridPowerIn: 0,
        gridPowerOut: 0,
        powerConsumption: 0,
        pvPower: 0,
      };
      dailySumList.value.forEach((day) => {
        const measureTime = day.count / 120;
        sum = {
          batteryPower: sum.batteryPower + (day.batteryPower / day.count) * measureTime,
          gridPowerIn: sum.gridPowerIn + (day.gridPowerIn / day.count) * measureTime,
          gridPowerOut: sum.gridPowerOut + (day.gridPowerOut / day.count) * measureTime,
          powerConsumption: sum.powerConsumption + (day.powerConsumption / day.count) * measureTime,
          pvPower: sum.pvPower + (day.pvPower / day.count) * measureTime,
        }
      });
      if (!sum) return {};
      return sum;
    });

    const liveText = (key, value) => {
      const x = {
        batteryPower: value < 0 ? "Akku entladen" : "Akku geladen",
        gridPowerIn: "Netzeinspeisung",
        gridPowerOut: "Netzbezug",
        powerConsumption: "Stromverbrauch",
        pvPower: "PV Erzeugung",
      };
      return x[key];
    };

    return {
      dailySumList,
      todaySum,
      liveText,
    };
  },
});
</script>