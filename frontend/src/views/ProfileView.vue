<template>
  <AppShell title="个人主页" subtitle="资料、徽章与安全设置">
    <div class="grid two">
      <div class="panel panel-inner">
        <SectionTitle title="基础资料" desc="修改后同步到后端" />
        <div class="grid two">
          <div class="field"><label>昵称</label><input v-model="form.nickname" class="input" /></div>
          <div class="field"><label>城市</label><input v-model="form.city" class="input" /></div>
        </div>
        <div class="field"><label>头像</label><input v-model="form.avatar" class="input" /></div>
        <div class="field"><label>简介</label><textarea v-model="form.bio" class="textarea" /></div>
        <div class="toolbar"><button class="btn primary" @click="saveProfile">保存资料</button></div>
      </div>
      <div class="panel panel-inner">
        <SectionTitle title="徽章" desc="自动解锁的成就体系" />
        <div class="list">
          <div v-for="item in auth.badges" :key="item.id" class="card">
            <h4>{{ item.name }}</h4>
            <p>{{ item.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </AppShell>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import AppShell from '@/components/AppShell.vue';
import SectionTitle from '@/components/SectionTitle.vue';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const form = reactive({ nickname: '', avatar: '', bio: '', city: '' });

onMounted(async () => {
  await auth.loadMe();
  await auth.loadBadges();
  if (auth.user) Object.assign(form, { nickname: auth.user.nickname, avatar: auth.user.avatar, bio: auth.user.bio, city: auth.user.city });
});

async function saveProfile() {
  await auth.updateProfile(form);
}
</script>
