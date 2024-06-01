<template>
  <ReserveBar @updateAllRooms="updateAllRooms"/>
  <ReserveList  @showDiv="showDiv" @hideDiv="hideDiv" :openForm="openForm" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction" :roomItems="roomItems"/>
  <ReserveForm  @showDiv="showDiv" @hideDiv="hideDiv" :formDisplay="formDisplay" :formInfo="formInfo" :roomName="roomName" :schedulesList="schedulesList" :users="users" @close-form="closeForm" @update-form="updateForm"/>
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
      updateVariables: {},

      timeList: ['9:0', '9:30', '10:0', '10:30', '11:0', '11:30', '12:0', '12:30', '13:0', '13:30', '14:0', '14:30', '15:0', '15:30', '16:0', '16:30', '17:0', '17:30', '18:0'],
      reservatorList: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
    
      eventInput: {
        title: "Team Meeting",
        description: "",
        startAt: 1716825600000,
        endAt: 1716831999999,
        roomId: "6655178de1dfe965fa4b1951",
        participantsIDs: ["6645ece136e2a0f035961bdd"],
        notes: "",
        remindAt: 1625074200
            },

      formInfo: {
        title: 'test title',
        description: 'test description',
        roomId: '',
        roomName: '',
        userId: ['6645ece136e2a0f035961bdd'],
        userName: ['Ivan Lee'],
        namesString: 'Ivan Lee',
        eventId: '',
        dayTime: this.getDaytime(),
        start_time: '10:00',
        end_time: '12:00',
        notes: 'test content',
        fileName: '',
        fileUrl: '',
        reservatorList: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ],
        schedulesList: [],
      },
      showDivStyle: {
        display: 'none',
        position: 'absolute',
        maxWidth: '18rem',
        zIndex: 1000,
      },
      roomItems: [],
      roomName: '',
      roomId: '',
      schedulesList: [],
    };
  },
  methods: {
    openForm(item) {
      this.formDisplay = true;
      this.roomName = item.name;
      this.roomId = item.id;
      this.formInfo.roomId = item.id;
      this.schedulesList = item.schedulesList;
    },

    updateForm() {
      this.updateAllRooms({})
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
    fetchAvailable(available) {
      this.startTimeStamp = available; // endTimeStamp
    },
    fetchAvailableRooms(rooms) {
      this.roomItems = rooms;

      this.roomItems.forEach(room => {
        let finalReservatiorList = Array(this.timeList.length).fill('');

        room.schedules.forEach(schedule => {
          const reservatorName = this.findUserNameById(schedule.creator.id);
          const eventTitle = this.findUserNameById(schedule.title);
          const { hours: startHours, minutes: startMinutes } = this.getHours(schedule.startAt);
          const { hours: endHours, minutes: endMinutes } = this.getHours(schedule.endAt);
          
          let nickName;
          if (reservatorName.length === 3) {
            nickName = reservatorName.substring(1);
          } else {
            const nameParts = reservatorName.split(' ');
            nickName = nameParts[0];
          }

          const newReservatorList = this.transferReservatorList(startHours, startMinutes, endHours, endMinutes, reservatorName);
          const scheduleInfo = {
              id: schedule.creator.id,
              name: reservatorName,
              nickName: nickName,
              title: schedule.title,
              startHours,
              startMinutes,
              endHours,
              endMinutes,
              eventTitle
            };
          finalReservatiorList = finalReservatiorList.map((slot, index) => slot.name === reservatorName ? slot : newReservatorList[index] ? scheduleInfo : slot);
          
        });

        room.schedulesList = finalReservatiorList.slice(0, -1);
      });

    },
    transferReservatorList(startHours, startMinutes, endHours, endMinutes, eventName) {
      const startTime = startHours * 60 + startMinutes;
      const endTime = endHours * 60 + endMinutes;

      return this.timeList.map(time => {
        const [hours, minutes] = time.split(':').map(Number);
        const timeInMinutes = hours * 60 + minutes;

        return timeInMinutes >= startTime && timeInMinutes < endTime ? eventName : '';
      });
    },
    queryUsers(users) {
      this.users = users
    },
    updateAllRooms(variables) {

      if (Object.keys(variables).length !== 0) {
        this.dayTime = variables.dayTime;
        this.updateVariables = variables;
      }

      this.loadPreLoader(1000).then(() => {
        this.$refs.commWithGql.fetchAvailableRooms(this.updateVariables);
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
    getHours(timestamp) {
      const date = new Date(timestamp);
      const day = date.getDate();
      const hours = date.getHours();
      const minutes = date.getMinutes();

      return { day, hours, minutes}
    },
    findUserNameById(inputId) {
      const user = this.users.find(user => user.id === inputId);
      return user ? user.name : 'User not found';
    },
    getDaytime() {
      const today = new Date();
      const year = today.getFullYear();
      const month = String(today.getMonth() + 1).padStart(2, '0'); // 月份從0開始，需要+1，並且補零
      const day = String(today.getDate()).padStart(2, '0'); // 日期補零
      return `${year}-${month}-${day}`;
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