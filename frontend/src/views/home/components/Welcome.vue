<template>
  <a-card class="card" :bordered="false">
    <a-row  wrap :gutter="[{ xs: 0, sm: 14, md: 14, lg: 14, xl: 14, xxl: 14 }, 16]" class="content">
      <a-space size="medium">
        <a-avatar style="padding: 2px;" :size="68">
          <img src="@/assets/logo.png">
        </a-avatar>
        <div class="welcome">
          <p class="hello">{{ goodTimeText() }}！{{ userStore.nickname }}</p>
          <p>当前时间：<span class="datetime">{{ appStore.sysTime?appStore.sysTime:"---" }}</span></p>
        </div>
      </a-space>
    </a-row>
  </a-card>
</template>

<script lang="ts" setup >
  // import { ref } from 'vue';
  import { useUserStore,useAppStore } from '@/store';
  import { goodTimeText } from '@/utils'
 

  const appStore = useAppStore();
  const userStore = useUserStore()
</script>

<style scoped lang="less">
:deep(.arco-statistic-title) {
  margin-bottom: 0;
}

.card {
  background: transparent;
  .content {
    padding: 8px 20px;
    .welcome {
      margin: 8px 0;
      color: var(--color-neutral-6);
      .hello {
        font-size: 1.25rem;
        color: var(--color-neutral-8);
        margin-bottom: 10px;
      }
    }
  }
}
</style>
