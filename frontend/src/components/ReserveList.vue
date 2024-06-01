<!-- ReserveList.vue -->
<template>
    <!-- <div id="detailInfo" :style="showDivStyle">{{ reservator }}</div> -->
    <comm-with-gql ref="commWithGql"></comm-with-gql>

    <div class="recent-listing" id="items">
      <div class="container">
        <div class="row">
          <div class="col-lg-12">
            <div class="">
              <div class="item">
                <div class="row">
                  <ReserveBlock v-for="item in roomItems" :key="item.name" :item="item" @showDiv="$emit('showDiv', $event)" @hideDiv="$emit('hideDiv')" @update-form="$emit('update-form')" @openForm="openForm" @openCommentForm="openCommentForm" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction" :pageState="pageState"/>
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
  import CommWithGql from '@/components/CommWithGql.vue'
  
  export default {
    name: 'ReserveList',
    emits: ['update-form', 'showDiv', 'hideDiv'],
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
      ReserveBlock,
      CommWithGql
    },
    props: ['pageState', 'openForm', 'openCommentForm', 'bookingAction', 'editAction', 'deleteAction', 'editCommentAction', 'roomItems']
  }
  </script>
  