<template>
  <a-card
    class="general-card demostat-cards"
    title="统计演示"
    :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '16px' }"
  >
    <a-row :gutter="12">
      <!-- 柱状图 -->
      <a-col :xs="24" :sm="24" :md="12" :lg="12" class="card-col">
        <div class="stat-box">
          <div class="stat-head">
            <span class="stat-title">月度访问量</span>
            <span class="stat-tag">柱状图</span>
          </div>
          <Chart height="200px" :option="barOption" />
        </div>
      </a-col>
      <!-- 环形占比 -->
      <a-col :xs="24" :sm="24" :md="12" :lg="12" class="card-col">
        <div class="stat-box">
          <div class="stat-head">
            <span class="stat-title">流量来源占比</span>
            <span class="stat-tag">环形图</span>
          </div>
          <Chart height="200px" :option="pieOption" />
        </div>
      </a-col>
      <!-- 折线趋势 -->
      <a-col :xs="24" :sm="24" :md="12" :lg="12" class="card-col">
        <div class="stat-box">
          <div class="stat-head">
            <span class="stat-title">活跃用户趋势</span>
            <span class="stat-tag">折线图</span>
          </div>
          <Chart height="200px" :option="lineOption" />
        </div>
      </a-col>
      <!-- 横向条形 -->
      <a-col :xs="24" :sm="24" :md="12" :lg="12" class="card-col">
        <div class="stat-box">
          <div class="stat-head">
            <span class="stat-title">地区销售排行</span>
            <span class="stat-tag">条形图</span>
          </div>
          <Chart height="200px" :option="hbarOption" />
        </div>
      </a-col>
    </a-row>
  </a-card>
</template>

<script lang="ts" setup>
  import useChartOption from '@/hooks/chart-option';

  const palette = ['#165DFF', '#00B42A', '#FF7D00', '#722ED1', '#F76965'];

  const { chartOption: barOption } = useChartOption((isDark) => ({
    grid: { left: 10, right: 10, top: 20, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月'],
      axisLine: { lineStyle: { color: isDark ? '#3a3a3a' : '#e5e6eb' } },
      axisLabel: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
      axisLabel: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    series: [
      {
        type: 'bar',
        data: [320, 432, 501, 434, 590, 620],
        barWidth: '45%',
        itemStyle: {
          color: palette[0],
          borderRadius: [4, 4, 0, 0],
        },
      },
    ],
  }));

  const { chartOption: pieOption } = useChartOption((isDark) => ({
    tooltip: { trigger: 'item' },
    legend: {
      bottom: 0,
      textStyle: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    series: [
      {
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: false,
        label: { show: false },
        data: [
          { value: 1048, name: '直接访问', itemStyle: { color: palette[0] } },
          { value: 735, name: '搜索引擎', itemStyle: { color: palette[1] } },
          { value: 580, name: '社交媒体', itemStyle: { color: palette[2] } },
          { value: 484, name: '外部链接', itemStyle: { color: palette[3] } },
        ],
      },
    ],
  }));

  const { chartOption: lineOption } = useChartOption((isDark) => ({
    grid: { left: 10, right: 10, top: 20, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: isDark ? '#3a3a3a' : '#e5e6eb' } },
      axisLabel: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
      axisLabel: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    series: [
      {
        type: 'line',
        smooth: true,
        data: [820, 932, 901, 934, 1290, 1330, 1320],
        symbol: 'circle',
        symbolSize: 6,
        itemStyle: { color: palette[1] },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(0,180,42,0.3)' },
              { offset: 1, color: 'rgba(0,180,42,0)' },
            ],
          },
        },
      },
    ],
  }));

  const { chartOption: hbarOption } = useChartOption((isDark) => ({
    grid: { left: 10, right: 30, top: 10, bottom: 10, containLabel: true },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    xAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
      axisLabel: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    yAxis: {
      type: 'category',
      data: ['华南', '华北', '华东', '西南', '东北'],
      axisLine: { lineStyle: { color: isDark ? '#3a3a3a' : '#e5e6eb' } },
      axisLabel: { color: isDark ? '#a9aeb8' : '#86909c' },
    },
    series: [
      {
        type: 'bar',
        data: [
          { value: 320, itemStyle: { color: palette[3] } },
          { value: 410, itemStyle: { color: palette[2] } },
          { value: 520, itemStyle: { color: palette[0] } },
          { value: 260, itemStyle: { color: palette[4] } },
          { value: 180, itemStyle: { color: palette[1] } },
        ],
        barWidth: '55%',
        itemStyle: { borderRadius: [0, 4, 4, 0] },
      },
    ],
  }));
</script>

<style scoped lang="less">
  .general-card {
    background-color: var(--color-neutral-1);
    border-radius: 4px;
  }
  .card-col {
    margin-bottom: 12px;
  }
  .stat-box {
    background: var(--color-bg-2);
    border: 1px solid var(--color-neutral-3);
    border-radius: 6px;
    padding: 14px 16px;
    height: 100%;
  }
  .stat-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }
  .stat-title {
    font-size: 14px;
    font-weight: 600;
    color: var(--color-text-1);
  }
  .stat-tag {
    font-size: 12px;
    color: rgb(var(--primary-6));
    background: var(--color-primary-light-1);
    border-radius: 4px;
    padding: 1px 8px;
  }
</style>
