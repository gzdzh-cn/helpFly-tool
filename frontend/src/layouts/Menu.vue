<template>
  <a-layout id="app-menu">
    <a-layout-sider
      theme="light"
      class="layout-sider"
      style="width: 120px;"
    >
      <a-menu 
        theme="light" 
        mode="inline" 
        :selectedKeys="[current]"
        @menu-item-click="changeMenu">
        <a-menu-item v-for="item in menulist" :key="item.name">
          <icon-font :type="item?.meta.icon" /><span >{{ item?.meta.title }}</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout class="submenu-layout">
      <a-layout-content class="submenu-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
  import { ref,computed } from 'vue';
  import { useRouter,useRoute } from 'vue-router';
  import routerMap from '@/router/routerMap';
  const route = useRoute();
  const router = useRouter();
  const current=ref(route.name)
  const menulist = computed(() => {
    const childrenlist= routerMap.find((item)=>item.path=="/")?.children
    const nowname=route.path.split("/")
    return childrenlist.find((item)=>item.name==nowname[1])?.children
  });
  const changeMenu=(key)=>{
    current.value = key;
    const linkInfo =menulist.value.find((item)=>item.name==key)
    if(linkInfo)
    router.push({ name: linkInfo.name, params: linkInfo.params})
  }
</script>
<style lang="less" scoped>
#app-menu {
  height: 100%;
  width: 100%;
  text-align: center;
  min-width: 0;
  overflow: hidden;
  overscroll-behavior: none;
  .submenu-layout,
  .submenu-content {
    flex: 1 1 auto;
    width: auto;
    min-width: 0;
    min-height: 0;
    overflow-y: auto;
    overscroll-behavior: none;
  }
  .layout-sider {
    // border-top: 1px solid var(--color-neutral-3);
    border-right: 1px solid var(--color-neutral-3);
    // background-color: #FAFAFA;
    overflow: auto;
    :deep(.arco-menu .arco-menu-item .arco-icon){
      margin-right: 5px;
      display: inline-block;
    }
  }
}
</style>
