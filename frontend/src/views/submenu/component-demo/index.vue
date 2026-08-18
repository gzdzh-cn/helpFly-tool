<template>
  <div class="container">
    <a-card class="general-card" title="通用组件演示" :body-style="{ padding: '16px' }">
      <a-alert type="info" class="tip">ArcoDesign 常用组件大杂烩，可直接对照复制用于业务页面。</a-alert>

      <a-row :gutter="16">
        <!-- Tabs -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="Tabs 标签页" :body-style="{ padding: '14px' }">
            <a-tabs default-active-key="1">
              <a-tab-pane key="1" title="基础信息">
                <p class="ph">基础信息内容区，可放置表单或描述。</p>
              </a-tab-pane>
              <a-tab-pane key="2" title="高级配置">
                <p class="ph">高级配置内容区。</p>
              </a-tab-pane>
              <a-tab-pane key="3" title="关联数据">
                <p class="ph">关联数据内容区。</p>
              </a-tab-pane>
            </a-tabs>
          </a-card>
        </a-col>

        <!-- Steps -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="Steps 步骤条" :body-style="{ padding: '14px' }">
            <a-steps :current="1" type="arrow">
              <a-step title="创建" description="填写基本信息" />
              <a-step title="审核" description="等待管理员审核" />
              <a-step title="完成" description="流程结束" />
            </a-steps>
          </a-card>
        </a-col>

        <!-- Timeline -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="Timeline 时间线" :body-style="{ padding: '14px' }">
            <a-timeline>
              <a-timeline-item label="09:00">系统启动完成</a-timeline-item>
              <a-timeline-item label="10:30" dot="success">数据同步成功</a-timeline-item>
              <a-timeline-item label="14:20" dot="warning">磁盘使用率偏高</a-timeline-item>
              <a-timeline-item label="18:00" dot="error">出现一次异常</a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>

        <!-- Tags & Progress -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="Tag 标签 / 进度" :body-style="{ padding: '14px' }">
            <div class="row">
              <a-tag color="red">紧急</a-tag>
              <a-tag color="orange">待处理</a-tag>
              <a-tag color="green">已完成</a-tag>
              <a-tag color="blue">进行中</a-tag>
              <a-tag color="gray">已归档</a-tag>
            </div>
            <div class="row">
              <a-progress :percent="30" />
            </div>
            <div class="row">
              <a-progress :percent="68" :stroke-color="{ '0%': '#165DFF', '100%': '#00B42A' }" />
            </div>
            <div class="row">
              <a-progress type="circle" :percent="85" />
              <a-progress type="dashboard" :percent="60" style="margin-left: 16px" />
            </div>
          </a-card>
        </a-col>

        <!-- Tree -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="Tree 树形" :body-style="{ padding: '14px' }">
            <a-tree :data="treeData" default-expand-all />
          </a-card>
        </a-col>

        <!-- Buttons & Empty & Result -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="按钮 / 空态 / 反馈" :body-style="{ padding: '14px' }">
            <div class="row">
              <a-button type="primary">主要按钮</a-button>
              <a-button>次要按钮</a-button>
              <a-button status="success">成功</a-button>
              <a-button status="warning">警告</a-button>
              <a-button status="danger">危险</a-button>
            </div>
            <div class="row">
              <a-button @click="msg">Message 提示</a-button>
              <a-button @click="noti">Notification 通知</a-button>
              <a-button @click="modal">Modal 弹窗</a-button>
            </div>
            <div class="row">
              <a-empty description="暂无数据" />
            </div>
          </a-card>
        </a-col>

        <!-- Badge & Avatar -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="徽标 / 头像 / 骨架" :body-style="{ padding: '14px' }">
            <div class="row">
              <a-badge :count="9" :dot="false"><a-avatar>U</a-avatar></a-badge>
              <a-badge dot><a-avatar :style="{ backgroundColor: '#165DFF' }">A</a-avatar></a-badge>
              <a-avatar :size="32" :style="{ backgroundColor: '#00B42A' }">B</a-avatar>
              <a-avatar :size="32" :style="{ backgroundColor: '#FF7D00' }">C</a-avatar>
            </div>
            <div class="row">
              <a-skeleton :loading="true" :animation="true">
                <a-skeleton-line :rows="3" />
              </a-skeleton>
            </div>
          </a-card>
        </a-col>

        <!-- Collapse -->
        <a-col :xs="24" :md="12" class="card-col">
          <a-card title="Collapse 折叠面板" :body-style="{ padding: '14px' }">
            <a-collapse :default-active-key="['1']">
              <a-collapse-item key="1" header="常见问题一">
                <p class="ph">这里是问题一的详细解答内容。</p>
              </a-collapse-item>
              <a-collapse-item key="2" header="常见问题二">
                <p class="ph">这里是问题二的详细解答内容。</p>
              </a-collapse-item>
            </a-collapse>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { Message, Notification, Modal } from '@arco-design/web-vue';

  const treeData = [
    {
      title: '根节点',
      key: 'root',
      children: [
        { title: '子节点 A', key: 'a' },
        {
          title: '子节点 B',
          key: 'b',
          children: [
            { title: '叶子 B-1', key: 'b1' },
            { title: '叶子 B-2', key: 'b2' },
          ],
        },
      ],
    },
  ];

  const msg = () => Message.info('这是一条 Message 提示');
  const noti = () =>
    Notification.info({ title: '通知标题', content: '这是右下角通知内容', position: 'bottomRight' });
  const modal = () =>
    Modal.info({ title: '弹窗标题', content: '这是一个信息弹窗，可用于确认或提示。' });
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
  .row {
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
  }
  .ph {
    color: var(--color-text-3);
    font-size: 13px;
    margin: 4px 0;
  }
</style>
