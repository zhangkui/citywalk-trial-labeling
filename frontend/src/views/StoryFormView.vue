<template>
  <AppShell :title="storyId ? '编辑故事' : '发布故事'" subtitle="图文与多媒体内容编辑">
    <div class="panel panel-inner">
      <div class="grid two">
        <div class="field"><label>标题</label><input v-model="form.title" class="input" /></div>
        <div class="field"><label>封面图</label><input v-model="form.coverImage" class="input" /></div>
        <div class="field"><label>关联路线 ID</label><input v-model="form.routeId" class="input" /></div>
        <div class="field"><label>地标 ID 列表</label><input v-model="form.landmarkIds" class="input" placeholder="1,2,3" /></div>
      </div>
      <div class="field" style="margin-top: 14px;"><label>正文</label><textarea v-model="form.content" class="textarea" /></div>
      <div class="field"><label>多媒体 JSON</label><textarea v-model="mediaText" class="textarea" /></div>
      <div class="toolbar" style="margin-top: 14px;"><button class="btn primary" @click="submit">保存</button></div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppShell from '@/components/AppShell.vue';
import { get, post, put } from '@/lib/api';
import type { StoryDetail } from '@/lib/types';

const route = useRoute();
const router = useRouter();
const storyId = computed(() => Number(route.params.id || 0) || 0);
const form = ref({ title: '', content: '', coverImage: '', routeId: '', landmarkIds: '' });
const mediaText = ref('[{"type":1,"url":"https://example.com/photo.jpg"}]');

onMounted(async () => {
  if (!storyId.value) return;
  const data = await get<StoryDetail>(`/v1/stories/${storyId.value}`);
  form.value = {
    title: data.title,
    content: data.content,
    coverImage: data.coverImage,
    routeId: String(data.routeId || ''),
    landmarkIds: data.landmarkIds,
  };
  mediaText.value = JSON.stringify(data.media ?? [], null, 2);
});

async function submit() {
  const body = {
    title: form.value.title,
    content: form.value.content,
    coverImage: form.value.coverImage,
    routeId: form.value.routeId ? Number(form.value.routeId) : null,
    landmarkIds: form.value.landmarkIds ? form.value.landmarkIds.split(',').map((item) => Number(item.trim())).filter(Boolean) : [],
    media: JSON.parse(mediaText.value || '[]'),
  };
  if (storyId.value) await put(`/v1/stories/${storyId.value}`, body);
  else await post('/v1/stories', body);
  await router.push('/stories');
}
</script>
