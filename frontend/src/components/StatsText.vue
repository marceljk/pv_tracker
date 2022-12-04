<template>
    <div v-if="response">
      {{ response }}
    </div>
    <div v-else>
      <div v-if="value < 0" class="red--text">
        {{ negResponse }}
      </div>
      <div v-else-if="value > 0" class="green--text">
        {{ posResponse }}
      </div>
    </div>
</template>

<script>
export default {
  name: 'StatsText',

  props: {
    value: null,
    negResponse: String,
    posResponse: String,
    response: null,
  }
}
</script>
