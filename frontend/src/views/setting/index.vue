<template>
  <div class="container">
    <div class="left-side">
      <div class="card">
        <a-form :model="formData" layout="vertical" size="large">
          <a-divider orientation="left">外观与窗口</a-divider>
          <a-row :gutter="16">
            <a-col :span="16">
              <a-form-item field="theme" label="主题">
                <a-radio-group v-model="formData.theme" type="button" @change="handleTheme">
                  <a-radio value="light">亮色</a-radio>
                  <a-radio value="dark">黑色</a-radio>
                  <a-radio value="auto">跟随系统</a-radio>
                </a-radio-group>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item field="SetAlwaysOnTop" label="置顶窗口">
                <a-switch type="round" @change="changeSetAlwaysOnTop">
                  <template #checked>开启</template>
                  <template #unchecked>关闭</template>
                </a-switch>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="SetFrameless" label="无边框窗口">
                <a-switch type="round" @change="changeSetFrameless">
                  <template #checked>开启</template>
                  <template #unchecked>关闭</template>
                </a-switch>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="SetResizable" label="窗口可调整大小">
                <a-switch type="round" :default-checked="true" @change="changeSetResizable">
                  <template #checked>可以</template>
                  <template #unchecked>不可</template>
                </a-switch>
              </a-form-item>
            </a-col>
          </a-row>

          <a-divider orientation="left">调试参数 · 缩放与布局</a-divider>
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="页面缩放">
                <a-space>
                  <a-button size="small" @click="()=>Window.ZoomOut()">
                    <template #icon><icon-minus /></template>
                  </a-button>
                  <a-button size="small" @click="()=>Window.ZoomReset()">重置</a-button>
                  <a-button size="small" @click="()=>Window.ZoomIn()">
                    <template #icon><icon-plus /></template>
                  </a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="全屏 / 最大化">
                <a-space>
                  <a-button size="small" @click="()=>Window.ToggleFullscreen()">全屏</a-button>
                  <a-button size="small" @click="()=>Window.ToggleMaximise()">最大化</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="窗口显示">
                <a-space>
                  <a-button size="small" @click="()=>Window.Hide()">隐藏</a-button>
                  <a-button size="small" @click="()=>Window.Show()">显示</a-button>
                  <a-button size="small" @click="()=>Window.Restore()">还原</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="title" label="窗口标题">
                <a-input v-model="titleText" placeholder="输入新窗口标题" allow-clear @press-enter="applyTitle">
                  <template #append>
                    <a-button size="mini" type="primary" @click="applyTitle">应用</a-button>
                  </template>
                </a-input>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="背景颜色 (RGBA)">
                <a-space>
                  <a-input-number v-model="bgColor.r" :min="0" :max="255" placeholder="R" />
                  <a-input-number v-model="bgColor.g" :min="0" :max="255" placeholder="G" />
                  <a-input-number v-model="bgColor.b" :min="0" :max="255" placeholder="B" />
                  <a-input-number v-model="bgColor.a" :min="0" :max="255" placeholder="A" />
                  <a-button size="mini" type="primary" @click="applyBgColor">应用</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="最小尺寸 (宽×高)">
                <a-space>
                  <a-input-number v-model="minSize.w" :min="200" :max="4000" placeholder="宽" />
                  <a-input-number v-model="minSize.h" :min="200" :max="4000" placeholder="高" />
                  <a-button size="mini" type="primary" @click="applyMinSize">应用</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="最大尺寸 (宽×高)">
                <a-space>
                  <a-input-number v-model="maxSize.w" :min="200" :max="8000" placeholder="宽" />
                  <a-input-number v-model="maxSize.h" :min="200" :max="8000" placeholder="高" />
                  <a-button size="mini" type="primary" @click="applyMaxSize">应用</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="打印页面">
                <a-button type="primary" size="small" @click="()=>Window.Print()">立即打印</a-button>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="重置窗口">
                <a-button type="primary" size="small" @click="()=>Window.Reload()">
                  <template #icon><icon-refresh /></template>立即重置
                </a-button>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="调试工具">
                <a-button type="primary" size="small" @click="()=>Window.OpenDevTools()">
                  <template #icon><icon-launch /></template>立即打开
                </a-button>
              </a-form-item>
            </a-col>
          </a-row>

          <a-divider orientation="left">调试参数 · 信息探测</a-divider>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="剪贴板测试">
                <a-space>
                  <a-input v-model="clipboardText" placeholder="写入内容" allow-clear />
                  <a-button size="mini" @click="setClipboard">写入</a-button>
                  <a-button size="mini" @click="readClipboard">读取</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="打开外部链接">
                <a-space>
                  <a-input v-model="urlText" placeholder="https://..." allow-clear />
                  <a-button size="mini" type="primary" @click="openUrl">打开</a-button>
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item label="实时状态">
                <a-space wrap>
                  <a-tag :color="focused ? 'green' : 'gray'">焦点: {{ focused ? '已聚焦' : '未聚焦' }}</a-tag>
                  <a-tag :color="fullscreen ? 'green' : 'gray'">全屏: {{ fullscreen ? '是' : '否' }}</a-tag>
                  <a-tag :color="maximised ? 'green' : 'gray'">最大化: {{ maximised ? '是' : '否' }}</a-tag>
                  <a-tag :color="minimised ? 'green' : 'gray'">最小化: {{ minimised ? '是' : '否' }}</a-tag>
                  <a-tag>缩放: {{ zoom }}%</a-tag>
                  <a-button size="mini" @click="refreshState">刷新状态</a-button>
                </a-space>
              </a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </div>
    </div>

    <div class="right-side">
      <div class="card" style="--wails-draggable:drag">
        <a-descriptions title="系统信息（无边框时拖动这里）" :column="1">
          <a-descriptions-item label="窗口name">{{ windowsname }}</a-descriptions-item>
          <a-descriptions-item label="系统(OS)">{{ Osinfo?.OS }}</a-descriptions-item>
          <a-descriptions-item label="系统品牌">{{ Osinfo?.OSInfo.Branding }}</a-descriptions-item>
          <a-descriptions-item label="系统版本">{{ Osinfo?.OSInfo.Version }}</a-descriptions-item>
          <a-descriptions-item label="Arch">{{ Osinfo?.Arch }}</a-descriptions-item>
          <a-descriptions-item label="调试状态">{{ Osinfo?.Debug }}</a-descriptions-item>
          <a-descriptions-item label="WebView2">{{ Osinfo?.PlatformInfo.WebView2 }}</a-descriptions-item>
        </a-descriptions>
      </div>
      <div class="card" style="margin-top:16px">
        <a-descriptions title="屏幕与运行信息" :column="1">
          <a-descriptions-item label="主屏分辨率">{{ screen?.Width }} × {{ screen?.Height }}</a-descriptions-item>
          <a-descriptions-item label="主屏缩放">{{ screen?.ScaleFactor }}</a-descriptions-item>
          <a-descriptions-item label="移动端">{{ isMobile ? '是' : '否' }}</a-descriptions-item>
          <a-descriptions-item label="桌面端">{{ isDesktop ? '是' : '否' }}</a-descriptions-item>
        </a-descriptions>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, computed, onMounted, reactive } from 'vue';
  import { useAppStore } from '@/store';
  import { System, Window, Screens, Clipboard, Browser } from '@wailsio/runtime';
  import { GreetService } from '/#/helpfly/internal/service';
  import { Message } from '@arco-design/web-vue';

  const appStore = useAppStore();
  const theme = computed(() => appStore.theme);
  const formData = ref({ theme: theme.value });

  const windowsname = ref('');
  const IsWindows = ref(false);
  const Osinfo = ref<any>();
  const screen = ref<any>();
  const isMobile = ref(false);
  const isDesktop = ref(false);

  // 实时窗口状态
  const focused = ref(false);
  const fullscreen = ref(false);
  const maximised = ref(false);
  const minimised = ref(false);
  const zoom = ref(100);

  // 标题 / 背景色 / 尺寸
  const titleText = ref('');
  const bgColor = reactive({ r: 22, g: 22, b: 22, a: 255 });
  const minSize = reactive({ w: 1100, h: 720 });
  const maxSize = reactive({ w: 4000, h: 3000 });

  // 剪贴板 / URL
  const clipboardText = ref('');
  const urlText = ref('https://wails.io');

  const refreshState = async () => {
    try {
      focused.value = await Window.IsFocused();
      fullscreen.value = await Window.IsFullscreen();
      maximised.value = await Window.IsMaximised();
      minimised.value = await Window.IsMinimised();
      zoom.value = Math.round((await Window.GetZoom()) * 100);
    } catch (e) {
      // 部分平台可能不支持
    }
  };

  onMounted(async () => {
    windowsname.value = await Window.Name();
    IsWindows.value = await System.IsWindows();
    Osinfo.value = await System.Environment();
    screen.value = await Screens.GetPrimary();
    isMobile.value = await System.IsMobile();
    isDesktop.value = await System.IsDesktop();
    titleText.value = windowsname.value;
    await refreshState();
  });

  // 切换主题
  const handleTheme = async () => {
    let isdark = false;
    if (formData.value.theme === 'dark') {
      isdark = true;
    } else if (formData.value.theme === 'auto') {
      isdark = await System.IsDarkMode();
    }
    appStore.toggleTheme(isdark);
    if (isdark) GreetService.SetTheme();
  };

  const changeSetAlwaysOnTop = (value: any) => {
    Window.SetAlwaysOnTop(value);
    Message.success({ content: value ? '窗口已位于顶部' : '窗口已取消置顶', id: 'setting' });
  };
  const changeSetFrameless = (value: any) => {
    Window.SetFrameless(value);
    Message.success({ content: value ? '窗口已无边框' : '窗口已取消无边框', id: 'setting' });
  };
  const changeSetResizable = (value: any) => {
    Window.SetResizable(value);
    Message.success({ content: value ? '窗口可调整大小' : '窗口不可调整大小', id: 'setting' });
  };

  const applyTitle = () => {
    if (titleText.value) Window.SetTitle(titleText.value);
  };
  const applyBgColor = () => {
    Window.SetBackgroundColour(bgColor.r, bgColor.g, bgColor.b, bgColor.a);
  };
  const applyMinSize = () => Window.SetMinSize(minSize.w, minSize.h);
  const applyMaxSize = () => Window.SetMaxSize(maxSize.w, maxSize.h);

  const setClipboard = () => {
    if (clipboardText.value) {
      Clipboard.SetText(clipboardText.value);
      Message.success('已写入剪贴板');
    }
  };
  const readClipboard = async () => {
    clipboardText.value = await Clipboard.Text();
    Message.success('已读取剪贴板');
  };
  const openUrl = () => {
    if (urlText.value) Browser.OpenURL(urlText.value);
  };
</script>

<style scoped lang="less">
  .container {
    padding: 8px;
    height: 100%;
    box-sizing: border-box;
    min-height: 0;
    overflow: auto;
    display: flex;
    .left-side {
      flex: 1;
      overflow: auto;
    }
    .right-side {
      width: 280px;
      margin-left: 16px;
    }
    .card {
      padding: 8px;
      border-radius: 5px;
      background: var(--color-bg-1);
      box-shadow: 0 10px 11px rgb(var(--arcoblue-3), 0.08), 0 6px 4px rgb(var(--arcoblue-3), 0.06),
        0 0 0 1px rgb(var(--arcoblue-3), 0.05), 0 2.89797px 2.12518px rgb(var(--arcoblue-3), 0.04),
        0 1.87823px 1.37737px rgb(var(--arcoblue-3), 0.03), 0 1.18233px 0.867039px rgb(var(--arcoblue-3), 0.02),
        0 0.67932px 0.498168px rgb(var(--arcoblue-3), 0.02), 0 0.298986px 0.219257px rgb(var(--arcoblue-3), 0.01);
    }
    :deep(.arco-form-item) {
      margin-bottom: 10px;
    }
    :deep(.arco-form-item-label-col) {
      margin-bottom: 0px;
    }
  }
</style>
