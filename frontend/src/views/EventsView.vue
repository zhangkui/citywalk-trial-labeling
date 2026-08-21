<template>
  <AppShell title="线下活动" subtitle="路线延伸到线下的组织入口">
    <div class="panel panel-inner">
      <div class="toolbar" style="margin-bottom: 14px;">
        <select v-model="status" class="select" style="max-width: 180px;">
          <option :value="-1">全部状态</option>
          <option :value="0">待开始</option>
          <option :value="1">进行中</option>
          <option :value="2">已结束</option>
        </select>
        <button class="btn primary" @click="reload">刷新</button>
        <RouterLink class="btn secondary" to="/events/new">创建活动</RouterLink>
      </div>
      <div class="grid two">
        <div class="list">
          <div v-for="item in items" :key="item.id" class="card">
            <h4>{{ item.title }}</h4>
            <p>{{ item.meetupAddress }} · {{ item.currentParticipants }}/{{ item.maxParticipants }} 人</p>
            <RouterLink class="btn secondary" :to="`/events/${item.id}`">详情</RouterLink>
          </div>
        </div>
        <div class="panel panel-inner"><SectionTitle title="活动说明" desc="支持报名、取消和签到" /></div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import { get } from '@/lib/api';
import type { EventItem } from '@/lib/types';

const status = ref(-1);
const items = ref<EventItem[]>([]);

async function reload() {
  const params: Record<string, unknown> = { page: 1, pageSize: 20 };
  if (status.value >= 0) params.status = status.value;
  const data = await get<{ list: EventItem[] }>('/v1/events', params);
  items.value = data.list;
}

onMounted(reload);
</script>
