<template>
  <v-app>
   <!-- <v-app-bar app> </v-app-bar> -->

    <v-main>
      <v-container>
        <v-row dense justify="center">
          <v-col class="ma-4" align="center">
            <v-container class="text-h2 mb-5">
              Photovoltaik
            </v-container>
            <v-progress-circular :color="percentColor" width="20" rotate="-90" size="300" :value="power.batteryPercent"
              :indeterminate="power.batteryPercent == null">

              <div class="text-h3" v-if="power.batteryPercent !== null">
                {{ power.batteryPercent + "%"}}
              </div>

            </v-progress-circular>
            <v-container v-if="power.batteryPower < 0" class="text-h5 red--text">
              Batterie entlädt: {{ Math.abs(power.batteryPower) }} W
            </v-container>
            <v-container v-else-if="power.batteryPower > 0" class="text-h5 green--text">
              Batterie ladet mit: {{ Math.abs(power.batteryPower) }} W
            </v-container>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>


export default {
  name: 'App',

  components: {
    //HelloWorld,
  },

  data() {
    return {
      power: {
        "batteryPercent": null,
        "gridPower": null,
        "pvPower": null,
        "powerConsumption": null,
        "batteryPower": null,
      },
      interval: 0,
    }
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
    }
  },

  created() {
    //this.login();
    this.fetchData();
    setInterval(this.fetchData, 2000);
  },

  methods: {
    async fetchData() {
      let res = await fetch('http://192.168.2.155/api/data');
      this.power = await res.json()
    }
  }
};
</script>
