<template>
  <div>
    <ReserveList :items="items" :openForm="openForm"/>
    <ReserveForm ref="reserveFormComponent"/>
  </div>
</template>
  
  <script>
  import ReserveList from './ReserveList.vue';
  import ReserveForm from './ReserveForm.vue';

  export default {
    name: 'QueryItems',
    data() {
      return {
        items: [
          {
            name: '001',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 12,
            can_eat: false,
            reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
          {
            name: '002',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 24,
            can_eat: false,
            reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
          {
            name: '003',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 32,
            can_eat: true,
            reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
        ]
      };
    },
    components: {
      ReserveList,
      ReserveForm
    },
    methods: {
    openForm(item) {
      if (this.$refs.reserveFormComponent) {
        this.$refs.reserveFormComponent.openForm(item);
      }
    }
  }
};
  </script>
  
  