<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <reserved-child-page></reserved-child-page>
  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReservedChildPage from '@/components/ReservedChildPage.vue';

export default {
  components: {
    HeadPage,
    ReservedChildPage
  }
}
</script>
