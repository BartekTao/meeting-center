<template>
  <ReserveBar @updateAllRooms="updateAllRooms"/>
  <ReserveList  @showDiv="showDiv" @hideDiv="hideDiv" :openForm="openForm" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction" :roomItems="roomItems"/>
  <ReserveForm  @showDiv="showDiv" @hideDiv="hideDiv" :formDisplay="formDisplay" :formInfo="formInfo" @close-form="closeForm" @update-form="updateForm"/>
  <EventInfo ref="eventInfo"/>
  <comm-with-gql @fetch-available-rooms="fetchAvailableRooms" @query-users="queryUsers" ref="commWithGql"></comm-with-gql>
  <js-preloader ref="jsPreloader"></js-preloader>
</template>

<script>
import ReserveBar from '@/components/ReserveBar.vue'
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';
import EventInfo from '@/components/EventInfo.vue';
import CommWithGql from '@/components/CommWithGql.vue'
import JsPreloader from '@/components/JsPreloader.vue';

export default {
  name: 'ReservePage',
  components: {
    ReserveBar,
    ReserveList,
    ReserveForm,
    EventInfo,
    CommWithGql,
    JsPreloader
  },
  data() {
    return {
      startTimeStamp: '',
      endTimeStamp: '',
      bookingAction: true,
      editAction: false,
      editCommentAction: false,
      deleteAction: false,
      formDisplay: false,
      oneHourInMilliseconds: 3600000,
      users: [],

      eventInput: {
        title: "Team Meeting",
        description: "",
        startAt: 1716825600000,
        endAt: 1716831999999,
        roomId: "6655178de1dfe965fa4b1951",
        participantsIDs: ["6645ece136e2a0f035961bdd"],
        notes: "Bring all relevant documents",
        remindAt: 1625074200
            },

      formInfo: {
        title: '',
        roomId: '',
        roomName: '',
        userId: ['6645ece136e2a0f035961bdd'],
        userName: ['Ivan Lee'],
        namesString: 'Ivan Lee',
        eventId: '',
        dayTime: '',
        start_time: '10:00',
        end_time: '12:00',
        notes: 'test content',
        fileName: '',
        fileUrl: '',
        reservatorList: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ],
      },
      showDivStyle: {
        display: 'none',
        position: 'absolute',
        maxWidth: '18rem',
        zIndex: 1000,
      },
      roomItems: [],
    };
  },
  methods: {
    openForm(item) {
      this.formDisplay = true;
      this.formInfo.roomId = item.id;
      this.formInfo.roomName = item.name;
    },
    updateForm(formInfo) {
      this.eventInput.title = formInfo.title;

      const startTime = this.dayTime+'-'+formInfo.start_time + ':00';
      this.eventInput.startAt = this.transferToTimestamp(startTime);

      const endTime = this.dayTime+'-'+formInfo.end_time + ':00';
      this.eventInput.endAt = this.transferToTimestamp(endTime);

      this.eventInput.roomId = this.formInfo.roomId;
      const namesArray = formInfo.namesString.split(',');
      const idsArray = namesArray.map(name => {
        const user = this.users.find(user => user.name === name);
        return user ? user.id : null;
      });
      this.eventInput.participantsIDs = idsArray
      
      this.eventInput.notes = formInfo.notes;
      this.eventInput.remindAt = this.eventInput.startAt + this.oneHourInMilliseconds

      this.$refs.commWithGql.createEvent(this.eventInput);
    },
    closeForm() {
      this.formDisplay = false;
    },
    showDiv(data) {
      this.$refs.eventInfo.showDiv(data);
    },
    hideDiv() {
      this.$refs.eventInfo.hideDiv();
    },
    fetchAvailable(availables) {
      this.startTimeStamp = availables; // endTimeStamp
    },
    fetchAvailableRooms(rooms) {
      this.roomItems = rooms
    },
    queryUsers(users) {
      this.users = users
    },
    updateAllRooms(variables) {
      // this.$refs.commWithGql.queryUsers();
      this.dayTime = variables.dayTime;
      this.loadPreLoader(1000).then(() => {
        this.$refs.commWithGql.fetchAvailableRooms(variables);
      });
    },
    loadPreLoader(duration) {
      this.$refs.jsPreloader.isLoaded = false;
      return new Promise(resolve => {
        setTimeout(() => {
          this.$refs.jsPreloader.isLoaded = true;
          resolve();
        }, duration);
      });
    },
    transferToTimestamp(time) {
      const formattedTime_ = time.replace(/(\d{4})-(\d{2})-(\d{2})-(\d{1,2}):(\d{2}):(\d{2})/, '$1-$2-$3T$4:$5:$6');
      const formattedTime = formattedTime_.replace(/T(\d):/, 'T0$1:');
      const date = new Date(formattedTime);
      return date.getTime();
    },
  },
  mounted() {
    this.$refs.commWithGql.queryUsers();
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