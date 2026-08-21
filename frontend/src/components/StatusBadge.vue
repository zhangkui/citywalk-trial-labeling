<template>
  <span :class="['badge', tone]">{{ label }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { eventStatusOptions, routeStatusOptions, storyStatusOptions } from '@/lib/catalog';

const props = defineProps<{ kind: 'route' | 'story' | 'event' | 'plain'; value?: number; label?: string }>();

const label = computed(() => {
  if (props.label) return props.label;
  const source = props.kind === 'route' ? routeStatusOptions : props.kind === 'story' ? storyStatusOptions : props.kind === 'event' ? eventStatusOptions : [];
  return source.find((item) => item.value === props.value)?.label ?? String(props.value ?? '');
});

const tone = computed(() => {
  if (props.kind === 'plain') return '';
  if (props.value === 2) return 'good';
  if (props.value === 1) return 'warn';
  if (props.value === 3) return 'bad';
  return '';
});
</script>
