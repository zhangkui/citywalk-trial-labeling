import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

type FavoriteBucket = 'routes' | 'stories' | 'landmarks';

const key = 'citywalk.favorites';

export const useFavoritesStore = defineStore('favorites', () => {
  const state = ref<Record<FavoriteBucket, number[]>>(load());

  function load(): Record<FavoriteBucket, number[]> {
    try {
      const raw = localStorage.getItem(key);
      if (!raw) return { routes: [], stories: [], landmarks: [] };
      return JSON.parse(raw) as Record<FavoriteBucket, number[]>;
    } catch {
      return { routes: [], stories: [], landmarks: [] };
    }
  }

  function persist() {
    localStorage.setItem(key, JSON.stringify(state.value));
  }

  function toggle(bucket: FavoriteBucket, id: number) {
    const items = state.value[bucket];
    const index = items.indexOf(id);
    if (index >= 0) items.splice(index, 1);
    else items.push(id);
    persist();
  }

  const routeIds = computed(() => state.value.routes);
  const storyIds = computed(() => state.value.stories);
  const landmarkIds = computed(() => state.value.landmarks);

  return { state, routeIds, storyIds, landmarkIds, toggle };
});
