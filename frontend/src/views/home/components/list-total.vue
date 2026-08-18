<template>
  <a-card class="general-card" title="数据统计" :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '24px 20px 0 20px' }">

    <a-row :gutter="8">
      <template v-for="link in links">
        <a-col :span="8" class="wrapper">
          <div class="icon-box">
            <icon-font :type="link.icon" class="icon" size="26" />
          </div>
          <a-typography-paragraph class="text">
            <div class="value-box">
              <span class="value">{{ link.value }} </span>
              <span class="unit">{{ link.unit }}</span>
            </div>
            <span class="text-title">
              {{ link.text }}
            </span>

          </a-typography-paragraph>
        </a-col>

      </template>
    </a-row>
    <!-- <a-divider class="split-line" style="margin: 0" /> -->
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import useAppStore from '@/store/modules/app';

const appStore = useAppStore();
const loading = ref(false);

onMounted(() => {
  loading.value = true;

});


const links = [
  { text: '用户总数', icon: 'icon-yonghushuliang1', type: "button", value: "538", unit: '个' },
  { text: '表单总数', icon: 'icon-biaodan2', type: "browser", value: "1020", unit: '个' },
  { text: '表单类型占比', icon: 'icon-icon_percent', type: "a", value: "50", unit: '%' },
];
</script>

<style scoped lang="less">
.general-card {
  background: transparent;
  height: 180px;

  .wrapper {
    cursor: pointer;

    .icon-box {
      display: flex;
      align-items: center;
      justify-content: center;

      .icon {
        color: #fff;
        font-size: 26px;
      }

  

    }

    .text {
      text-align: center;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;

      .value-box {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;

        .value {
          font-size: 24px;
          font-weight: bold;
          color: #fff;
        }

        .unit {
          font-size: 12px;
          color: #c2c1c1;
        }
      }

      .text-title {
        font-size: 12px;
        // font-weight: bold;
        color: #fff;
      }

    }
  }
}
</style>