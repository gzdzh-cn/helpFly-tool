<template>
  <a-card
    class="general-card"
    title="剪贴板"
    :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '10px 10px 0 10px' }"
  >
    <template #extra>
      <!-- <a-link>管理</a-link> -->
    </template>
    <div class="clipboard">
      <a-form :model="formData" layout="vertical" size="large">
        <a-form-item field="text" label="剪贴板内容">
          <a-input v-model="formData.text" placeholder="输入剪贴板内容"  />
          <a-button type="primary" @click="handleReplication">复制</a-button>
        </a-form-item>
        <a-form-item field="text" label="剪贴板内容">
          <a-input v-model="formData.copyText" placeholder="粘贴的内容" allow-clear />
          <a-button type="primary" @click="handleCopy">粘贴</a-button>
        </a-form-item>
      </a-form>
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import {Clipboard} from "@wailsio/runtime";
  const formData=ref({
    text: "这是一段剪贴板内容，可以粘贴到文本框中。",
    copyText: "",
  })
  //复制
  const handleReplication=async()=>{
     await Clipboard.SetText(formData.value.text)
     Message.success({content:"复制成功",id:"clipboard",duration:2000})
  }
  //粘贴
  const handleCopy=async()=>{
     const clipboardText=await Clipboard.Text()
     formData.value.copyText=clipboardText
     Message.success({content:"粘贴成功",id:"clipboard",duration:2000})
  }
</script>

<style scoped lang="less">
.general-card{
  background: var(--color-neutral-1);
  .wrapper{
    cursor: pointer;
    margin-bottom: 10px;
    padding: 5px 0px;
    border-radius: 3px;
    .icon{
      text-align: center;
    }
    .text{
      text-align: center;
      margin-bottom: 0px;
    }
  }
  :deep(.arco-form-item){
    margin-bottom: 10px;
  }
  :deep(.arco-form-item-label-col){
    margin-bottom: 5px;
  }
}
</style>
