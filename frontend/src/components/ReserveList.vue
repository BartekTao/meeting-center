<!-- ReserveList.vue -->
<template>
    <div id="detailInfo" :style="showDivStyle">{{ content }}</div>
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
  import { ref, reactive } from 'vue';
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
      const showDivStyle = reactive({
        width: '100px',
        height: '100px',
        backgroundColor: 'red',
        display: 'none',
        position: 'absolute',
        zIndex: 1000,
      });
      const content = ref('');

      function showDiv(data) {
        showDivStyle.display = 'block';
        showDivStyle.left = data.event.pageX + 'px';
        showDivStyle.top = data.event.pageY + 'px';
        content.value = data.unit;
        console.log(data.index);
      }

      function hideDiv() {
        showDivStyle.display = 'none';
      }

      return {
        showDivStyle,
        showDiv,
        hideDiv,
        content
      };
    }
  }
  </script>
  