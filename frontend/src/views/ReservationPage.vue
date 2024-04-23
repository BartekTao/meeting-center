<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <reserve-child-page></reserve-child-page>
  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReserveChildPage from '@/components/ReserveChildPage.vue';

export default {
  components: {
    HeadPage,
    ReserveChildPage
  }
}
</script>
