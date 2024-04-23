<template>
  <QueryBar />
  <QueryItems />
</template>

<script>
import QueryBar from '@/components/QueryBar.vue'
import QueryItems from '@/components/QueryItems.vue'

export default {
  name: 'QueryPage',
    data() {
      return {
        userName: "Ray",
        time_period: ['9:00', '9:30', '10:00', '10:30', '11:00', '11:30',
          '12:00', '12:30', '13:00', '13:30', '14:00', '14:30', '15:00',
          '15:30', '16:00', '16:30', '17:00', '17:30', '18:00']
      };
    },
  components: {
    QueryBar,
    QueryItems
  }
}
</script>
