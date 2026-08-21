<template>
  <AppShell :title="route?.title || '路线详情'" subtitle="路线、途经点和互动信息">
    <div v-if="route" class="split">
      <div class="panel panel-inner">
        <div class="toolbar" style="justify-content: space-between;">
          <StatusBadge kind="route" :value="route.status" />
          <div class="toolbar">
            <button class="btn secondary" @click="toggleFavorite">收藏</button>
            <button class="btn ghost" @click="rate">评分 5 分</button>
          </div>
        </div>
        <p>{{ route.city }} · {{ route.favoriteCount }} 收藏 · {{ route.rating.toFixed?.(1) ?? route.rating }} 分</p>
        <p>{{ route.description }}</p>
        <SimpleMap :center="[route.startLat || 31.2304, route.startLng || 121.4737]" :markers="markers" :path="path" />
      </div>
      <div class="panel panel-inner">
        <SectionTitle title="途经点" desc="路线节点信息" />
        <div class="list">
          <div v-for="item in route.waypoints" :key="item.id ?? item.order" class="card">
            <h4>{{ item.name }}</h4>
            <p>{{ item.lat }}, {{ item.lng }} · 停留 {{ item.stayDuration }} 分钟</p>
          </div>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import SimpleMap from '@/components/SimpleMap.vue';
import { useRouteStore } from '@/stores/routes';

const routeStore = useRouteStore();
const route = computed(() => routeStore.current);
const currentRoute = useRoute();

onMounted(async () => {
  await routeStore.loadDetail(Number(currentRoute.params.id));
});

const markers = computed(() => (route.value?.waypoints ?? []).map((item, index) => ({ lat: item.lat, lng: item.lng, label: item.name, color: index === 0 ? '#6ee7b7' : '#f4a261' })));
const path = computed(() => (route.value?.waypoints ?? []).map((item) => [item.lat, item.lng] as [number, number]));

async function toggleFavorite() {
  if (!route.value) return;
  await routeStore.toggleFavorite(route.value.id, true);
}

async function rate() {
  if (!route.value) return;
  await routeStore.rate(route.value.id, 5);
  await routeStore.loadDetail(route.value.id);
}
</script>
