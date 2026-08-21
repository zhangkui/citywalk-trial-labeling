<template>
  <AppShell title="首页推荐" subtitle="热度、距离和兴趣的综合入口">
    <div class="grid two">
      <div class="panel panel-inner">
        <SectionTitle title="推荐路线" desc="按热度排列的路线精选" />
        <div class="list">
          <div v-for="item in recommend.routes?.slice(0, 4) ?? []" :key="item.id" class="card">
            <h4>{{ item.title }}</h4>
            <p>{{ item.city }} · {{ item.rating.toFixed?.(1) ?? item.rating }} 分 · 收藏 {{ item.favoriteCount }}</p>
            <RouterLink class="btn secondary" :to="`/routes/${item.id}`">查看路线</RouterLink>
          </div>
        </div>
      </div>
      <div class="panel panel-inner">
        <SectionTitle title="附近地标" desc="基于定位的地标推荐" />
        <div class="list">
          <div v-for="item in recommend.nearbyLandmarks?.slice(0, 4) ?? []" :key="item.id" class="card">
            <h4>{{ item.name }}</h4>
            <p>{{ item.address }} · {{ item.distance?.toFixed?.(0) ?? item.distance }} m</p>
          </div>
        </div>
      </div>
    </div>
    <div class="panel panel-inner" style="margin-top: 18px;">
      <SectionTitle title="热点故事" desc="最新与最热内容汇聚" />
      <div class="grid three">
        <div v-for="item in recommend.stories?.slice(0, 6) ?? []" :key="item.id" class="card">
          <h4>{{ item.title }}</h4>
          <p>{{ item.content.slice(0, 80) }}...</p>
          <RouterLink class="btn secondary" :to="`/stories/${item.id}`">查看故事</RouterLink>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import { get } from '@/lib/api';

const recommend = reactive<{ routes: any[]; stories: any[]; nearbyLandmarks: any[] }>({ routes: [], stories: [], nearbyLandmarks: [] });

onMounted(async () => {
  const data = await get<typeof recommend>('/v1/recommend', { lat: 31.2304, lng: 121.4737 }).catch(() => null);
  if (data) Object.assign(recommend, data);
});
</script>
