<!-- ReserveList.vue -->
<template>
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
                  <ReserveBlock v-for="item in room_list" :key="item.name" :item="item" @openForm="openForm" :bookingAction="bookingAction" :editAction="editAction" :deleteAction="deleteAction"/>
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
  }
  </script>
  