<template>
  <div class="auth-wrap">
    <div class="auth-card">
      <section class="auth-hero">
        <h1>CityWalk</h1>
        <p>记录路线、故事、地标与活动，让城市探索有内容、有组织，也有社区感。</p>
        <div class="grid two" style="margin-top: 24px;">
          <div class="kpi"><div class="label">路线</div><div class="value">{{ stats.routes }}</div></div>
          <div class="kpi"><div class="label">故事</div><div class="value">{{ stats.stories }}</div></div>
        </div>
      </section>
      <section class="auth-form">
        <div class="section-head" style="margin: 0;">
          <div><h3>{{ mode === 'login' ? '欢迎回来' : '创建账户' }}</h3><p>{{ mode === 'login' ? '使用用户名和密码登录' : '注册后即可创建内容' }}</p></div>
        </div>
        <div v-if="mode === 'register'" class="field"><label>昵称</label><input v-model="form.nickname" class="input" placeholder="城市探索家" /></div>
        <div class="field"><label>用户名</label><input v-model="form.username" class="input" placeholder="citywalker" /></div>
        <div class="field"><label>密码</label><input v-model="form.password" type="password" class="input" placeholder="至少 8 位" /></div>
        <div v-if="mode === 'register'" class="field"><label>所在城市</label><input v-model="form.city" class="input" placeholder="上海" /></div>
        <button class="btn primary" @click="submit">{{ mode === 'login' ? '登录' : '注册' }}</button>
        <button class="btn ghost" @click="mode = mode === 'login' ? 'register' : 'login'">
          {{ mode === 'login' ? '去注册' : '已有账号，去登录' }}
        </button>
        <p class="muted">默认管理员：`admin / Admin123!`</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { get } from '@/lib/api';

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const mode = ref<'login' | 'register'>('login');
const form = reactive({ username: '', password: '', nickname: '', city: '' });
const stats = reactive({ routes: 0, stories: 0 });

onMounted(async () => {
  const data = await get<{ routes: { list: unknown[] }; stories: { list: unknown[] } }>('/v1/recommend').catch(() => null);
  stats.routes = data?.routes?.list?.length ?? 0;
  stats.stories = data?.stories?.list?.length ?? 0;
});

async function submit() {
  if (mode.value === 'login') await auth.login(form.username, form.password);
  else await auth.register(form);
  await router.push(String(route.query.redirect ?? '/'));
}
</script>
