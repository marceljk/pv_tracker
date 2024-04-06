<template>
  <v-card title="Vorhersage">
    <template v-slot:text>
      <v-row>
        <v-col>
          <v-data-table
            :expanded="expandedItems"
            density="compact"
            :items="dailyForecast"
            :headers="headers"
            show-expand
            @update:expanded="expandSingleItem"
            @click:row="clickRow"
          >
            <template v-slot:expanded-row="{ columns, item }">
              <tr>
                <td :colspan="columns.length" :key="item.id">
                  <Line
                    :id="item.id"
                    :options="chartOptions"
                    :data="chartData(item.id)"
                  />
                </td>
              </tr>
            </template>
            <template v-slot:bottom></template>
          </v-data-table>
        </v-col>
      </v-row>
    </template>
  </v-card>
</template>

<script setup>
import { ref } from "vue";
import { useDatabase, useDatabaseList } from "vuefire";
import { ref as dbRef, query, orderByKey } from "firebase/database";
import { Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
} from "chart.js";
ChartJS.register(
  Title,
  Tooltip,
  Legend,
  PointElement,
  LineElement,
  CategoryScale,
  LinearScale
);

const expandedItems = ref([]);
const headers = ref([
  {
    title: "Tag",
    value: "day",
  },
  {
    title: "kWh",
    value: "estimate",
  },
]);

const chartOptions = ref({
  maintainAspectRatio: true,
});

const db = useDatabase();
const dailyForecast = useDatabaseList(
  query(dbRef(db, "dailyForecast"), orderByKey())
);

const hourlyForecast = useDatabaseList(
  query(dbRef(db, "hourlyForecast"), orderByKey())
);

const expandSingleItem = (items) => {
  const filteredItems = items.filter(
    (item) => !expandedItems.value.includes(item)
  );
  expandedItems.value = filteredItems;
};

const chartData = (day) => {
  const date = new Date(Number(day)).toLocaleDateString();

  const dataset = () => {
    const dataOfCurrentDay = hourlyForecast.value.filter(
      (item) => new Date(item.period_end).toLocaleDateString() == date
    );
    const labels = [];
    const data = [];
    dataOfCurrentDay.forEach((item) => {
      if (item.pv_estimate == 0) return;
      labels.push(
        new Date(item.period_end).toLocaleTimeString("de").substring(0, 5)
      );
      data.push(item.pv_estimate);
    });

    return {
      labels,
      data,
    };
  };

  const data = dataset();

  return {
    labels: data.labels,
    datasets: [
      {
        data: data.data,
        label: date,
        backgroundColor: "rgb(0, 255, 132)",
        borderColor: "rgb(0, 255, 132)",
        tension: 0.5,
        radius: 0,
      },
    ],
  };
};

const clickRow = (event, item) => {
  const { id } = item.item;
  if (expandedItems.value.includes(id)) expandedItems.value = [];
  else expandedItems.value = [id];
}
</script>
