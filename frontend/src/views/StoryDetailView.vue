<template>
  <AppShell :title="story?.title || '故事详情'" subtitle="完整内容、媒体和互动">
    <div v-if="story" class="split">
      <div class="panel panel-inner">
        <p>{{ story.content }}</p>
        <div class="toolbar">
          <button class="btn secondary" @click="like">点赞</button>
          <RouterLink class="btn ghost" :to="`/stories/${story.id}/edit`">编辑</RouterLink>
        </div>
      </div>
      <div class="panel panel-inner">
        <SectionTitle title="多媒体" desc="图片与音频资源" />
        <div class="list">
          <div v-for="item in story.media" :key="item.id ?? item.order" class="card">
            <h4>{{ item.type === 1 ? '图片' : '音频' }}</h4>
            <p>{{ item.url }}</p>
          </div>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import { get, post } from '@/lib/api';
import type { StoryDetail } from '@/lib/types';

const route = useRoute();
const story = computed<StoryDetail | null>(() => state.value);
const state = ref<StoryDetail | null>(null);

onMounted(async () => {
  state.value = await get<StoryDetail>(`/v1/stories/${route.params.id}`);
});

async function like() {
  if (!state.value) return;
  await post(`/v1/stories/${state.value.id}/like`);
  state.value = await get<StoryDetail>(`/v1/stories/${route.params.id}`);
}
</script>
