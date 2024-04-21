<template>
    <div>
      <head-page pageContent="後台管理"></head-page>
      <manager-item></manager-item>
    </div>
</template>
  
  <script>
  import HeadPage from '@/components/HeadPage.vue';
  import ManagerItem from '@/components/ManagerItem.vue';
  
  export default {
    components: {
      HeadPage,
      ManagerItem
    }
  }
  </script>
  