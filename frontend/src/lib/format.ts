export const formatDateTime = (value?: string) => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('zh-CN', { hour12: false });
};

export const formatMoney = (value?: number) => {
  if (value == null) return '-';
  return value === 0 ? '免费' : `¥${(value / 100).toFixed(2)}`;
};

export const formatDistance = (value?: number) => {
  if (!value && value !== 0) return '-';
  return value >= 1000 ? `${(value / 1000).toFixed(1)} km` : `${value} m`;
};

export const clampText = (value: string, max = 100) => {
  if (value.length <= max) return value;
  return `${value.slice(0, max)}...`;
};
