<template>
    <div class="recent-listing" id="items">
      <div class="container">
        <div class="row">
          <div class="col-lg-12">
            <div class="section-heading">
              <h2>查詢空間</h2>
            </div>
          </div>
          <div class="col-lg-12">
            <div class="">
              <div class="item">
                <div class="row">
                  <EventItem v-for="item in items" :key="item.name" :item="item" @openForm="openForm(item)"/>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ReserveForm ref="reserveFormComponent"/>
  </template>
  
  <script>
  import EventItem from './EventItem.vue';
  import ReserveForm from './ReserveForm.vue';

  export default {
    name: 'QueryItems',
    data() {
      return {
        tapImage: require('@/assets/images/tap.png'),
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
      EventItem,
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
  
  