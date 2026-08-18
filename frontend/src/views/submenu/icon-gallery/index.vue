<template>
  <div class="container">
    <a-card class="general-card" title="图标专栏" :body-style="{ padding: '16px' }">
      <a-alert type="info" class="tip">
        内置本地 SVG 图标库（离线可用）+ 在线阿里 iconfont。支持分类筛选与搜索，本地图标分页展示。
      </a-alert>

      <div class="filter-bar">
        <a-radio-group v-model="source" type="button" size="small">
          <a-radio value="local">本地图标</a-radio>
          <a-radio value="online">在线字体</a-radio>
        </a-radio-group>
        <a-input v-model="keyword" placeholder="搜索图标名称" allow-clear style="width: 240px">
          <template #prefix><icon-search /></template>
        </a-input>
      </div>

      <!-- 在线图标 -->
      <template v-if="source === 'online'">
        <a-tabs v-model:active-key="onlineCat">
          <a-tab-pane v-for="c in onlineCats" :key="c" :title="c" />
        </a-tabs>
        <div class="icon-grid">
          <div v-for="ic in filteredOnline" :key="ic" class="icon-cell" @click="copy(ic)">
            <icon-font :type="ic" :size="26" />
            <span class="icon-name">{{ ic }}</span>
          </div>
        </div>
      </template>

      <!-- 本地图标 -->
      <template v-else>
        <a-tabs v-model:active-key="localCat">
          <a-tab-pane key="全部" title="全部" />
          <a-tab-pane v-for="c in iconCategories" :key="c" :title="c" />
        </a-tabs>
        <div class="icon-grid">
          <div v-for="ic in pagedLocal" :key="ic.name" class="icon-cell" @click="copy(ic.name)">
            <LocalIcon :name="ic.name.replace(/-v\d+$/, '')" :size="26" />
            <span class="icon-name">{{ ic.name }}</span>
          </div>
        </div>
        <div class="pager">
          <a-pagination
            v-model:current="page"
            :page-size="pageSize"
            :total="filteredLocal.length"
            show-total
          />
        </div>
      </template>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, computed, watch } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import LocalIcon from '@/components/icons/LocalIcon.vue';
  import { localIcons, iconCategories, expandLocalIcons } from '@/components/icons/local-icons';

  const source = ref<'local' | 'online'>('local');
  const keyword = ref('');
  const page = ref(1);
  const pageSize = ref(100);

  const localCat = ref('全部');
  const onlineCat = ref('常用');

  // 在线 iconfont 演示集合（项目内已使用过的确定图标）
  const onlineIconMap: Record<string, string[]> = {
    常用: ['icon-filled', 'icon-warning', 'icon-task', 'icon-baobiao', 'icon-pdf', 'icon-date', 'icon-shouye1', 'icon-gform1'],
    导航: ['icon-renwu1', 'icon-xiaofeijilu', 'icon-daochuwenjian-', 'icon-yonghushuliang1', 'icon-biaodan2', 'icon-icon_percent', 'icon-canshuguanli'],
    业务: ['icon-tongzhi1', 'icon-wenjian', 'icon-xitongguanli', 'icon-lirunfenxikuangjia', 'icon-dialog', 'icon-tubiao', 'icon-zujian', 'icon-biaoge'],
  };
  const onlineCats = Object.keys(onlineIconMap);
  const filteredOnline = computed(() => {
    const list = onlineIconMap[onlineCat.value] || [];
    if (!keyword.value) return list;
    return list.filter((i) => i.includes(keyword.value));
  });

  // 本地图标：扩展成数百个用于演示
  const allLocal = expandLocalIcons(5);
  const filteredLocal = computed(() => {
    let list = allLocal;
    if (localCat.value !== '全部') list = list.filter((i) => i.category === localCat.value);
    if (keyword.value) list = list.filter((i) => i.name.includes(keyword.value));
    return list;
  });
  const pagedLocal = computed(() => {
    const start = (page.value - 1) * pageSize.value;
    return filteredLocal.value.slice(start, start + pageSize.value);
  });

  watch([source, localCat, keyword], () => {
    page.value = 1;
  });

  function copy(name: string) {
    navigator.clipboard?.writeText(name).then(
      () => Message.success(`已复制图标名：${name}`),
      () => Message.info(name),
    );
  }
</script>

<style scoped lang="less">
  .container {
    padding: 16px;
    height: 100%;
    box-sizing: border-box;
    overflow: auto;
  }
  .general-card {
    background-color: var(--color-bg-2);
    border-radius: 4px;
  }
  .tip {
    margin-bottom: 14px;
  }
  .filter-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 14px;
  }
  .icon-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(96px, 1fr));
    gap: 12px;
  }
  .icon-cell {
    border: 1px solid var(--color-neutral-3);
    border-radius: 6px;
    padding: 14px 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    transition: all 0.2s;
    background: var(--color-bg-1);
    &:hover {
      border-color: rgb(var(--primary-6));
      box-shadow: 0 3px 10px rgba(0, 0, 0, 0.08);
      transform: translateY(-2px);
    }
    .icon-name {
      font-size: 11px;
      color: var(--color-text-3);
      word-break: break-all;
      text-align: center;
      line-height: 1.3;
    }
  }
  .pager {
    display: flex;
    justify-content: flex-end;
    margin-top: 14px;
  }
</style>
