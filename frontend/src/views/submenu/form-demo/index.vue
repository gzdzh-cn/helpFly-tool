<template>
  <div class="container">
    <a-card class="general-card" title="表单 / 详情演示" :body-style="{ padding: '16px' }">
      <a-alert type="info" class="tip">
        本页面为纯前端脚手架演示，涵盖 ArcoDesign 常用输入控件、表单校验，以及新增 / 编辑 / 详情三种模式。
      </a-alert>

      <!-- 模式切换 -->
      <div class="mode-bar">
        <a-radio-group v-model="mode" type="button" size="small">
          <a-radio value="create">新增模式</a-radio>
          <a-radio value="edit">编辑模式</a-radio>
          <a-radio value="detail">详情模式</a-radio>
        </a-radio-group>
        <a-button type="primary" @click="submit">提交保存</a-button>
      </div>

      <!-- 详情模式只读展示 -->
      <a-descriptions v-if="mode === 'detail'" :data="detailData" :column="2" bordered class="form-area" />

      <!-- 新增 / 编辑模式表单 -->
      <a-form v-else ref="formRef" :model="form" layout="vertical" class="form-area">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="name" label="姓名" :rules="[{ required: true, message: '请输入姓名' }]">
              <a-input v-model="form.name" placeholder="请输入姓名" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="category" label="所属类别" :rules="[{ required: true, message: '请选择类别' }]">
              <a-select v-model="form.category" placeholder="请选择" allow-clear>
                <a-option v-for="c in categoryOptions" :key="c" :value="c">{{ c }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="email" label="邮箱" :rules="[{ type: 'email', message: '邮箱格式不正确' }]">
              <a-input v-model="form.email" placeholder="name@example.com" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="phone" label="手机号" :rules="[{ match: /^1\d{10}$/, message: '手机号格式不正确' }]">
              <a-input v-model="form.phone" placeholder="11 位手机号" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="gender" label="性别">
              <a-radio-group v-model="form.gender">
                <a-radio value="male">男</a-radio>
                <a-radio value="female">女</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="birthday" label="出生日期">
              <a-date-picker v-model="form.birthday" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="amount" label="金额(元)">
              <a-input-number v-model="form.amount" :precision="2" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="score" label="评分">
              <a-rate v-model="form.score" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="switchOn" label="启用状态">
              <a-switch v-model="form.switchOn" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="tags" label="标签">
              <a-select v-model="form.tags" multiple placeholder="可多选" allow-clear>
                <a-option v-for="t in tagOptions" :key="t" :value="t">{{ t }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="remark" label="备注">
              <a-textarea v-model="form.remark" placeholder="请输入备注信息" :auto-size="{ minRows: 3, maxRows: 6 }" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, computed } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { categoryOptions, tagOptions, emptyForm, type DemoForm } from './mock';

  const mode = ref<'create' | 'edit' | 'detail'>('create');
  const formRef = ref();
  const form = reactive<DemoForm>({ ...emptyForm(), id: 'D1001' });

  const detailData = computed(() => [
    { label: '编号', value: form.id },
    { label: '姓名', value: form.name || '—' },
    { label: '类别', value: form.category || '—' },
    { label: '邮箱', value: form.email || '—' },
    { label: '手机号', value: form.phone || '—' },
    { label: '性别', value: form.gender === 'male' ? '男' : form.gender === 'female' ? '女' : '—' },
    { label: '出生日期', value: form.birthday || '—' },
    { label: '金额(元)', value: form.amount.toFixed(2) },
    { label: '评分', value: `${form.score} 星` },
    { label: '启用状态', value: form.switchOn ? '已启用' : '已停用' },
    { label: '标签', value: form.tags.join('、') || '—' },
    { label: '备注', value: form.remark || '—' },
  ]);

  async function submit() {
    if (mode.value === 'detail') {
      Message.info('当前为详情模式，不可编辑');
      return;
    }
    const ok = await formRef.value?.validate();
    if (ok) {
      Message.error('请完善表单必填项');
      return;
    }
    Message.success(mode.value === 'create' ? '新增成功（演示）' : '保存成功（演示）');
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
  .mode-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }
  .form-area {
    margin-top: 4px;
  }
</style>
