<template>
  <v-app-bar flat class="px-4">
    <v-icon color="emerald" size="32" class="mr-3">mdi-solar-power-variant</v-icon>
    <v-app-bar-title class="font-weight-bold">
      PV Tracker
    </v-app-bar-title>
    
    <v-spacer></v-spacer>
    
    <div class="d-flex align-center">
      <span class="pulse-dot mr-2"></span>
      <span class="text-caption text-secondary font-weight-medium d-none d-sm-inline">Live-Verbindung</span>
    </div>
  </v-app-bar>
</template>

<script setup>
  //
</script>

<style scoped>
.pulse-dot {
  width: 8px;
  height: 8px;
  background-color: var(--color-pv);
  border-radius: 50%;
  display: inline-block;
  box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
  }
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 6px rgba(16, 185, 129, 0);
  }
  100% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0);
  }
}
</style>
