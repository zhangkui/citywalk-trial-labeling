<template>
  <AppShell title="地标地图" subtitle="展示地标、分类和附近信息">
    <div class="split">
      <div class="panel panel-inner">
        <SectionTitle title="地标列表" desc="真实调用后端接口" />
        <div class="toolbar" style="margin-bottom: 14px;">
          <select v-model.number="categoryId" class="select" style="max-width: 180px;">
            <option :value="0">全部分类</option>
            <option v-for="item in landmarkCategories" :key="item.value" :value="item.value">{{ item.label }}</option>
          </select>
          <button class="btn primary" @click="reload">刷新</button>
          <RouterLink class="btn secondary" to="/landmarks/new">添加地标</RouterLink>
        </div>
        <div class="list">
          <div v-for="item in items" :key="item.id" class="card">
            <h4>{{ item.name }}</h4>
            <p>{{ item.address }} · {{ item.rating.toFixed?.(1) ?? item.rating }} 分</p>
          </div>
        </div>
      </div>
      <SimpleMap :center="center" :markers="markers" />
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import SimpleMap from '@/components/SimpleMap.vue';
import { get } from '@/lib/api';
import { landmarkCategories } from '@/lib/catalog';
import type { LandmarkItem } from '@/lib/types';

const categoryId = ref(0);
const items = ref<LandmarkItem[]>([]);
const center = ref<[number, number]>([31.2304, 121.4737]);
const markers = computed(() => (Array.isArray(items.value) ? items.value : []).map((item) => ({ lat: item.lat, lng: item.lng, label: item.name })));

async function reload() {
  const data = await get<{ list: LandmarkItem[] }>('/v1/landmarks', { categoryId: categoryId.value, page: 1, pageSize: 30 });
  items.value = Array.isArray(data.list) ? data.list : [];
}

onMounted(reload);
</script>
