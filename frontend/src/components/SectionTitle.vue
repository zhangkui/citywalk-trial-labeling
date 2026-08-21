<template>
  <div class="section-head">
    <div>
      <h3>{{ title }}</h3>
      <p v-if="desc">{{ desc }}</p>
    </div>
    <div class="toolbar"><slot /></div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; desc?: string }>();
</script>
