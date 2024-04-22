<template>
  <QueryBar />
  <QueryItems />
  <ItemPeriod />
</template>

<script>
import QueryBar from '@/components/QueryBar.vue'
import QueryItems from '@/components/QueryItems.vue'
import ItemPeriod from '@/components/ItemPeriod.vue';

export default {
  name: 'QueryPage',
  components: {
    QueryBar,
    QueryItems,
    ItemPeriod
  }
}
</script>
