<template>
  <AppShell :title="event?.title || '活动详情'" subtitle="活动信息与报名名单">
    <div v-if="event" class="split">
      <div class="panel panel-inner">
        <p>{{ event.description }}</p>
        <p>{{ event.meetupAddress }} · {{ event.meetupTime }}</p>
        <div class="toolbar">
          <button class="btn primary" @click="join">报名活动</button>
          <button class="btn ghost" @click="cancel">取消报名</button>
        </div>
      </div>
      <div class="panel panel-inner">
        <SectionTitle title="报名名单" desc="活动参与情况" />
        <div class="list">
          <div v-for="item in event.participations" :key="item.id" class="card">
            <h4>{{ item.name }}</h4>
            <p>{{ item.remark }} · 状态 {{ item.status }}</p>
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
import { get, post, del } from '@/lib/api';
import { useAuthStore } from '@/stores/auth';
import type { EventDetail } from '@/lib/types';

const auth = useAuthStore();
const route = useRoute();
const state = ref<EventDetail | null>(null);
const event = computed(() => state.value);

onMounted(async () => {
  state.value = await get<EventDetail>(`/v1/events/${route.params.id}`);
});

async function join() {
  if (!auth.isAuthed) return;
  if (!state.value) return;
  await post(`/v1/events/${state.value.id}/join`, { name: auth.user?.nickname || auth.user?.username || '访客', remark: '' });
  state.value = await get<EventDetail>(`/v1/events/${route.params.id}`);
}

async function cancel() {
  if (!state.value) return;
  await del(`/v1/events/${state.value.id}/join`);
  state.value = await get<EventDetail>(`/v1/events/${route.params.id}`);
}
</script>
