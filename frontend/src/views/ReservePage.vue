<template>
  <ReserveBar/>
  <ReserveList  @showDiv="showDiv" @hideDiv="hideDiv" :openForm="openForm" :bookingAction="bookingAction" :editAction="editAction" :deleteAction="deleteAction"/>
  <ReserveForm  @showDiv="showDiv" @hideDiv="hideDiv" :formDisplay="formDisplay" :roomInfo="roomInfo" :userName="userName" @close-form="closeForm"/>
  <EventInfo ref="eventInfo"/>
</template>

<script>
import ReserveBar from '@/components/ReserveBar.vue'
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';
import EventInfo from '@/components/EventInfo.vue';

export default {
  name: 'ReservePage',
  components: {
    ReserveBar,
    ReserveList,
    ReserveForm,
    EventInfo
  },
  data() {
    return {
      bookingAction: true,
      editAction: false,
      deleteAction: false,
      formDisplay: false,
      roomInfo: {
          reservatorList: [],
          roomName: '',
      },
      userName: "Ray",
      showDivStyle: {
        display: 'none',
        position: 'absolute',
        maxWidth: '18rem',
        zIndex: 1000,
      },
      time_period: [],
      reservator: '',
      start_time: '',
      end_time: ''
    };
  },
  methods: {
    openForm(item) {
      this.formDisplay = true;
      this.roomInfo.id = item.id;
      this.roomInfo.roomName = item.roomId;
      this.roomInfo.reservatorList = item.reservatorList;
    },
    closeForm() {
      this.formDisplay = false;
    },
    showDiv(data) {
      this.$refs.eventInfo.showDiv(data);
    },
    hideDiv() {
      this.$refs.eventInfo.hideDiv();
    }
  },
  mounted() {
    this.time_period = this.$names;
  }
}
</script>


<style>

.info {
  display: flex;        /* 使用 Flexbox 佈局 */
  flex-wrap: wrap;      /* 允許項目換行 */
  list-style: none;     /* 去除列表前的標記 */
  padding: 0;           /* 去除預設的內距 */
  margin: 0;            /* 去除預設的外距 */
}

.info li {
  flex: 0 0 20%;        /* 每個 li 佔總寬度的 20% */
  box-sizing: border-box; /* 盒子模型調整，邊界和內距不再另外計算 */
  padding: 0.5rem;      /* 略微增加內距以便閱讀 */
  border: 1px solid #ccc; /* 加上邊框以視覺區分 */
  text-align: center;   /* 文字置中 */
  background-color: white; /* 淺藍色背景 */
}

.bordered {
  border: 1px solid white;
}

.unselectable {
    user-select: none; /* 防止文本被选取 */
    -webkit-user-select: none; /* Safari 浏览器 */
    -moz-user-select: none; /* Firefox 浏览器 */
    -ms-user-select: none; /* Internet Explorer/Edge 浏览器 */
  }
  .white-text {
    color: white !important;
  }

</style>