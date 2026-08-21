<template>
  <div ref="el" class="map-box"></div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';
import type { LatLngExpression } from 'leaflet';

type Marker = { lat: number; lng: number; label?: string; color?: string };

const props = defineProps<{
  center: [number, number];
  zoom?: number;
  markers?: Marker[];
  path?: [number, number][];
}>();

const el = ref<HTMLDivElement | null>(null);
let map: any = null;

async function renderMap() {
  if (!el.value) return;
  const L = await import('leaflet');
  if (!map) {
    map = L.map(el.value, { zoomControl: true }).setView(props.center as LatLngExpression, props.zoom ?? 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(map);
  }
  map.setView(props.center as LatLngExpression, props.zoom ?? 13);
  map.eachLayer((layer: any) => {
    if (layer.options?.pane === 'overlayPane' || layer instanceof L.Polyline || layer instanceof L.CircleMarker) {
      map.removeLayer(layer);
    }
  });
  (props.markers ?? []).forEach((marker) => {
    L.circleMarker([marker.lat, marker.lng], {
      radius: 7,
      color: marker.color ?? '#f4a261',
      fillColor: marker.color ?? '#f4a261',
      fillOpacity: 0.9,
      weight: 2,
    })
      .addTo(map)
      .bindPopup(marker.label ?? '标记点');
  });
  if ((props.path ?? []).length > 1) {
    L.polyline(props.path as LatLngExpression[], { color: '#6ee7b7', weight: 4, opacity: 0.9 }).addTo(map);
  }
}

onMounted(renderMap);
watch(() => [props.center, props.markers, props.path], renderMap, { deep: true });
onBeforeUnmount(() => {
  if (map) map.remove();
  map = null;
});
</script>
