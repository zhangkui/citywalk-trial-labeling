<template>
  <AppShell title="城市故事" subtitle="内容创作与阅读">
    <div class="panel panel-inner">
      <div class="toolbar" style="margin-bottom: 14px;">
        <select v-model="sort" class="select" style="max-width: 180px;">
          <option value="hot">热门</option>
          <option value="latest">最新</option>
        </select>
        <button class="btn primary" @click="reload">刷新</button>
        <RouterLink class="btn secondary" to="/stories/new">发布故事</RouterLink>
      </div>
      <div class="grid two">
        <div class="list">
          <div v-for="item in items" :key="item.id" class="card">
            <h4>{{ item.title }}</h4>
            <p>{{ item.content.slice(0, 120) }}...</p>
            <p>{{ item.likeCount }} 赞 · {{ item.viewCount }} 阅读</p>
            <RouterLink class="btn secondary" :to="`/stories/${item.id}`">查看</RouterLink>
          </div>
        </div>
        <div class="panel panel-inner"><SectionTitle title="故事说明" desc="支持发布图文与多媒体" /></div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import { get } from '@/lib/api';
import type { StoryItem } from '@/lib/types';

const sort = ref('hot');
const items = ref<StoryItem[]>([]);

async function reload() {
  const data = await get<{ list: StoryItem[] }>('/v1/stories', { sort: sort.value, page: 1, pageSize: 20 });
  items.value = data.list;
}

onMounted(reload);
</script>
