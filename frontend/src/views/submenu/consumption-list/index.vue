<template>
  <div class="container">
    <a-card class="general-card" :body-style="{ padding: '16px' }">
      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="toolbar-left">
          <a-radio-group :model-value="dataMode" type="button" size="small" @change="onModeChange">
            <a-radio value="backend">后端数据</a-radio>
            <a-radio value="mock">本地演示</a-radio>
          </a-radio-group>
          <a-input
            v-model="keyword"
            placeholder="搜索摘要/渠道/金额"
            allow-clear
            style="width: 220px"
            @press-enter="load"
            @clear="load"
          >
            <template #prefix><icon-search /></template>
          </a-input>
          <a-button @click="openAdvance">高级搜索</a-button>
          <a-button type="primary" @click="openCreate">新增</a-button>
          <a-button @click="exportData">导出</a-button>
          <a-button
            status="danger"
            :disabled="!selectedKeys.length"
            @click="batchDelete"
          >
            批量删除{{ selectedKeys.length ? `(${selectedKeys.length})` : '' }}
          </a-button>
          <a-button @click="load"><icon-refresh /></a-button>
        </div>
        <div class="toolbar-right">
          <a-tag v-if="dataMode === 'backend'" color="green">已连接 Go 后端</a-tag>
          <a-tag v-else color="arcoblue">本地 Mock 数据</a-tag>
        </div>
      </div>

      <!-- 表格 -->
      <a-table
        row-key="key"
        :loading="loading"
        :data="pagedData"
        :pagination="false"
        :scroll="{ x: '100%', y: '100%' }"
        :row-selection="{ type: 'checkbox', showCheckedAll: true, onlyCurrent: false }"
        :selected-keys="selectedKeys"
        @selection-change="onSelect"
      >
        <template #columns>
          <a-table-column title="序号" :width="70">
            <template #cell="{ rowIndex }">{{ (page - 1) * pageSize + rowIndex + 1 }}</template>
          </a-table-column>
          <a-table-column title="交易日期" data-index="tradeDate" :width="180" />
          <a-table-column title="渠道" data-index="channel" :width="120">
            <template #cell="{ record }">
              <a-tag :color="channelColor(record.channel)">{{ record.channel }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="摘要" data-index="summary" :width="160" />
          <a-table-column title="金额(元)" :width="130">
            <template #cell="{ record }">
              <span :style="{ color: record.incomeOrExpenseAmount >= 0 ? '#00B42A' : '#F53F3F', fontWeight: 600 }">
                {{ record.incomeOrExpenseAmount >= 0 ? '+' : '' }}{{ record.incomeOrExpenseAmount.toFixed(2) }}
              </span>
            </template>
          </a-table-column>
          <a-table-column title="余额(元)" data-index="balance" :width="130">
            <template #cell="{ record }">{{ Number(record.balance).toFixed(2) }}</template>
          </a-table-column>
          <a-table-column title="操作" :width="200" fixed="right">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
                <a-popconfirm content="确认删除该记录？" @ok="singleDelete(record)">
                  <a-button type="text" size="small" status="danger">删除</a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <div class="pager">
        <a-pagination
          v-model:current="page"
          v-model:page-size="pageSize"
          :total="filteredData.length"
          show-total
          show-page-size
        />
      </div>
    </a-card>

    <!-- 高级搜索抽屉 -->
    <a-drawer v-model:visible="advVisible" title="高级搜索" :width="360" :footer="false">
      <a-form :model="advForm" layout="vertical">
        <a-form-item label="关键字">
          <a-input v-model="advForm.keyword" placeholder="摘要/渠道/金额" allow-clear />
        </a-form-item>
        <a-form-item label="渠道">
          <a-select v-model="advForm.channel" placeholder="全部渠道" allow-clear>
            <a-option v-for="c in channelOptions" :key="c" :value="c">{{ c }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="收支类型">
          <a-select v-model="advForm.type" placeholder="全部" allow-clear>
            <a-option value="income">收入</a-option>
            <a-option value="expense">支出</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="交易日期">
          <a-range-picker v-model="advForm.dateRange" style="width: 100%" />
        </a-form-item>
        <a-form-item label="金额区间(绝对值)">
          <a-space>
            <a-input-number v-model="advForm.min" placeholder="最小" :min="0" style="width: 100px" />
            <span>~</span>
            <a-input-number v-model="advForm.max" placeholder="最大" :min="0" style="width: 100px" />
          </a-space>
        </a-form-item>
        <a-space>
          <a-button type="primary" @click="applyAdvance">搜索</a-button>
          <a-button @click="resetAdvance">重置</a-button>
        </a-space>
      </a-form>
    </a-drawer>

    <!-- 新增/编辑/详情 抽屉 -->
    <a-drawer
      v-model:visible="formVisible"
      :title="drawerTitle"
      :width="420"
      @cancel="formVisible = false"
      @before-ok="submitForm"
    >
      <a-form v-if="mode !== 'detail'" :model="form" layout="vertical" ref="formRef">
        <a-form-item field="tradeDate" label="交易日期" :rules="[{ required: true, message: '请选择交易日期' }]">
          <a-date-picker v-model="form.tradeDate" show-time style="width: 100%" />
        </a-form-item>
        <a-form-item field="channel" label="渠道" :rules="[{ required: true, message: '请选择渠道' }]">
          <a-select v-model="form.channel" placeholder="请选择渠道">
            <a-option v-for="c in channelOptions" :key="c" :value="c">{{ c }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="summary" label="摘要" :rules="[{ required: true, message: '请输入摘要' }]">
          <a-input v-model="form.summary" placeholder="如：超市购物" />
        </a-form-item>
        <a-form-item field="incomeOrExpenseAmount" label="金额(元，正为收入/负为支出)" :rules="[{ required: true, message: '请输入金额' }]">
          <a-input-number v-model="form.incomeOrExpenseAmount" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item field="balance" label="余额(元)">
          <a-input-number v-model="form.balance" :precision="2" style="width: 100%" />
        </a-form-item>
      </a-form>
      <a-descriptions v-else :data="detailData" :column="1" bordered />
    </a-drawer>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, computed, onMounted } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { DbService } from '/#/helpfly/internal/service';
  import {
    fetchList,
    addRecord,
    updateRecord,
    deleteRecords,
    applyAdvancedFilter,
    type DataMode,
  } from './api';
  import { channelOptions, type ConsumptionMock } from './mock';

  const dataMode = ref<DataMode>('backend');
  const keyword = ref('');
  const loading = ref(false);
  const allData = ref<ConsumptionMock[]>([]);
  const selectedKeys = ref<string[]>([]);

  const page = ref(1);
  const pageSize = ref(10);

  const advVisible = ref(false);
  const advForm = reactive({
    keyword: '',
    channel: '',
    type: '',
    dateRange: [] as string[],
    min: undefined as number | undefined,
    max: undefined as number | undefined,
  });

  const formVisible = ref(false);
  const mode = ref<'create' | 'edit' | 'detail'>('create');
  const formRef = ref();
  const form = reactive<ConsumptionMock>({
    key: '',
    taskId: '',
    tradeDate: '',
    channel: '',
    summary: '',
    incomeOrExpenseAmount: 0,
    balance: 0,
  });

  const drawerTitle = computed(() =>
    mode.value === 'create' ? '新增账单' : mode.value === 'edit' ? '编辑账单' : '账单详情',
  );

  const channelColor = (c: string) => {
    const map: Record<string, string> = {
      微信: 'green',
      支付宝: 'blue',
      银行卡: 'purple',
      现金: 'orange',
      云闪付: 'cyan',
    };
    return map[c] || 'gray';
  };

  // 先按关键字过滤，再按高级条件过滤
  const filteredData = computed(() => {
    let list = allData.value;
    if (keyword.value) list = list.filter((i) => JSON.stringify(i).includes(keyword.value));
    return applyAdvancedFilter(list, {
      keyword: advForm.keyword,
      channel: advForm.channel,
      type: advForm.type,
      dateRange: advForm.dateRange,
      min: advForm.min,
      max: advForm.max,
    });
  });

  const pagedData = computed(() => {
    const start = (page.value - 1) * pageSize.value;
    return filteredData.value.slice(start, start + pageSize.value);
  });

  async function load() {
    loading.value = true;
    try {
      allData.value = await fetchList(dataMode.value, keyword.value);
    } catch (e) {
      Message.error('后端连接失败，已自动切换到本地演示');
      dataMode.value = 'mock';
      allData.value = await fetchList('mock', keyword.value);
    } finally {
      loading.value = false;
    }
  }

  function onModeChange(v: string | number | boolean) {
    dataMode.value = v as DataMode;
    page.value = 1;
    load();
  }

  async function singleDelete(record: ConsumptionMock) {
    await deleteRecords(dataMode.value, [record.key]);
    Message.success('删除成功');
    load();
  }

  async function batchDelete() {
    if (!selectedKeys.value.length) return;
    await deleteRecords(dataMode.value, selectedKeys.value);
    Message.success('批量删除成功');
    selectedKeys.value = [];
    load();
  }

  function onSelect(keys: string[]) {
    selectedKeys.value = keys;
  }

  function openAdvance() {
    advVisible.value = true;
  }
  function applyAdvance() {
    advVisible.value = false;
    page.value = 1;
  }
  function resetAdvance() {
    Object.assign(advForm, { keyword: '', channel: '', type: '', dateRange: [], min: undefined, max: undefined });
  }

  function openCreate() {
    mode.value = 'create';
    Object.assign(form, {
      key: '',
      taskId: '',
      tradeDate: '',
      channel: '',
      summary: '',
      incomeOrExpenseAmount: 0,
      balance: 0,
    });
    formVisible.value = true;
  }
  function openEdit(record: ConsumptionMock) {
    mode.value = 'edit';
    Object.assign(form, record);
    formVisible.value = true;
  }
  const detailData = computed(() => [
    { label: '交易日期', value: form.tradeDate },
    { label: '渠道', value: form.channel },
    { label: '摘要', value: form.summary },
    { label: '金额(元)', value: `${form.incomeOrExpenseAmount.toFixed(2)}` },
    { label: '余额(元)', value: `${Number(form.balance).toFixed(2)}` },
  ]);
  function openDetail(record: ConsumptionMock) {
    mode.value = 'detail';
    Object.assign(form, record);
    formVisible.value = true;
  }

  async function submitForm() {
    if (mode.value === 'detail') return true;
    const ok = await formRef.value?.validate();
    if (ok) return false;
    if (mode.value === 'create') {
      await addRecord(dataMode.value, { ...form });
      Message.success('新增成功');
    } else {
      await updateRecord(dataMode.value, { ...form });
      Message.success('编辑成功');
    }
    formVisible.value = false;
    load();
    return true;
  }

  function exportData() {
    // 导出当前筛选结果为 CSV（前端纯演示，不依赖后端导出能力）
    const rows = filteredData.value;
    const header = ['交易日期', '渠道', '摘要', '金额', '余额'];
    const body = rows.map((r) => [r.tradeDate, r.channel, r.summary, r.incomeOrExpenseAmount, r.balance]);
    const csv = [header, ...body].map((r) => r.join(',')).join('\n');
    const blob = new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = '账单导出.csv';
    a.click();
    URL.revokeObjectURL(url);
    Message.success(`已导出 ${rows.length} 条记录`);
  }

  onMounted(load);
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
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 14px;
    .toolbar-left {
      display: flex;
      align-items: center;
      flex-wrap: nowrap;
      gap: 8px;
      :deep(.arco-btn),
      :deep(.arco-select),
      :deep(.arco-switch) {
        flex-shrink: 0;
      }
      :deep(.arco-input-wrapper) {
        flex: 1;
        min-width: 160px;
      }
    }
  }
  .pager {
    display: flex;
    justify-content: flex-end;
    margin-top: 14px;
  }
</style>
