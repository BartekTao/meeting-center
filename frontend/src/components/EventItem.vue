<template>
  
    <div class="col-lg-12">
      <div class="listing-item">
        <div class="left-image">
          <a><img :src="item.image_url" :alt="item.name"></a>
        </div>
        <div class="right-content align-self-center">
          <a><h4>會議室名稱：{{ item.name }}</h4></a>
          <!-- <div>{{ showReservator }}</div> -->
          <ItemPeriod 
            period-name="早上："
            :reservator-list="item.reservatorList.slice(0, 6)"
            :info-progress-width="250"
            @update-show-reservator="updateShowReservator"
          />
          <ItemPeriod 
            period-name="下午："
            :reservator-list="item.reservatorList.slice(6)"
            :info-progress-width="500"
            @update-show-reservator="updateShowReservator"
          />
          <div style="height: 20px;"></div>
          <ul class="info" style="padding-left: 0rem;">
            <li>人數限制：{{ item.people_limit }}</li>
            <li>可否飲食：{{ item.can_eat ? '是' : '否' }}</li>
          </ul><br>
          <div class="main-white-button">
            <a class="openFormBtn" @click="$emit('openForm', item)"><img :src="$parent.tapImage" alt="Booking">預約</a>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import ItemPeriod from './ItemPeriod.vue';

  export default {
    name: 'EventItem',
    props: ['item'],
    methods: {
      openForm() {
        this.$emit('openForm');
      },
      updateShowReservator(value) {
        this.showReservator = value;
      }
    },
    data() {
        return {
        showReservator: '',
        };
    },
    components: {
        ItemPeriod
    }
  }
  </script>
  