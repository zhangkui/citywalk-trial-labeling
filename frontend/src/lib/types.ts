export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PageResult<T> {
  list: T[];
  total: number;
}

export interface UserSummary {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  bio: string;
  city: string;
  roles: string[];
  permissions: string[];
}

export interface CurrentUser extends UserSummary {
  status: number;
  createdAt: string;
  updatedAt: string;
}

export interface BadgeItem {
  id: number;
  name: string;
  description: string;
  icon: string;
  conditionType: string;
  conditionValue: number;
  unlockedAt: string;
}

export interface SessionData {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: UserSummary;
}

export interface RouteWaypoint {
  id?: number;
  name: string;
  lat: number;
  lng: number;
  stayDuration: number;
  order?: number;
}

export interface RouteItem {
  id: number;
  title: string;
  description: string;
  themeId: number;
  city: string;
  startLat: number;
  startLng: number;
  endLat: number;
  endLng: number;
  totalDistance: number;
  duration: number;
  difficulty: number;
  coverImage: string;
  createdBy: number;
  status: number;
  favoriteCount: number;
  rating: number;
  ratingCount: number;
  viewCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface RouteDetail extends RouteItem {
  waypoints: RouteWaypoint[];
}

export interface StoryMedia {
  id?: number;
  storyId?: number;
  type: number;
  url: string;
  order?: number;
}

export interface StoryItem {
  id: number;
  title: string;
  content: string;
  coverImage: string;
  landmarkIds: string;
  routeId: number;
  createdBy: number;
  status: number;
  likeCount: number;
  viewCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface StoryDetail extends StoryItem {
  media: StoryMedia[];
}

export interface LandmarkItem {
  id: number;
  name: string;
  address: string;
  lat: number;
  lng: number;
  categoryId: number;
  description: string;
  coverImage: string;
  createdBy: number;
  status: number;
  rating: number;
  ratingCount: number;
  viewCount: number;
  createdAt: string;
  updatedAt: string;
  distance?: number;
}

export interface EventParticipation {
  id: number;
  eventId: number;
  userId: number;
  name: string;
  remark: string;
  status: number;
  createdAt: string;
  updatedAt: string;
}

export interface EventItem {
  id: number;
  title: string;
  routeId: number;
  meetupAddress: string;
  meetupLat: number;
  meetupLng: number;
  meetupTime: string;
  startTime: string;
  endTime: string;
  maxParticipants: number;
  currentParticipants: number;
  fee: number;
  description: string;
  createdBy: number;
  status: number;
  createdAt: string;
  updatedAt: string;
}

export interface EventDetail extends EventItem {
  participations: EventParticipation[];
}

export interface AdminUser {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  bio: string;
  city: string;
  status: number;
  roles: string[];
  createdAt: string;
  updatedAt: string;
}

export interface RoleItem {
  id: number;
  name: string;
  description: string;
  createdAt: string;
}

export interface PermissionItem {
  id: number;
  name: string;
  resource: string;
  action: string;
  description: string;
  createdAt: string;
}
