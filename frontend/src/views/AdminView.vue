<template>
  <AppShell title="管理后台" subtitle="用户、角色和权限">
    <div class="grid two">
      <div class="panel panel-inner">
        <SectionTitle title="用户管理" desc="启用、停用和新建用户" />
        <div class="toolbar" style="margin-bottom: 14px;">
          <input v-model="keyword" class="input" placeholder="搜索用户" style="max-width: 180px;" />
          <button class="btn primary" @click="loadUsers">刷新</button>
        </div>
        <div class="list">
          <div v-for="item in users" :key="item.id" class="card">
            <h4>{{ item.username }}</h4>
            <p>{{ item.nickname }} · {{ item.city }} · {{ item.roles.join(', ') }}</p>
            <div class="toolbar">
              <button class="btn secondary" @click="toggleStatus(item)">切换状态</button>
            </div>
          </div>
        </div>
      </div>
      <div class="panel panel-inner">
        <SectionTitle title="角色与权限" desc="查看系统配置" />
        <div class="grid two">
          <div>
            <h4>角色</h4>
            <div class="list">
              <div v-for="item in roles" :key="item.id" class="card"><strong>{{ item.name }}</strong><p>{{ item.description }}</p></div>
            </div>
          </div>
          <div>
            <h4>权限</h4>
            <div class="list">
              <div v-for="item in permissions" :key="item.id" class="card"><strong>{{ item.name }}</strong><p>{{ item.description }}</p></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import { get, put } from '@/lib/api';
import type { AdminUser, PermissionItem, RoleItem } from '@/lib/types';

const keyword = ref('');
const users = ref<AdminUser[]>([]);
const roles = ref<RoleItem[]>([]);
const permissions = ref<PermissionItem[]>([]);

async function loadUsers() {
  const data = await get<{ list: AdminUser[] }>('/v1/admin/users', { keyword: keyword.value, page: 1, pageSize: 20 });
  users.value = data.list;
}

async function loadStatic() {
  roles.value = await get<RoleItem[]>('/v1/admin/roles');
  permissions.value = await get<PermissionItem[]>('/v1/admin/permissions');
}

async function toggleStatus(item: AdminUser) {
  await put(`/v1/admin/users/${item.id}/status`, { status: item.status === 1 ? 0 : 1 });
  await loadUsers();
}

onMounted(async () => {
  await Promise.all([loadUsers(), loadStatic()]);
});
</script>
