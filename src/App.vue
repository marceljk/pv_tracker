<template>
  <v-app>
    <!-- <v-app-bar app> </v-app-bar> -->

    <v-main>
      <v-container>
        <v-row dense>
          <v-col class="ma-2" align="center">
            <v-container class="text-h4 mb-2"> Photovoltaik </v-container>
            <v-card class="ma-2 px-6 py-4">
              <v-row>
                <v-col>
                  <v-progress-circular
                    :color="percentColor"
                    width="10"
                    rotate="-90"
                    size="100"
                    :value="power.batteryPercent"
                    :indeterminate="power.batteryPercent == null"
                  >
                    <div class="text-h4" v-if="power.batteryPercent !== null">
                      {{ power.batteryPercent + "%" }}
                    </div>
                  </v-progress-circular>
                </v-col>

                <v-row v-if="power.batteryPercent != null">
                  <v-col class="text-left" align-self="center">
                    <StatsText
                      :value="power.batteryPower"
                      negResponse="Akku entlädt"
                      posResponse="Akku lädt"
                    />
                    <StatsText
                      :value="power.gridPower"
                      negResponse="Netzbezug"
                      posResponse="Netzeinspeisung"
                    />
                    <div>PV Erzeugung</div>
                    <div>Stromverbrauch</div>
                  </v-col>

                  <v-col class="text-right" align-self="center">
                    <div v-if="power.batteryPower">
                      {{ Math.abs(power.batteryPower) }} W
                    </div>
                    <div>{{ Math.abs(power.gridPower) }} W</div>
                    <div>{{ Math.abs(power.pvPower) }} W</div>
                    <div>{{ Math.abs(power.powerConsumption) }} W</div>
                  </v-col>
                </v-row>
              </v-row>
            </v-card>
            <br />
            <v-card class="ma-2 pa-2">
              <div class="text-h6">Prognose</div>
              <v-sparkline
                class="ma-2"
                padding="10"
                :value="pvForecast.pv_estimate"
                :labels="pvForecast.period_end"
                smooth="10"
                height="150"
                auto-draw
                stroke-linecap="round"
              >
                <template v-slot:label="item">
                  <template v-if="item.index % 2 == 0"> 
                    {{ item.value }}
                  </template>
                </template>
              </v-sparkline>
              <v-data-table
                class="ma-2"
                dense
                :loading="dailyForecast.length == 0"
                mobile-breakpoint="0"
                hide-default-footer
                :headers="headers"
                :items="dailyForecast"
              >
              </v-data-table>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import StatsText from "./components/StatsText.vue";

export default {
  name: "App",

  components: {
    StatsText,
  },

  data() {
    return {
      power: {
        batteryPercent: null,
        gridPower: null,
        pvPower: null,
        powerConsumption: null,
        batteryPower: null,
      },
      interval: 0,
      pvForecast: {
        pv_estimate: [],
        period_end: [],
        period: "",
      },
      dailyForecast: [],
      headers: [
        {
          text: "Tag",
          value: "day",
          sortable: false,
        },
        { text: "kWh", value: "estimate", align: "end" },
      ],
    };
    //
  },

  computed: {
    percentColor() {
      if (this.power.batteryPercent > 65) {
        return "green";
      } else if (this.power.batteryPercent > 35) {
        return "orange";
      }
      return "red";
    },
  },

  created() {
    this.fetchData();
    setInterval(this.fetchData, 1000);
    this.fetchForecast();
    this.fetchDailyForecast();
    setInterval(() => {
      this.fetchForecast();
      this.fetchDailyForecast();
    }, 1000 * 60);
  },

  methods: {
    async fetchData() {
      let res = await fetch("api/data");
      this.power = await res.json();
    },
    async fetchForecast() {
      let res = await fetch("api/forecast");
      const data = await res.json();

      this.pvForecast = {
        pv_estimate: data.pv_estimate,
        period_end: data.period_end.map((x) => {
          const date = new Date(x);
          return (
            date.getHours() +
            ":" +
            (date.getMinutes() == 0 ? "00" : date.getMinutes())
          );
        }),
        period: data.period,
      };
    },
    async fetchDailyForecast() {
      const res = await fetch("api/forecast/daily");
      let data = await res.json();
      data = data.dailyForecast;
      this.dailyForecast = [];
      Object.keys(data).forEach((x, index) => {
        if (index == 0) {
          if (data[x].estimate > 0) {
            this.dailyForecast.push({
              day: "Heute",
              estimate: data[x].estimate,
            });
          }
        } else if (index == 1) {
          this.dailyForecast.push({
            day: "Morgen",
            estimate: data[x].estimate,
          });
        } else {
          this.dailyForecast.push(data[x]);
        }
      });
    },
  },
};
</script>
