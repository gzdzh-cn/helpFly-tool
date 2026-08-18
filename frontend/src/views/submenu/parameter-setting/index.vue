<template>
  <div class="container">
    <a-card class="general-card" :bordered="false" title="系统参数设置">
      <a-form ref="formRef" :model="form" layout="vertical" class="setting-form">
        <a-divider orientation="left">基础信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item field="systemName" label="系统名称" required>
              <a-input v-model="form.systemName" placeholder="请输入系统名称" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="orgCode" label="机构编码">
              <a-input v-model="form.orgCode" placeholder="请输入机构编码" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="adminEmail" label="管理员邮箱" required>
              <a-input v-model="form.adminEmail" placeholder="name@example.com" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="timezone" label="时区">
              <a-select v-model="form.timezone" placeholder="请选择时区" allow-clear>
                <a-option value="Asia/Shanghai"> Asia/Shanghai (UTC+8)</a-option>
                <a-option value="Asia/Tokyo">Asia/Tokyo (UTC+9)</a-option>
                <a-option value="Europe/London">Europe/London (UTC+0)</a-option>
                <a-option value="America/New_York">America/New_York (UTC-5)</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">外观与主题</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item field="themeMode" label="主题模式">
              <a-radio-group v-model="form.themeMode" type="button">
                <a-radio value="light">浅色</a-radio>
                <a-radio value="dark">深色</a-radio>
                <a-radio value="auto">跟随系统</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="accentColor" label="主题色">
              <a-color-picker v-model="form.accentColor" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="language" label="界面语言">
              <a-select v-model="form.language" placeholder="请选择语言">
                <a-option value="zh-CN">简体中文</a-option>
                <a-option value="zh-TW">繁体中文</a-option>
                <a-option value="en-US">English</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="fontSize" label="字体大小">
              <a-slider v-model="form.fontSize" :min="12" :max="20" :step="1" :format-tooltip="(v)=>v+'px'" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">数据同步</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item field="autoSync" label="自动同步">
              <a-switch v-model="form.autoSync" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="syncInterval" label="同步间隔（分钟）">
              <a-input-number v-model="form.syncInterval" :min="1" :max="1440" :disabled="!form.autoSync" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="endpoint" label="同步地址">
              <a-input v-model="form.endpoint" placeholder="https://sync.example.com" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="retainDays" label="数据保留天数">
              <a-input-number v-model="form.retainDays" :min="1" :max="3650" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">通知与安全</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item field="notifyEmail" label="邮件通知">
              <a-switch v-model="form.notifyEmail" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="notifySms" label="短信通知">
              <a-switch v-model="form.notifySms" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="lockOnIdle" label="空闲锁定">
              <a-switch v-model="form.lockOnIdle" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="idleMinutes" label="空闲时间（分钟）">
              <a-input-number v-model="form.idleMinutes" :min="1" :max="120" :disabled="!form.lockOnIdle" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item field="remark" label="备注说明">
              <a-textarea v-model="form.remark" placeholder="可填写参数说明" :auto-size="{ minRows: 2, maxRows: 5 }" allow-clear />
            </a-form-item>
          </a-col>
        </a-row>

        <div class="actions">
          <a-space>
            <a-button @click="handleReset">重置</a-button>
            <a-button type="primary" :loading="saving" @click="handleSave">保存设置</a-button>
          </a-space>
        </div>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
  import { reactive, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useAppStore } from '@/store';

  const appStore = useAppStore();
  const formRef = ref();
  const saving = ref(false);

  // 当前实际主题 -> 表单选项
  const currentThemeToMode = (): 'light' | 'dark' | 'auto' => {
    if (appStore.theme === 'light') return 'light';
    if (appStore.theme === 'dark') return 'dark';
    return 'auto';
  };

  const defaultForm = () => ({
    systemName: 'HelpFly 助手',
    orgCode: 'ORG-0001',
    adminEmail: 'admin@example.com',
    timezone: 'Asia/Shanghai',
    themeMode: currentThemeToMode(),
    accentColor: '#165DFF',
    language: 'zh-CN',
    fontSize: 14,
    autoSync: true,
    syncInterval: 30,
    endpoint: 'https://sync.example.com',
    retainDays: 90,
    notifyEmail: true,
    notifySms: false,
    lockOnIdle: false,
    idleMinutes: 15,
    remark: '',
  });

  const form = reactive(defaultForm());

  // 真实应用主题到全局
  const applyTheme = (mode: 'light' | 'dark' | 'auto') => {
    if (mode === 'light') {
      appStore.toggleTheme(false);
    } else if (mode === 'dark') {
      appStore.toggleTheme(true);
    } else {
      // 跟随系统
      const dark = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false;
      appStore.toggleTheme(dark);
    }
  };

  const handleReset = () => {
    Object.assign(form, defaultForm());
  };

  const handleSave = async () => {
    const err = await formRef.value?.validate();
    if (err) {
      Message.warning('请完善必填项');
      return;
    }
    applyTheme(form.themeMode);
    saving.value = true;
    setTimeout(() => {
      saving.value = false;
      Message.success('参数已保存');
    }, 600);
  };
</script>

<style scoped lang="less">
  .container {
    padding: 16px;
    box-sizing: border-box;
  }
  .general-card {
    background: var(--color-bg-2);
  }
  .setting-form {
    max-width: 960px;
  }
  .actions {
    margin-top: 8px;
    text-align: right;
  }
</style>
