<template>
    <div>
      <head-page pageContent="後台管理"></head-page>
      <create-room></create-room>
      <manager-item></manager-item>
    </div>
</template>
  
  <script>
  import HeadPage from '@/components/HeadPage.vue';
  import ManagerItem from '@/components/ManagerItem.vue';
  import CreateRoom from '@/components/CreateRoom.vue'
  
  export default {
    components: {
      HeadPage,
      ManagerItem,
      CreateRoom
    }
  }
  </script>
  