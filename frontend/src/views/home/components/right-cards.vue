<template>
  <div class="right-cards">
    <!-- 系统状态 -->
    <a-card class="rc-card" title="系统状态" :body-style="{ padding: '16px' }">
      <div v-for="item in statusList" :key="item.label" class="status-row">
        <div class="status-top">
          <span class="status-label">{{ item.label }}</span>
          <span class="status-value">{{ item.value }}%</span>
        </div>
        <a-progress
          :percent="item.value"
          :color="item.color"
          :show-text="false"
          size="small"
        />
      </div>
    </a-card>

    <!-- 实时指标 -->
    <a-card class="rc-card" title="实时指标" :body-style="{ padding: '16px' }">
      <a-row :gutter="12">
        <a-col :span="12" v-for="m in metrics" :key="m.label" class="metric-col">
          <div class="metric">
            <div class="metric-value" :style="{ color: m.color }">{{ m.value }}</div>
            <div class="metric-label">{{ m.label }}</div>
          </div>
        </a-col>
      </a-row>
    </a-card>

    <!-- 热门标签 -->
    <a-card class="rc-card" title="热门标签" :body-style="{ padding: '16px' }">
      <div class="tag-wrap">
        <a-tag v-for="t in tags" :key="t.text" :color="t.color" class="tag-item">
          {{ t.text }}
        </a-tag>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  const statusList = [
    { label: 'CPU 使用率', value: 42, color: '#165DFF' },
    { label: '内存占用', value: 68, color: '#00B42A' },
    { label: '磁盘空间', value: 35, color: '#FF7D00' },
    { label: '网络负载', value: 57, color: '#722ED1' },
  ];

  const metrics = [
    { label: '今日访问', value: '1,284', color: '#165DFF' },
    { label: '新增用户', value: '86', color: '#00B42A' },
    { label: '订单数', value: '312', color: '#FF7D00' },
    { label: '转化率', value: '3.8%', color: '#722ED1' },
  ];

  const tags = [
    { text: 'Go', color: 'arcoblue' },
    { text: 'Vue3', color: 'green' },
    { text: 'Wails', color: 'orange' },
    { text: 'ECharts', color: 'purple' },
    { text: 'ArcoDesign', color: 'red' },
    { text: 'TypeScript', color: 'cyan' },
    { text: 'bbolt', color: 'lime' },
    { text: '桌面应用', color: 'magenta' },
    { text: '数据可视化', color: 'arcoblue' },
  ];
</script>

<style scoped lang="less">
  .rc-card {
    background-color: var(--color-neutral-1);
    border-radius: 4px;
    margin-bottom: 10px;
  }
  .status-row {
    margin-bottom: 14px;
    &:last-child {
      margin-bottom: 0;
    }
    .status-top {
      display: flex;
      justify-content: space-between;
      margin-bottom: 6px;
      font-size: 13px;
      .status-label {
        color: var(--color-text-2);
      }
      .status-value {
        color: var(--color-text-1);
        font-weight: 600;
      }
    }
  }
  .metric-col {
    margin-bottom: 12px;
  }
  .metric {
    background: var(--color-fill-1);
    border-radius: 6px;
    padding: 12px;
    text-align: center;
    .metric-value {
      font-size: 20px;
      font-weight: 700;
      line-height: 1.2;
    }
    .metric-label {
      font-size: 12px;
      color: var(--color-text-3);
      margin-top: 2px;
    }
  }
  .tag-wrap {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    .tag-item {
      cursor: default;
    }
  }
</style>
