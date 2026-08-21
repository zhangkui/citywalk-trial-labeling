import { computed, ref } from 'vue';
import { defineStore } from 'pinia';
import { del, get, post, put, loadSession, saveSession } from '@/lib/api';
import type { BadgeItem, CurrentUser, SessionData, UserSummary } from '@/lib/types';

export const useAuthStore = defineStore('auth', () => {
  const session = ref<SessionData | null>(loadSession());
  const user = ref<CurrentUser | null>(null);
  const badges = ref<BadgeItem[]>([]);

  const isAuthed = computed(() => Boolean(session.value?.accessToken));
  const roles = computed(() => user.value?.roles ?? session.value?.user.roles ?? []);
  const permissions = computed(() => user.value?.permissions ?? session.value?.user.permissions ?? []);

  const syncSession = (next: SessionData | null) => {
    session.value = next;
    saveSession(next);
  };

  const applyUser = (data: UserSummary & Partial<CurrentUser>) => {
    user.value = {
      id: data.id,
      username: data.username,
      nickname: data.nickname,
      avatar: data.avatar,
      bio: data.bio,
      city: data.city,
      status: data.status ?? 1,
      roles: data.roles ?? [],
      permissions: data.permissions ?? [],
      createdAt: data.createdAt ?? new Date().toISOString(),
      updatedAt: data.updatedAt ?? new Date().toISOString(),
    };
  };

  async function bootstrap() {
    if (!session.value?.accessToken) return;
    try {
      await loadMe();
      await loadBadges();
    } catch {
      await logout(false);
    }
  }

  async function login(username: string, password: string) {
    const result = await post<SessionData>('/v1/auth/login', { username, password });
    syncSession(result);
    applyUser(result.user);
    await loadBadges();
  }

  async function register(payload: { username: string; password: string; nickname: string; city: string }) {
    const result = await post<SessionData>('/v1/auth/register', payload);
    syncSession(result);
    applyUser(result.user);
    await loadBadges();
  }

  async function loadMe() {
    const data = await get<CurrentUser>('/v1/auth/me');
    applyUser(data);
    return data;
  }

  async function loadBadges() {
    if (!isAuthed.value) return [];
    badges.value = await get<BadgeItem[]>('/v1/auth/badges');
    return badges.value;
  }

  async function updateProfile(payload: { nickname: string; avatar: string; bio: string; city: string }) {
    const data = await put<CurrentUser>('/v1/auth/profile', payload);
    applyUser(data);
    return data;
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    await put('/v1/auth/password', { oldPassword, newPassword });
  }

  async function logout(callApi = true) {
    if (callApi && session.value?.accessToken) {
      try {
        await del('/v1/auth/logout', { refreshToken: session.value.refreshToken });
      } catch {
        // ignore
      }
    }
    syncSession(null);
    user.value = null;
    badges.value = [];
  }

  const hasRole = (name: string) => roles.value.includes(name);
  const hasPermission = (name: string) => permissions.value.includes(name);

  return {
    session,
    user,
    badges,
    isAuthed,
    roles,
    permissions,
    bootstrap,
    login,
    register,
    loadMe,
    loadBadges,
    updateProfile,
    changePassword,
    logout,
    hasRole,
    hasPermission,
  };
});
