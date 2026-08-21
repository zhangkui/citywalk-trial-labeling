<template>
  <AppShell :title="landmarkId ? '编辑地标' : '添加地标'" subtitle="维护地标内容与坐标">
    <div class="panel panel-inner">
      <div class="grid two">
        <div class="field"><label>名称</label><input v-model="form.name" class="input" /></div>
        <div class="field"><label>地址</label><input v-model="form.address" class="input" /></div>
        <div class="field"><label>经度</label><input v-model="form.lng" class="input" /></div>
        <div class="field"><label>纬度</label><input v-model="form.lat" class="input" /></div>
        <div class="field"><label>分类 ID</label><input v-model="form.categoryId" class="input" /></div>
        <div class="field"><label>封面图</label><input v-model="form.coverImage" class="input" /></div>
      </div>
      <div class="field" style="margin-top: 14px;"><label>简介</label><textarea v-model="form.description" class="textarea" /></div>
      <div class="toolbar" style="margin-top: 14px;"><button class="btn primary" @click="submit">保存</button></div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppShell from '@/components/AppShell.vue';
import { get, post, put } from '@/lib/api';
import type { LandmarkItem } from '@/lib/types';

const route = useRoute();
const router = useRouter();
const landmarkId = computed(() => Number(route.params.id || 0) || 0);
const form = ref({ name: '', address: '', lat: '31.2304', lng: '121.4737', categoryId: '1', description: '', coverImage: '' });

onMounted(async () => {
  if (!landmarkId.value) return;
  const data = await get<LandmarkItem>(`/v1/landmarks/${landmarkId.value}`);
  form.value = {
    name: data.name,
    address: data.address,
    lat: String(data.lat),
    lng: String(data.lng),
    categoryId: String(data.categoryId || ''),
    description: data.description,
    coverImage: data.coverImage,
  };
});

async function submit() {
  const body = { ...form.value, lat: Number(form.value.lat), lng: Number(form.value.lng), categoryId: form.value.categoryId ? Number(form.value.categoryId) : null };
  if (landmarkId.value) await put(`/v1/landmarks/${landmarkId.value}`, body);
  else await post('/v1/landmarks', body);
  await router.push('/landmarks');
}
</script>
