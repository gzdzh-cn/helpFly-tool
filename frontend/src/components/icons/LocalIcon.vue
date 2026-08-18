<template>
  <svg
    class="local-icon"
    viewBox="0 0 24 24"
    :width="size"
    :height="size"
    :style="{ color: color }"
    fill="currentColor"
    aria-hidden="true"
  >
    <path :d="path" />
  </svg>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import { localIcons } from './local-icons';

  const props = withDefaults(
    defineProps<{
      name: string;
      size?: number | string;
      color?: string;
    }>(),
    { size: 18, color: 'currentColor' },
  );

  const path = computed(() => {
    const found = localIcons.find((i) => i.name === props.name);
    return found ? found.path : '';
  });
</script>

<style scoped lang="less">
  .local-icon {
    display: inline-block;
    vertical-align: middle;
  }
</style>
