<template>
  <v-card class="forecast-card h-100">
    <template v-slot:title>
      <div class="d-flex align-center justify-space-between w-100">
        <div class="d-flex align-center">
          <v-icon color="emerald" class="mr-2">mdi-chart-timeline-variant</v-icon>
          <span class="font-weight-bold">Ertragsprognose</span>
        </div>
        <v-chip size="x-small" color="emerald" variant="tonal" class="font-weight-medium">48h Ausblick</v-chip>
      </div>
    </template>
    
    <template v-slot:text>
      <v-row>
        <v-col class="pa-0">
          <v-data-table
            :expanded="expandedItems"
            density="compact"
            :items="dailyForecast"
            :headers="headers"
            show-expand
            @update:expanded="expandSingleItem"
            @click:row="clickRow"
            class="elevation-0"
          >
            <!-- Custom Row Cells -->
            <template #[`item.day`]="{ item }">
              <span class="font-weight-medium text-slate">{{ formatDay(item.day) }}</span>
            </template>
            
            <template #[`item.estimate`]="{ item }">
              <v-chip size="small" color="emerald" variant="tonal" class="font-weight-bold px-2">
                {{ formatEstimate(item.estimate) }} kWh
              </v-chip>
            </template>

            <!-- Expanded Line Chart -->
            <template v-slot:expanded-row="{ columns, item }">
              <tr>
                <td :colspan="columns.length" :key="item.id" class="pa-3" style="background: rgba(10, 13, 22, 0.45);">
                  <div style="height: 200px;" class="w-100">
                    <Line
                      :id="item.id"
                      :options="chartOptions"
                      :data="chartData(item.id)"
                    />
                  </div>
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
  Filler,
} from "chart.js";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  PointElement,
  LineElement,
  CategoryScale,
  LinearScale,
  Filler
);

const expandedItems = ref([]);
const headers = ref([
  {
    title: "Tag",
    value: "day",
    align: "start",
    sortable: false,
  },
  {
    title: "Erwarteter Ertrag",
    value: "estimate",
    align: "end",
    sortable: false,
  },
]);

const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      backgroundColor: "rgba(15, 23, 42, 0.95)",
      titleColor: "#f8fafc",
      bodyColor: "#94a3b8",
      borderColor: "rgba(255, 255, 255, 0.08)",
      borderWidth: 1,
      padding: 10,
      cornerRadius: 8,
      bodyFont: {
        family: "Inter",
        size: 11,
      },
      titleFont: {
        family: "Outfit",
        weight: "bold",
        size: 12,
      },
      displayColors: false,
    },
  },
  scales: {
    x: {
      grid: {
        color: "rgba(255, 255, 255, 0.04)",
      },
      ticks: {
        color: "#94a3b8",
        font: {
          family: "Inter",
          size: 10,
        },
      },
    },
    y: {
      grid: {
        color: "rgba(255, 255, 255, 0.04)",
      },
      ticks: {
        color: "#94a3b8",
        font: {
          family: "Inter",
          size: 10,
        },
      },
    },
  },
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

const formatDay = (val) => {
  if (!val) return "";
  const num = Number(val);
  if (!isNaN(num) && num > 1000000000000) {
    return new Date(num).toLocaleDateString("de-DE", {
      weekday: "long",
      day: "2-digit",
      month: "2-digit",
    });
  }
  return val;
};

const formatEstimate = (val) => {
  if (val === undefined || val === null || isNaN(val)) return "0";
  return Number(val).toLocaleString("de-DE", { minimumFractionDigits: 1, maximumFractionDigits: 1 });
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
        label: "Prognose (kW)",
        backgroundColor: "rgba(16, 185, 129, 0.12)",
        borderColor: "#10b981",
        borderWidth: 2,
        tension: 0.4,
        fill: true,
        pointBackgroundColor: "#10b981",
        pointBorderColor: "#fff",
        pointBorderWidth: 1.5,
        pointRadius: 1.5,
        pointHoverRadius: 4,
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
