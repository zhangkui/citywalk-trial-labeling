<template>
  <AppShell title="路线广场" subtitle="发现城市漫步路线">
    <div class="panel panel-inner">
      <div class="toolbar" style="margin-bottom: 14px;">
        <input v-model="filters.city" class="input" placeholder="城市" style="max-width: 180px;" />
        <select v-model.number="filters.themeId" class="select" style="max-width: 180px;">
          <option :value="0">全部主题</option>
          <option v-for="item in themeOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
        </select>
        <select v-model="filters.sort" class="select" style="max-width: 180px;">
          <option value="hot">热门</option>
          <option value="latest">最新</option>
          <option value="rating">评分</option>
        </select>
        <button class="btn primary" @click="reload">筛选</button>
        <RouterLink class="btn secondary" to="/routes/new">创建路线</RouterLink>
      </div>
      <div class="grid two">
        <div class="list">
          <div v-for="item in routes.items" :key="item.id" class="card">
            <div class="toolbar" style="justify-content: space-between;">
              <h4>{{ item.title }}</h4>
              <StatusBadge kind="route" :value="item.status" />
            </div>
            <p>{{ item.city }} · {{ item.favoriteCount }} 收藏 · {{ item.rating.toFixed?.(1) ?? item.rating }} 分</p>
            <p>{{ item.description.slice(0, 120) }}...</p>
            <div class="toolbar">
              <RouterLink class="btn secondary" :to="`/routes/${item.id}`">详情</RouterLink>
              <RouterLink class="btn ghost" :to="`/routes/${item.id}/edit`">编辑</RouterLink>
            </div>
          </div>
        </div>
        <div class="panel panel-inner">
          <SectionTitle title="列表说明" desc="所有数据都来自 Go API" />
          <div class="kpi"><div class="label">总数</div><div class="value">{{ routes.total }}</div></div>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import { themeOptions } from '@/lib/catalog';
import { useRouteStore } from '@/stores/routes';

const routes = useRouteStore();
const filters = reactive({ city: '', themeId: 0, sort: 'hot' });

async function reload() {
  await routes.loadList(filters);
}

onMounted(reload);
</script>
