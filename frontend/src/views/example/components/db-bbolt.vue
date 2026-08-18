<template>
  <a-card
    class="general-card"
    title="boltDB数据库"
    :header-style="{ paddingBottom: '0' }"
    :body-style="{ padding: '10px 10px 0 10px' }"
  >
    <template #extra>
      <!-- <a-link>管理</a-link> -->
    </template>
    <div class="clipboard">
      <a-form :model="formData" size="large">
        <a-form-item field="text" hide-label :extra="'获取DB数据：'+resService">
          <a-input v-model="formData.text" placeholder="输入存储数据"  />
          <a-button type="primary" @click="handleReplication">保存</a-button>
        </a-form-item>
      </a-form>
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import {DbService} from "/#/helpfly/internal/service";
  const formData=ref({
    text: "",
  })
  //复制
  const resService=ref("")
  const handleReplication=async()=>{
    Message.loading({content:"调用中·请稍后",id:"http",duration:0})
    DbService.SetData(formData.value.text).then((result) => {
        resService.value= result
        Message.success({content:"测试保存DB数据成功",id:"http",duration:2000})
      }).catch((err:any) => {
        Message.error({content:err,id:"http",duration:2000})
      });
  }
  //获取DB中的值
  const GetDBData=async()=>{
    Message.loading({content:"调用中·请稍后",id:"http",duration:0})
    DbService.GetData().then((result) => {
        resService.value= result
        Message.success({content:"测试获取DB数据成功",id:"http",duration:2000})
      }).catch((err:any) => {
        Message.error({content:err,id:"http",duration:2000})
      });
  }
  GetDBData()
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
}
</style>
