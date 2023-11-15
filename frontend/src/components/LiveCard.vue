<template>
  <v-card>
    <template v-slot:title>
      Live
    </template>
    <template v-slot:text>
      <v-row v-if="liveData">
        <v-col cols="9">
          <template v-for="(key) in Object.keys(liveData)" :key="key">
            <v-row dense v-if="key !== 'batteryPercent'">
              <v-col cols="8">{{ liveText(key, liveData[key]) }}</v-col>
              <v-col class="text-right" cols="4">{{ Math.abs(liveData[key]) }} W</v-col>
            </v-row>
          </template>
        </v-col>
        <v-col cols="3" class="d-inline-flex align-center">
          <v-progress-circular :model-value="liveData.batteryPercent" :size="75" :width="10">
            {{ liveData.batteryPercent }}
          </v-progress-circular>
        </v-col>
      </v-row>
    </template>
  </v-card>
</template>

<script>
import { defineComponent } from 'vue';
import { useDatabaseObject, useDatabase } from "vuefire";
import { ref as dbRef, query, orderByChild, limitToLast } from 'firebase/database'

export default defineComponent({
  setup() {
    const db = useDatabase();
    const liveData = useDatabaseObject(query(dbRef(db, "live"), orderByChild("batteryPercent"), limitToLast(10)));
    const liveText = (key, value) => {
      const x = {
        batteryPower: value < 0 ? "Akku entlädt" : "Akku lädt",
        gridPower: value < 0 ? "Netzbezug" : "Netzeinspeisung",
        powerConsumption: "Stromverbrauch",
        pvPower: "PV Erzeugung",
      };
      return x[key];
    };

    return {
      liveData,
      liveText,
    }
  }
});
</script>