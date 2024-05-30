<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <ReserveList :openForm="openForm" :openCommentForm="openCommentForm" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction" :roomItems="roomItems"/>
    <ReserveForm @close-form="closeForm" :formDisplay="formDisplay" :roomInfo="roomInfo" :userName="userName"/>
    <CommentForm @close-comment-form="closeCommentForm" :commentDisplay="commentDisplay" :roomInfo="roomInfo" :userName="userName"/>
    <comm-with-gql @fetch-available-rooms="fetchAvailableRooms" @get-event-list="getEventList" @query-users="queryUsers" ref="commWithGql"></comm-with-gql>
    <js-preloader ref="jsPreloader"></js-preloader>

  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';
import CommentForm from '@/components/CommentForm.vue';
import CommWithGql from '@/components/CommWithGql.vue';
import JsPreloader from '@/components/JsPreloader.vue';
// import ReservedChildPage from '@/components/ReservedChildPage.vue';

export default {
  components: {
    HeadPage,
    ReserveList,
    ReserveForm,
    CommentForm,
    CommWithGql,
    JsPreloader
  },
  data() {
    return {
      bookingAction: false,
      editAction: true,
      editCommentAction: true,
      deleteAction: true,

      formDisplay: false,
      commentDisplay: false,
      ignoreRoom: [],

      roomItems: [],
      users: [],
      eventList: [],
      updateVariables: {
        userIDs: ['6645ece136e2a0f035961bdd'],
        startAt: this.getCurrentDateStartTimeStamp(),
        endAt: 1748599200000,
      },

      timeList: ['9:0', '9:30', '10:0', '10:30', '11:0', '11:30', '12:0', '12:30', '13:0', '13:30', '14:0', '14:30', '15:0', '15:30', '16:0', '16:30', '17:0', '17:30', '18:0'],
      reservatorList: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
    
      roomInfo: {
          reservatorList: [],
          id: '',
          roomName: '',
      },
      userName: "Ray",
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
    openCommentForm(item) {
      this.commentDisplay = true;
      this.roomInfo.id = item.id;
      this.roomInfo.roomName = item.roomId;
    },
    closeCommentForm() {
      this.commentDisplay = false;
    },
    queryUsers(users) {
      this.users = users
    },
    updateAllRooms() {

      this.loadPreLoader(100).then(() => {
        this.$refs.commWithGql.getUserEvents(this.updateVariables);
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
    getCurrentDateStartTimeStamp() {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      return today.getTime();
    },
    getEventList(eventList) {
      
      this.eventList = eventList;

      eventList.forEach(event => {
        if (Array.isArray(event.roomsData)) {
          this.roomItems.push(...event.roomsData);
        }
      });


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
    fetchAvailableRooms(room) {
      this.ignoreRoom = room;
    },
    findUserNameById(inputId) {
      const user = this.users.find(user => user.id === inputId);
      return user ? user.name : 'User not found';
    },
    getHours(timestamp) {
      const date = new Date(timestamp);
      const day = date.getDate();
      const hours = date.getHours();
      const minutes = date.getMinutes();

      return { day, hours, minutes}
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
  },
  mounted() {
    this.$refs.commWithGql.queryUsers();
    this.updateAllRooms();
  }
}
</script>
