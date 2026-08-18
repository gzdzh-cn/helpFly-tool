<template>
    <router-view/>
</template>
<script lang="ts" setup>
import { useAppStore } from '@/store';
const appStore = useAppStore();
//Events.On事件需要在这里监听一次才能使他在router路由界面生效
//获取建议所有的监听都在这里写，然后用pinia同步数据
import {Events} from "@wailsio/runtime";
Events.On('time', (time:any) => {
    appStore.sysTime=time.data;
});
</script>
<style lang="less">
#app{
    color: var(--color-neutral-10);
}
</style>  