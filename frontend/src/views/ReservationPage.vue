<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <reservation-item></reservation-item>
  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReservationItem from '@/components/ReservationItem.vue';

export default {
  components: {
    HeadPage,
    ReservationItem
  }
}
</script>
