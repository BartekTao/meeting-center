<!-- ReserveList.vue -->
<template>
    <!-- <div id="detailInfo" :style="showDivStyle">{{ reservator }}</div> -->

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
                  <ReserveBlock v-for="item in room_list" :key="item.name" :item="item" @showDiv="$emit('showDiv', $event)" @hideDiv="$emit('hideDiv')" @openForm="openForm" @openCommentForm="openCommentForm" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction"/>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>
  
  <script>
  import ReserveBlock from './ReserveBlock.vue';
  
  export default {
    name: 'ReserveList',
    emits: ['showDiv', 'hideDiv'],
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
    props: ['openForm', 'openCommentForm', 'bookingAction', 'editAction', 'deleteAction', 'editCommentAction']
  }
  </script>
  