<template>
  <a-card
    class="general-card"
    title="http请求"
    :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '15px' }"
  >
    <template #extra>
      <!-- <a-link>管理</a-link> -->
    </template>
    <div class="clipboard">
      <a-form :model="formData" layout="vertical" >
        <a-form-item field="text" :extra="resAxios">
          <template #label>
            <span style="font-weight: 600;padding-right: 10px;font-size: 15px;">前端请求</span>
            <a-radio-group v-model="requestType.axios">
              <a-radio value="GET">GET</a-radio>
              <a-radio value="POST">POST</a-radio>
            </a-radio-group>
          </template>
          <a-input v-model="formData.paramAxios" placeholder="输入请求参数"  />
          <a-button type="primary" @click="handleAxios">请求</a-button>
        </a-form-item>
        <a-form-item field="text" :extra="resService">
          <template #label>
            <span style="font-weight: 600;padding-right: 10px;font-size: 15px;">service端请求</span>
            <a-radio-group v-model="requestType.service">
              <a-radio value="GET">GET</a-radio>
              <a-radio value="POST">POST</a-radio>
            </a-radio-group>
          </template>
          <a-input v-model="formData.paramService" placeholder="输入请求参数" allow-clear />
          <a-button type="primary" @click="handleService">请求</a-button>
        </a-form-item>
      </a-form>
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Message } from '@arco-design/web-vue';
//APi
import { getData,postData} from './api';
  import {HttpService} from "/#/helpfly/internal/service";
  const formData=ref({
    paramAxios: "",
    paramService: "",
  })
  const requestType=ref({
    axios: "GET",
    service: "GET",
  })

  //前端axios-http请求(GET)
  const resAxios=ref()
  const handleAxios=async()=>{
    try {
      Message.loading({content:"请求中，请稍后",id:"http",duration:0})
       if(requestType.value.axios=="GET"){
        resAxios.value= await getData({data:formData.value.paramAxios});
        Message.success({content:"前端axios-http-GET请求",id:"http",duration:2000})
       }else{
        resAxios.value= await postData({data:formData.value.paramAxios});
        Message.success({content:"前端axios-http-POST请求",id:"http",duration:2000})
       }
    } catch (err:any) {
      Message.loading({content:"请求中，请稍后",id:"http",duration:1})
    } finally {
    }
  }
  //Go的Service端请求http
  const resService=ref()
  const handleService=async()=>{
    Message.loading({content:"调用中·请稍后",id:"http",duration:0})
    if(requestType.value.service=="GET"){
      HttpService.GetData(formData.value.paramService).then((result) => {
        resService.value= result
        Message.success({content:"调用Service端请求http-GET成功",id:"http",duration:2000})
      }).catch((err:any) => {
        Message.error({content:err,id:"http",duration:2000})
      });
    }else{
      HttpService.PostData(formData.value.paramService).then((result) => {
        resService.value= result
        Message.success({content:"调用Service端请求http-POST成功",id:"http",duration:2000})
      }).catch((err:any) => {
        Message.error({content:err,id:"http",duration:2000})
      });
    }
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
    margin-bottom: 15px;
  }
  :deep(.arco-form-item-label-col){
    margin-bottom: 5px;
  }
  :deep(.arco-radio-group .arco-radio){
    margin-right: 8px;
  }
}
</style>
