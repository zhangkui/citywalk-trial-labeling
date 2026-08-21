<template>
  <div class="app-shell">
    <aside class="side">
      <div class="brand">
        <h1>CityWalk</h1>
        <p>城市漫步与文化导览平台</p>
      </div>
      <div class="nav-group">
        <RouterLink class="nav-link" to="/">首页 <span>01</span></RouterLink>
        <RouterLink class="nav-link" to="/routes">路线 <span>02</span></RouterLink>
        <RouterLink class="nav-link" to="/stories">故事 <span>03</span></RouterLink>
        <RouterLink class="nav-link" to="/landmarks">地标 <span>04</span></RouterLink>
        <RouterLink class="nav-link" to="/events">活动 <span>05</span></RouterLink>
        <RouterLink class="nav-link" to="/profile">个人页 <span>06</span></RouterLink>
        <RouterLink v-if="adminVisible" class="nav-link" to="/admin">管理后台 <span>07</span></RouterLink>
      </div>
      <div class="spacer"></div>
      <div v-if="auth.user" class="card">
        <h4>{{ auth.user.nickname || auth.user.username }}</h4>
        <p class="muted">{{ auth.user.city || '未设置城市' }}</p>
        <div class="toolbar" style="margin-top: 12px;">
          <button class="btn secondary" @click="auth.logout()">退出</button>
        </div>
      </div>
      <RouterLink v-else class="btn primary" to="/login">登录 / 注册</RouterLink>
    </aside>
    <main class="shell-main">
      <div class="topbar">
        <div class="headline">
          <h2>{{ title }}</h2>
          <p>{{ subtitle }}</p>
        </div>
        <div class="toolbar">
          <slot name="actions" />
        </div>
      </div>
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useAuthStore } from '@/stores/auth';

defineProps<{ title: string; subtitle?: string }>();

const auth = useAuthStore();
const adminVisible = computed(() => auth.hasRole('系统管理员'));
</script>
