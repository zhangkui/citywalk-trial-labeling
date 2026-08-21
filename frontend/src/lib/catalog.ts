export const themeOptions = [
  { value: 1, label: '城市风貌' },
  { value: 2, label: '历史建筑' },
  { value: 3, label: '美食探店' },
  { value: 4, label: '文艺书店' },
  { value: 5, label: '咖啡地图' },
  { value: 6, label: '街头艺术' },
];

export const landmarkCategories = [
  { value: 1, label: '历史建筑' },
  { value: 2, label: '博物馆' },
  { value: 3, label: '书店' },
  { value: 4, label: '咖啡馆' },
  { value: 5, label: '公园' },
  { value: 6, label: '艺术空间' },
  { value: 7, label: '美食' },
];

export const routeStatusOptions = [
  { value: 0, label: '草稿' },
  { value: 1, label: '待审' },
  { value: 2, label: '已发布' },
  { value: 3, label: '已下架' },
];

export const storyStatusOptions = [
  { value: 0, label: '草稿' },
  { value: 1, label: '待审' },
  { value: 2, label: '已发布' },
];

export const eventStatusOptions = [
  { value: 0, label: '待开始' },
  { value: 1, label: '进行中' },
  { value: 2, label: '已结束' },
  { value: 3, label: '已取消' },
];

export const statusLabel = (items: { value: number; label: string }[], value: number) =>
  items.find((item) => item.value === value)?.label ?? String(value);
