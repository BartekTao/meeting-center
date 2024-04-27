<template>
  
    <div class="col-lg-12">
      <div class="listing-item">
        <div class="left-image">
          <a><img :src="image_url" :alt="item.name"></a>
        </div>
        <div class="right-content align-self-center">
          <a><h4>會議室名稱：{{ item.roomId }}</h4></a>
          <!-- <div>{{ showReservator }}</div> -->
          <ItemPeriod 
            period-name="早上："
            :reservator-list="item.reservatorList.slice(0, 6)"
            :info-progress-width="250"
            :margin-left='0'
            @update-show-reservator="updateShowReservator"
          />
          <ItemPeriod 
            period-name="下午："
            :reservator-list="item.reservatorList.slice(6)"
            :info-progress-width="500"
            :margin-left='0'
            @update-show-reservator="updateShowReservator"
          />
          <div style="height: 20px;"></div>
          <ul class="info" style="padding-left: 0rem;">
            <li>人數限制：{{ item.capacity }}</li>  
            <li>有大桌子：{{ item.equipment.includes('big table') ? '是' : '否' }}</li>
            <li>有投影機：{{ item.equipment.includes('projector') ? '是' : '否' }}</li>
            <li>可否進食：{{ item.rules.includes('no food') ? '否' : '是' }}</li>
            <li>可否喝水：{{ item.rules.includes('no drinks') ? '否' : '是' }}</li>
          </ul><br>

          <div class="flex-container">
            <div class="main-white-button">
              <a class="openFormBtn" v-if="bookingAction" @click="$emit('openForm', item)"><img :src="tapImage" alt="Booking">預約</a>
            </div>
            <div class="main-white-button">
              <a class="openFormBtn" v-if="editAction"  @click="$emit('openForm', item)"><img :src="docImage" alt="Edit">編輯</a>
            </div>
            <div class="main-white-button">
              <a class="openFormBtn" v-if="deleteAction" ><img :src="deleteImage" alt="Delete">刪除</a>
            </div>
          </div>

        </div>
      </div>
    </div>
  </template>
  
  <script>
  import ItemPeriod from './ItemPeriod.vue';

  export default {
    name: 'ReserveBlock',
    props: ['item', 'bookingAction', 'editAction', 'deleteAction'],
    methods: {
      updateShowReservator(value) {
        this.showReservator = value;
      }
    },
    data() {
        return {
        booking_action: true,
        edit_action: true,
        delete_action: true,
        tapImage: require('@/assets/images/tap.png'),
        docImage: require('@/assets/images/google-docs.png'),
        deleteImage: require('@/assets/images/delete.png'),
        image_url: require('../assets/images/listing-01.jpg'),
        showReservator: '',
        };
    },
    components: {
        ItemPeriod
    }
  }
  </script>
  