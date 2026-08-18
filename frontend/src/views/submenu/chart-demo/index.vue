<template>
  <div class="container">
    <a-card class="general-card" title="图表分析" :body-style="{ padding: '16px' }">
      <a-alert type="info" class="tip">基于 ECharts 的可视化专题演示，全部使用本地 mock 数据，亮/暗主题自动适配。</a-alert>
      <a-row :gutter="16">
        <a-col v-for="c in charts" :key="c.title" :xs="24" :sm="24" :md="12" class="card-col">
          <a-card :title="c.title" :body-style="{ padding: '8px' }" class="chart-card">
            <Chart height="260px" :option="c.option" />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import useChartOption from '@/hooks/chart-option';

  const palette = ['#165DFF', '#00B42A', '#FF7D00', '#722ED1', '#F53F3F', '#14C9C9'];

  // 柱状图
  const { chartOption: barOption } = useChartOption((isDark) => ({
    grid: { left: 10, right: 10, top: 20, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: ['一月', '二月', '三月', '四月', '五月', '六月'],
      axisLine: { lineStyle: { color: isDark ? '#3a3a3a' : '#e5e6eb' } },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
    },
    series: [{ type: 'bar', data: [320, 432, 501, 434, 590, 620], barWidth: '45%', itemStyle: { color: palette[0], borderRadius: [4, 4, 0, 0] } }],
  }));

  // 折线图
  const { chartOption: lineOption } = useChartOption((isDark) => ({
    grid: { left: 10, right: 10, top: 20, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: isDark ? '#3a3a3a' : '#e5e6eb' } },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
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
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(0,180,42,0.3)' },
              { offset: 1, color: 'rgba(0,180,42,0)' },
            ],
          },
        },
      },
    ],
  }));

  // 饼图
  const { chartOption: pieOption } = useChartOption(() => ({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: '55%',
        center: ['50%', '45%'],
        data: palette.map((c, i) => ({ value: [320, 280, 210, 160, 120, 90][i], name: ['搜索', '直接', '社交', '广告', '邮件', '其他'][i], itemStyle: { color: c } })),
      },
    ],
  }));

  // 环形图
  const { chartOption: ringOption } = useChartOption(() => ({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
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

  // 散点图
  const { chartOption: scatterOption } = useChartOption((isDark) => ({
    grid: { left: 10, right: 10, top: 20, bottom: 20, containLabel: true },
    tooltip: { trigger: 'item' },
    xAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
      name: 'X',
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: isDark ? '#2a2a2a' : '#f2f3f5' } },
      name: 'Y',
    },
    series: [
      {
        type: 'scatter',
        symbolSize: 14,
        itemStyle: { color: palette[4] },
        data: Array.from({ length: 30 }, () => [Math.round(Math.random() * 100), Math.round(Math.random() * 100)]),
      },
    ],
  }));

  // 雷达图
  const { chartOption: radarOption } = useChartOption(() => ({
    tooltip: {},
    legend: { bottom: 0 },
    radar: {
      indicator: [
        { name: '性能', max: 100 },
        { name: '功能', max: 100 },
        { name: '易用', max: 100 },
        { name: '稳定', max: 100 },
        { name: '安全', max: 100 },
      ],
    },
    series: [
      {
        type: 'radar',
        data: [
          { value: [90, 80, 85, 70, 95], name: '产品A', itemStyle: { color: palette[0] } },
          { value: [70, 95, 75, 90, 80], name: '产品B', itemStyle: { color: palette[2] } },
        ],
      },
    ],
  }));

  // 仪表盘
  const { chartOption: gaugeOption } = useChartOption(() => ({
    series: [
      {
        type: 'gauge',
        progress: { show: true, width: 14 },
        axisLine: { lineStyle: { width: 14 } },
        detail: { valueAnimation: true, formatter: '{value}%', fontSize: 22 },
        data: [{ value: 76, name: '完成率' }],
        itemStyle: { color: palette[1] },
      },
    ],
  }));

  const charts = [
    { title: '柱状图', option: barOption.value },
    { title: '折线图', option: lineOption.value },
    { title: '饼图', option: pieOption.value },
    { title: '环形图', option: ringOption.value },
    { title: '散点图', option: scatterOption.value },
    { title: '雷达图', option: radarOption.value },
    { title: '仪表盘', option: gaugeOption.value },
    { title: '指标概览', option: gaugeOption.value },
  ];
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
  .card-col {
    margin-bottom: 16px;
  }
  .chart-card {
    background: var(--color-bg-1);
    border-radius: 4px;
  }
</style>
