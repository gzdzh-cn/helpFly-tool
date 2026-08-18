<template>
  <a-card
    class="general-card"
    title="事件演示"
    :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '15px 15px 0px 15px' }"
  >
    <template #extra>
      <!-- <a-link>管理</a-link> -->
    </template>
    <a-row :gutter="8">
      <a-col v-for="link in links" :key="link.text" @click="handleEvent(link.event)" :span="8" class="wrapper">
        <div class="icon">
          <icon-font :type="link.icon" class="icon" size="26"/>
        </div>
        <a-typography-paragraph class="text">
          {{ link.text }}
        </a-typography-paragraph>
      </a-col>
    </a-row>
  </a-card>
</template>

<script lang="ts" setup>
  import {Events,Browser} from "@wailsio/runtime";
  import { Message } from '@arco-design/web-vue';
  import {OpenWindow} from "/#/helpfly/internal/service";
  const links = [
    { text: '触发事件', icon: 'icon-zhihangshijian',event:"emitgo" },
    { text: '打开新窗口', icon: 'icon-dakaixinchuangkou-' ,event:"newwindow"},
    { text: '浏览器网页', icon: 'icon-wangye' ,event:"weburl"},
  ];
  //执行事件
  const handleEvent=async(event:string)=>{
    if(event=="emitgo"){
    const Emitname = new Events.WailsEvent("myevent","data88")
     const res= await Events.Emit(Emitname)
     console.log("事件执行结果：",res)
     Message.success({content:"触发service事件执行成功",id:"emit",duration:2000})
    }else if(event=="newwindow"){
      OpenWindow.Open("https://goflys.cn")
    }else if(event=="weburl"){
      Browser.OpenURL("https://goflys.cn")
    }
  }
</script>

<style scoped lang="less">
.general-card{
  background: transparent;
  .wrapper{
    cursor: pointer;
    margin-bottom: 10px;
    padding: 5px 0px;
    border-radius: 3px;
    &:hover{
      background: var(--color-neutral-3);
    }
    .icon{
      text-align: center;
    }
    .text{
      text-align: center;
      margin-bottom: 0px;
    }
  }
}
</style>
