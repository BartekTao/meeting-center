<template>
  
  <comm-with-gql ref="commWithGql"></comm-with-gql>
    <div class="col-lg-12">
      <div class="listing-item">
        <div class="left-image">
          <a><img :src="image_url" :alt="item.name"></a>
        </div>
        <div class="right-content align-self-center">
          <a v-if="this.pageState === 'reserved'">
            <h4>會議室：{{ item.name }} | 會議名稱：{{ item.title }} </h4>
          </a>
          <a v-else>
            <h4>會議室：{{ item.name }}</h4>
          </a>
          <!-- <a>日期：5/28 | 時段：10:00~12:00 </a> -->
          <!-- <a>日期：{{ item.eventDay }} | 時段：{{ item.start_time }}~{{ item.end_time }} </a> -->
          <a v-if="this.pageState === 'reserved'">
            日期：{{ item.eventDay }} | 時段：{{ item.start_time }}~{{ item.end_time }}
          </a>
          <a></a>
          <ItemPeriod 
            period-name="早上："
            :reservator-list="item.schedulesList.slice(0, 6)"
            :info-progress-width="300"
            :margin-left='0'
            @update-show-reservator="updateShowReservator"
            @showDiv="$emit('showDiv', $event)"
            @hideDiv="$emit('hideDiv')"
          />
          <ItemPeriod 
            period-name="下午："
            :reservator-list="item.schedulesList.slice(6)"
            :info-progress-width="600"
            :margin-left='0'
            @update-show-reservator="updateShowReservator"
            @showDiv="$emit('showDiv', $event)"
            @hideDiv="$emit('hideDiv')"
          />
          <div style="height: 20px;"></div>
          <ul class="info" style="padding-left: 0rem;">
            <li>人數限制：{{ item.capacity }}</li>  
            <li>有大桌子：{{ item.equipments.includes('TABLE') ? '是' : '否' }}</li>
            <li>有投影機：{{ item.equipments.includes('PROJECTOR') ? '是' : '否' }}</li>
            <li>可否進食：{{ item.rules.includes('NO_FOOD') ? '否' : '是' }}</li>
            <li>可否喝水：{{ item.rules.includes('NO_DRINK') ? '否' : '是' }}</li>
            <!-- <li>schedules：{{ item.schedulesList }}</li> -->
          </ul><br>

          <div class="flex-container">
            <div class="main-white-button">
              <a class="openFormBtn" v-if="bookingAction" @click="$emit('openForm', item)"><img :src="tapImage" alt="Booking">預約</a>
            </div>
            <div class="main-white-button">
              <a class="openFormBtn" v-if="editAction"  @click="$emit('openForm', item)"><img :src="docImage" alt="Edit">編輯</a>
            </div>
            <div class="main-white-button">
              <a class="openFormBtn" v-if="deleteAction" @click="deleteEvent(item)" ><img :src="deleteImage" alt="Delete">刪除</a>
            </div>
            <div class="main-white-button">
              <a class="openFormBtn" v-if="editCommentAction"  @click="$emit('openCommentForm', item)"><img :src="docImage" alt="Edit">編輯會議結論</a>
            </div>
          </div>

        </div>
      </div>
    </div>
</template>
  
  <script>
  import ItemPeriod from './ItemPeriod.vue';
  import CommWithGql from '@/components/CommWithGql.vue'

  export default {
    name: 'ReserveBlock',
    emits: ['showDiv', 'hideDiv', 'openForm', 'openCommentForm', 'update-form'],
    props: ['pageState', 'item', 'bookingAction', 'editAction', 'deleteAction', 'editCommentAction'],
    methods: {
      updateShowReservator(value) {
        this.showReservator = value;
      },
      deleteEvent(item) {

        console.log('item:', item);
        this.$refs.commWithGql.deleteEvent(item.eventId);
        this.$emit('update-form');
      },
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
        ItemPeriod,
        CommWithGql
    }
  }
  </script>
  