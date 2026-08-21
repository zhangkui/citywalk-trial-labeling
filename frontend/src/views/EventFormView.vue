<template>
  <AppShell :title="eventId ? '编辑活动' : '创建活动'" subtitle="配置路线、时间、人数和费用">
    <div class="panel panel-inner">
      <div class="grid two">
        <div class="field"><label>标题</label><input v-model="form.title" class="input" /></div>
        <div class="field"><label>路线 ID</label><input v-model="form.routeId" class="input" /></div>
        <div class="field"><label>集合地点</label><input v-model="form.meetupAddress" class="input" /></div>
        <div class="field"><label>集合时间</label><input v-model="form.meetupTime" class="input" placeholder="2026-08-20T10:00:00" /></div>
        <div class="field"><label>开始时间</label><input v-model="form.startTime" class="input" /></div>
        <div class="field"><label>结束时间</label><input v-model="form.endTime" class="input" /></div>
        <div class="field"><label>人数上限</label><input v-model.number="form.maxParticipants" class="input" /></div>
        <div class="field"><label>费用(分)</label><input v-model.number="form.fee" class="input" /></div>
      </div>
      <div class="field" style="margin-top: 14px;"><label>描述</label><textarea v-model="form.description" class="textarea" /></div>
      <div class="toolbar" style="margin-top: 14px;"><button class="btn primary" @click="submit">保存</button></div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppShell from '@/components/AppShell.vue';
import { get, post, put } from '@/lib/api';
import type { EventDetail } from '@/lib/types';

const route = useRoute();
const router = useRouter();
const eventId = computed(() => Number(route.params.id || 0) || 0);
const form = ref({ title: '', routeId: '', meetupAddress: '', meetupTime: '', startTime: '', endTime: '', maxParticipants: 20, fee: 0, description: '' });

onMounted(async () => {
  if (!eventId.value) return;
  const data = await get<EventDetail>(`/v1/events/${eventId.value}`);
  form.value = {
    title: data.title,
    routeId: String(data.routeId || ''),
    meetupAddress: data.meetupAddress,
    meetupTime: data.meetupTime,
    startTime: data.startTime,
    endTime: data.endTime,
    maxParticipants: data.maxParticipants,
    fee: data.fee,
    description: data.description,
  };
});

async function submit() {
  const body = {
    ...form.value,
    routeId: form.value.routeId ? Number(form.value.routeId) : null,
    maxParticipants: Number(form.value.maxParticipants),
    fee: Number(form.value.fee),
    meetupTime: form.value.meetupTime,
    startTime: form.value.startTime,
    endTime: form.value.endTime,
  };
  if (eventId.value) await put(`/v1/events/${eventId.value}`, body);
  else await post('/v1/events', body);
  await router.push('/events');
}
</script>
