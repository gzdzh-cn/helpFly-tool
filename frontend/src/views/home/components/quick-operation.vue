<template>
  <a-card
    class="general-card"
    title="快捷操作"
    :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '24px 20px 0 20px' }"
  >
    <!-- <template #extra>
      <a-link>管理</a-link>
    </template> -->
    <a-row :gutter="8">
      <template v-for="link in links" >
        <a-col :span="8" v-if="link.type=='button'" class="wrapper" @click="handleOpen(link.value)">
          <div class="icon">
            <icon-font :type="link.icon" class="icon" size="26"/>
          </div>
          <a-typography-paragraph class="text">
            {{ link.text }}
          </a-typography-paragraph>
        </a-col>
        <a-col :span="8" v-else-if="link.type=='browser'" class="wrapper">
            <a :wml-openURL="link.value">
              <div class="icon">
                <icon-font :type="link.icon" class="icon" size="26"/>
              </div>
              <a-typography-paragraph class="text">
                {{ link.text }}
              </a-typography-paragraph>
            </a>
        </a-col>
        <a-col :span="8" v-else-if="link.type=='a'" class="wrapper">
            <a :href="link.value"  target="_blank">
              <div class="icon">
                <icon-font :type="link.icon" class="icon" size="26"/>
              </div>
              <a-typography-paragraph class="text">
                {{ link.text }}
              </a-typography-paragraph>
            </a>
        </a-col>
      </template>
    </a-row>
    <!-- <a-divider class="split-line" style="margin: 0" /> -->
  </a-card>
</template>

<script lang="ts" setup>
import {onMounted } from 'vue';
import { Notification } from '@arco-design/web-vue';
import {WML} from "@wailsio/runtime";
onMounted(async()=>{
    WML.Reload()
})
  const links = [
    { text: '通知提醒框', icon: 'icon-filled',type:"button",value:"notification" },
    { text: '默认浏览器', icon: 'icon-wangye',type:"browser",value:"https://goflys.cn"},
    { text: 'Webview', icon: 'icon-wangye',type:"a",value:"https://v3alpha.wails.io"},
    // { text: 'workplace.onlinePromotion', icon: 'icon-mobile' },
    // { text: 'workplace.contentPutIn', icon: 'icon-fire' },
  ];
  //执行UI组件
  const handleOpen=(val:string)=>{
    if(val=="notification"){
        Notification.info({
          title: 'Notification',
          content: 'This is a notification!',
          position: 'bottomRight',
          closable: true,
        })
    }
  }
</script>

<style scoped lang="less">
.general-card{
  background: transparent;
  height: 180px;
 
  .wrapper{
    cursor: pointer;
    .icon{
      text-align: center;
    }
    .text{
      text-align: center;
    }
  }
}
</style>
