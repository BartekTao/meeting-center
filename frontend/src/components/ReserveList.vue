<!-- ReserveList.vue -->
<template>
    <!-- <div id="detailInfo" :style="showDivStyle">{{ reservator }}</div> -->

    <div class="card text-bg-info mb-3 white-text"  :style="showDivStyle">
      <div class="card-header">詳細資訊</div>
      <div class="card-body">
        <h5 class="card-title">預約人：{{ reservator }}</h5>
        <h5 class="card-title">開始時間：{{ start_time }}</h5>
        <h5 class="card-title">結束時間：{{ end_time }}</h5>
      </div>
    </div>

    <div class="recent-listing" id="items">
      <div class="container">
        <div class="row">
          <!-- <div class="col-lg-12">
            <div class="section-heading">
              <h2>查詢空間</h2>
            </div>
          </div> -->
          <div class="col-lg-12">
            <div class="">
              <div class="item">
                <div class="row">
                  <ReserveBlock v-for="item in room_list" :key="item.name" :item="item" @showDiv="showDiv" @hideDiv="hideDiv" @openForm="openForm" :bookingAction="bookingAction" :editAction="editAction" :deleteAction="deleteAction"/>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>
  
  <script>
  import { ref, reactive, getCurrentInstance } from 'vue';
  import ReserveBlock from './ReserveBlock.vue';
  
  export default {
    name: 'ReserveList',
    data() {
      return {
        room_list: [
          {
              id: '6629c2edd7d285f521a5d787',
              roomId: "test",
              capacity: 10,
              equipment: [],
              rules: ["no food", "no drinks" ],
              reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
          {
              id: '6629c2edd7d285f521a5d787',
              roomId: "test",
              capacity: 10,
              equipment: [],
              rules: ["no food", "no drinks" ],
              reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
          {
              id: '6629c2edd7d285f521a5d787',
              roomId: "test",
              capacity: 10,
              equipment: [],
              rules: ["no food", "no drinks" ],
              reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
        ]
      };
    },
    components: {
      ReserveBlock
    },
    props: ['openForm', 'bookingAction', 'editAction', 'deleteAction'],
    setup() {
      const instance = getCurrentInstance();
      const time_period = ref([]);
      const reservator = ref('');
      const start_time = ref('');
      const end_time = ref('');

      if (instance && instance.appContext.config.globalProperties.$names) {
        time_period.value = instance.appContext.config.globalProperties.$names;
      }

      const showDivStyle = reactive({
        display: 'none',
        position: 'absolute',
        maxWidth: '18rem',
        zIndex: 1000,
      });

      function showDiv(data) {
        showDivStyle.display = 'block';
        showDivStyle.left = data.event.pageX + 'px';
        showDivStyle.top = data.event.pageY + 'px';
        reservator.value = data.unit;
        start_time.value = time_period.value[data.index];
        end_time.value = time_period.value[data.index+1];
      }

      function hideDiv() {
        showDivStyle.display = 'none';
      }

      return {
        showDivStyle,
        time_period,
        start_time, 
        end_time,
        reservator,
        showDiv,
        hideDiv
      };
    }
  }
  </script>
  