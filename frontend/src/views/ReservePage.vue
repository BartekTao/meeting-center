<template>
  <ReserveBar />
  <ReserveList :items="items" :openForm="openForm"/>
  <ReserveForm ref="reserveFormComponent"/>
</template>

<script>
import ReserveBar from '@/components/ReserveBar.vue'
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';

export default {
  name: 'ReservePage',
  components: {
    ReserveBar,
    ReserveList,
    ReserveForm
  },
  data() {
    return {
      userName: "Ray",
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
        }
      ],
      time_period: ['9:00', '9:30', '10:00', '10:30', '11:00', '11:30',
        '12:00', '12:30', '13:00', '13:30', '14:00', '14:30', '15:00',
        '15:30', '16:00', '16:30', '17:00', '17:30', '18:00']
    };
  },
  methods: {
    openForm(item) {
      if (this.$refs.reserveFormComponent) {
        this.$refs.reserveFormComponent.openForm(item);
      }
    }
  }
}
</script>
