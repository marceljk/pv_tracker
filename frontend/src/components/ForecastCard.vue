<template>
  <v-card title="Vorhersage">
    <template v-slot:text>
      <v-row>
        <v-col>
          <table v-if="dailyForecast" width="100%">
            <thead>
              <tr>
                <th class="text-left">Tag</th>
                <th class="text-right">kWh</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="day in dailyForecast" :key="day.id">
                <td>{{ day.day }}</td>
                <td class="text-right">{{ day.estimate }}</td>
              </tr>
            </tbody>
          </table>
        </v-col>
      </v-row>
    </template>
  </v-card>
</template>

<script>
import { defineComponent } from 'vue';
import { useDatabase, useDatabaseList } from "vuefire";
import { ref as dbRef, query, orderByKey } from 'firebase/database'

export default defineComponent({
  setup() {
    const db = useDatabase();
    const dailyForecast = useDatabaseList(query(dbRef(db, "dailyForecast"), orderByKey()));

    return {
      dailyForecast,
    }
  }
});
</script>