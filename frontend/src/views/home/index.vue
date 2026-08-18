<template>
  <div class="container">
    <div class="left-side">
      <div class="panel">
        <!-- <Welcome /> -->
        <ListTotal />
      </div>
      <div class="panel" style="margin-top: 10px;">
        <ContentChart />
      </div>
      <div class="panel" style="margin-top: 10px;">
        <DemoCards />
      </div>
    </div>
    <div class="right-side">
      <a-grid :cols="24" :row-gap="10">
        <a-grid-item :span="24">
          <div class="panel moduler-wrap">
            <QuickOperation />
          </div>
        </a-grid-item>
        <a-grid-item class="panel" :span="24">
          <CategoriesPercent />
        </a-grid-item>
        <a-grid-item class="panel" :span="24">
          <RightCards />
        </a-grid-item>
      </a-grid>
    </div>
  </div>
</template>

<script setup>
// import Welcome from './components/Welcome.vue'
import QuickOperation from './components/quick-operation.vue';
import ContentChart from './components/content-chart.vue';
import CategoriesPercent from './components/categories-percent.vue';
import ListTotal from './components/list-total.vue';
import DemoCards from './components/demo-cards.vue';
import RightCards from './components/right-cards.vue';
</script>
<style scoped lang="less">
  .container {
    padding: 15px;
    display: flex;
    height: 100%;
    box-sizing: border-box;
    overflow-y: auto;
    .left-side {
      flex: 1;
      min-width: 0;
    }
    .right-side {
      width: 320px;
      flex-shrink: 0;
      margin-left: 16px;
      .panel {
        background-color: var(--color-neutral-1);
        border-radius: 4px;
        overflow: auto;
      }
    }
  }
 
</style>