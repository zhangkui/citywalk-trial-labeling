import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import LoginView from '@/views/LoginView.vue';
import HomeView from '@/views/HomeView.vue';
import RoutesView from '@/views/RoutesView.vue';
import RouteDetailView from '@/views/RouteDetailView.vue';
import RouteFormView from '@/views/RouteFormView.vue';
import StoriesView from '@/views/StoriesView.vue';
import StoryDetailView from '@/views/StoryDetailView.vue';
import StoryFormView from '@/views/StoryFormView.vue';
import LandmarksView from '@/views/LandmarksView.vue';
import LandmarkFormView from '@/views/LandmarkFormView.vue';
import EventsView from '@/views/EventsView.vue';
import EventDetailView from '@/views/EventDetailView.vue';
import EventFormView from '@/views/EventFormView.vue';
import ProfileView from '@/views/ProfileView.vue';
import AdminView from '@/views/AdminView.vue';
import NotFoundView from '@/views/NotFoundView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: LoginView },
    { path: '/', name: 'home', component: HomeView, meta: { title: '首页推荐' } },
    { path: '/routes', name: 'routes', component: RoutesView, meta: { title: '路线广场' } },
    { path: '/routes/new', name: 'route-create', component: RouteFormView, meta: { requiresAuth: true, title: '创建路线' } },
    { path: '/routes/:id', name: 'route-detail', component: RouteDetailView, meta: { title: '路线详情' } },
    { path: '/routes/:id/edit', name: 'route-edit', component: RouteFormView, meta: { requiresAuth: true, title: '编辑路线' } },
    { path: '/stories', name: 'stories', component: StoriesView, meta: { title: '城市故事' } },
    { path: '/stories/new', name: 'story-create', component: StoryFormView, meta: { requiresAuth: true, title: '发布故事' } },
    { path: '/stories/:id', name: 'story-detail', component: StoryDetailView, meta: { title: '故事详情' } },
    { path: '/stories/:id/edit', name: 'story-edit', component: StoryFormView, meta: { requiresAuth: true, title: '编辑故事' } },
    { path: '/landmarks', name: 'landmarks', component: LandmarksView, meta: { title: '地标地图' } },
    { path: '/landmarks/new', name: 'landmark-create', component: LandmarkFormView, meta: { requiresAuth: true, title: '添加地标' } },
    { path: '/landmarks/:id/edit', name: 'landmark-edit', component: LandmarkFormView, meta: { requiresAuth: true, title: '编辑地标' } },
    { path: '/events', name: 'events', component: EventsView, meta: { title: '线下活动' } },
    { path: '/events/new', name: 'event-create', component: EventFormView, meta: { requiresAuth: true, title: '创建活动' } },
    { path: '/events/:id', name: 'event-detail', component: EventDetailView, meta: { title: '活动详情' } },
    { path: '/events/:id/edit', name: 'event-edit', component: EventFormView, meta: { requiresAuth: true, title: '编辑活动' } },
    { path: '/profile', name: 'profile', component: ProfileView, meta: { requiresAuth: true, title: '个人主页' } },
    { path: '/admin', name: 'admin', component: AdminView, meta: { requiresAuth: true, requiresRole: '系统管理员', title: '管理后台' } },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundView, meta: { title: '页面不存在' } },
  ],
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();
  if (!auth.user && auth.session?.accessToken) {
    await auth.bootstrap();
  }
  if (to.meta.requiresAuth && !auth.isAuthed) {
    return { name: 'login', query: { redirect: to.fullPath } };
  }
  if (to.meta.requiresRole && !auth.hasRole(String(to.meta.requiresRole))) {
    return { name: 'home' };
  }
  if (to.meta.title) document.title = `${to.meta.title} · CityWalk`;
  return true;
});

export default router;
