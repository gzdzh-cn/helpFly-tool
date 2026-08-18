<template>
  <div class="container">
    <a-row :gutter="8">
      <a-col :span="12">
        <div class="panel moduler-wrap">
          <baseOperation />
        </div>
        <div class="panel moduler-wrap">
          <httpRequest />
        </div>
      </a-col>
      <a-col :span="12">
        <div class="panel moduler-wrap">
          <eventOperation />
        </div>
        <div class="panel moduler-wrap">
          <clipboard />
        </div>
        <div class="panel moduler-wrap">
          <dbBbolt />
        </div>
    </a-col>
    </a-row>
  </div>
</template>

<script setup>
  import baseOperation from './components/base-operation.vue';
  import eventOperation from './components/event-operation.vue';
  import clipboard from './components/clipboard.vue';
  import httpRequest from './components/http-request.vue';
  import dbBbolt from './components/db-bbolt.vue';

</script>
<style scoped lang="less">
  .container {
    padding: 15px;
    .panel {
      background-color: var(--color-neutral-1);
      border-radius: 4px;
      overflow: auto;
      margin-bottom: 10px;
    }
  }

</style>