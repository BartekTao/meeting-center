<template>
  <QueryBar />
  <QueryItem />
</template>

<script>
import QueryBar from '@/components/QueryBar.vue'
import QueryItem from '@/components/QueryItem.vue'

export default {
  name: 'QueryPage',
  components: {
    QueryBar,
    QueryItem
  }
}
</script>
