<template>
  <AppShell :title="routeId ? '编辑路线' : '创建路线'" subtitle="填写路线基本信息和途经点">
    <div class="panel panel-inner">
      <div class="grid two">
        <div class="field"><label>标题</label><input v-model="form.title" class="input" /></div>
        <div class="field"><label>城市</label><input v-model="form.city" class="input" /></div>
        <div class="field"><label>主题 ID</label><input v-model.number="form.themeId" class="input" type="number" /></div>
        <div class="field"><label>封面图</label><input v-model="form.coverImage" class="input" /></div>
        <div class="field"><label>起点经纬度</label><input v-model="form.start" class="input" placeholder="31.23,121.47" /></div>
        <div class="field"><label>终点经纬度</label><input v-model="form.end" class="input" placeholder="31.24,121.48" /></div>
      </div>
      <div class="field" style="margin-top: 14px;"><label>描述</label><textarea v-model="form.description" class="textarea" /></div>
      <div class="field"><label>途经点 JSON</label><textarea v-model="waypointsText" class="textarea" /></div>
      <div class="toolbar" style="margin-top: 14px;"><button class="btn primary" @click="submit">保存</button></div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppShell from '@/components/AppShell.vue';
import { useRouteStore } from '@/stores/routes';

const router = useRouter();
const route = useRoute();
const routeStore = useRouteStore();
const routeId = computed(() => Number(route.params.id || 0) || 0);
const form = ref({ title: '', description: '', city: '', themeId: 1, coverImage: '', start: '', end: '' });
const waypointsText = ref('[{"name":"起点","lat":31.2304,"lng":121.4737,"stayDuration":15}]');

onMounted(async () => {
  if (!routeId.value) return;
  const data = await routeStore.loadDetail(routeId.value);
  form.value = {
    title: data.title,
    description: data.description,
    city: data.city,
    themeId: data.themeId || 1,
    coverImage: data.coverImage,
    start: `${data.startLat},${data.startLng}`,
    end: `${data.endLat},${data.endLng}`,
  };
  waypointsText.value = JSON.stringify(data.waypoints ?? [], null, 2);
});

async function submit() {
  const [startLat, startLng] = form.value.start.split(',').map(Number);
  const [endLat, endLng] = form.value.end.split(',').map(Number);
  const waypoints = JSON.parse(waypointsText.value || '[]');
  await routeStore.save({
    id: routeId.value || undefined,
    title: form.value.title,
    description: form.value.description,
    city: form.value.city,
    themeId: form.value.themeId,
    coverImage: form.value.coverImage,
    startLat,
    startLng,
    endLat,
    endLng,
    waypoints,
  });
  await router.push('/routes');
}
</script>
