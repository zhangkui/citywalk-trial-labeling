import { ref } from 'vue';
import { defineStore } from 'pinia';
import { del, get, post, put } from '@/lib/api';
import type { PageResult, RouteDetail, RouteItem, RouteWaypoint } from '@/lib/types';

export const useRouteStore = defineStore('routes', () => {
  const items = ref<RouteItem[]>([]);
  const current = ref<RouteDetail | null>(null);
  const total = ref(0);
  const loading = ref(false);

  async function loadList(params: { city?: string; themeId?: number; sort?: string; page?: number; pageSize?: number; minRate?: number } = {}) {
    loading.value = true;
    try {
      const data = await get<PageResult<RouteItem>>('/v1/routes', params);
      items.value = data.list;
      total.value = data.total;
      return data;
    } finally {
      loading.value = false;
    }
  }

  async function loadDetail(id: number) {
    current.value = await get<RouteDetail>(`/v1/routes/${id}`);
    return current.value;
  }

  async function save(payload: Partial<RouteDetail> & { title: string; description: string; city: string; waypoints: RouteWaypoint[] }) {
    const body = {
      title: payload.title,
      description: payload.description,
      themeId: payload.themeId || null,
      city: payload.city,
      startLat: payload.startLat ?? null,
      startLng: payload.startLng ?? null,
      endLat: payload.endLat ?? null,
      endLng: payload.endLng ?? null,
      totalDistance: payload.totalDistance ?? null,
      duration: payload.duration ?? null,
      difficulty: payload.difficulty ?? 1,
      coverImage: payload.coverImage ?? '',
      waypoints: payload.waypoints,
      status: payload.status ?? 1,
    };
    if (payload.id) {
      return put<RouteDetail>(`/v1/routes/${payload.id}`, body);
    }
    return post<RouteDetail>('/v1/routes', body);
  }

  async function remove(id: number) {
    await del(`/v1/routes/${id}`);
  }

  async function toggleFavorite(id: number, add: boolean) {
    return add ? post(`/v1/routes/${id}/favorite`) : del(`/v1/routes/${id}/favorite`);
  }

  async function rate(id: number, score: number) {
    await post(`/v1/routes/${id}/rate`, { score });
  }

  return { items, current, total, loading, loadList, loadDetail, save, remove, toggleFavorite, rate };
});
