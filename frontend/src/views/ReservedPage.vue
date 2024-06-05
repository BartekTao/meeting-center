<template>
  <div>
    <background-header></background-header>
    <head-page pageContent="已預約空間"></head-page>
    <ReserveList @showDiv="showDiv" @hideDiv="hideDiv" :openForm="openForm" :openCommentForm="openCommentForm" @update-form="updateAllRooms" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction" :roomItems="roomItems" :pageState="pageState"/>
    <ReserveForm @showDiv="showDiv" @hideDiv="hideDiv" :formDisplay="formDisplay" :formInfo="formInfo" :schedulesList="schedulesList" @close-form="closeForm" @update-form="updateAllRooms" :users="users"/>
    <CommentForm @close-comment-form="closeCommentForm" @update-form="updateAllRooms" :commentDisplay="commentDisplay" :formInfo="formInfo" :users="users"/>
    <EventInfo ref="eventInfo"/>
    <comm-with-gql @fetch-available-rooms="fetchAvailableRooms" @get-event-list="getEventList" @query-users="queryUsers" ref="commWithGql"></comm-with-gql>
    <js-preloader ref="jsPreloader"></js-preloader>
  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';
import EventInfo from '@/components/EventInfo.vue';
import CommentForm from '@/components/CommentForm.vue';
import CommWithGql from '@/components/CommWithGql.vue';
import JsPreloader from '@/components/JsPreloader.vue';
import BackgroundHeader from '@/components/BackgroundHeader.vue';
// import ReservedChildPage from '@/components/ReservedChildPage.vue';

export default {
  components: {
    HeadPage,
    ReserveList,
    ReserveForm,
    EventInfo,
    CommentForm,
    CommWithGql,
    JsPreloader,
    BackgroundHeader,
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
      pageState: 'reserved',
      
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
        notes: '',
        summary: '',
        fileName: '',
        fileUrl: '',
        reservatorList: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ],
        pageState: 'reserved',
      },

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
      roomItems: [],
      roomName: "",
      roomId: "",
      schedulesList: [],
    };
  },
  methods: {
    openForm(item) {
      this.formDisplay = true;
      this.formInfo.roomName = item.name;
      this.roomId = item.id;
      this.formInfo.roomId = item.id;
      this.schedulesList = item.schedulesList;

      this.formInfo.title = item.title;
      this.formInfo.description = item.description;
      this.formInfo.summary = item.summary;

      this.formInfo.start_time = item.start_time;
      this.formInfo.end_time = item.end_time;

      const namesArray = item.participants.map(participant => participant.name);
      this.formInfo.namesString = namesArray.join(', ');
      this.formInfo.eventId = item.eventId;
      this.formInfo.fileName = item.fileName;
      this.formInfo.fileUrl = item.fileUrl;

    },
    closeForm() {
      this.formDisplay = false;
    },
    openCommentForm(item) {
      this.commentDisplay = true;
      this.roomInfo.id = item.id;
      this.roomInfo.roomName = item.name;
      this.formInfo.eventId = item.eventId;
      this.formInfo.summary = item.summary;
      this.formInfo.title = item.title;
    },
    closeCommentForm() {
      this.commentDisplay = false;
    },
    queryUsers(users) {
      this.users = users
    },
    updateAllRooms() {
      console.log('this.updateVariables:', this.updateVariables);
      this.loadPreLoader(500).then(() => {
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
    showDiv(data) {
      this.$refs.eventInfo.showDiv(data);
    },
    hideDiv() {
      this.$refs.eventInfo.hideDiv();
    },
    getEventList(eventList) {
      
      this.roomItems = []; // clear
      this.eventList = eventList;

      // Iterate through eventList and extract roomsData along with originalData
      eventList.forEach(event => {
        if (Array.isArray(event.roomsData)) {

          event.roomsData.forEach(room => {

            const {month_day: eventDay, hours: startHours, minutes: startMinutes } = this.getHours(event.originalData.startAt);
            const { hours: endHours, minutes: endMinutes } = this.getHours(event.originalData.endAt);

            const start_time = `${startHours}:${startMinutes.toString().padStart(2, '0')}`;
            const end_time = `${endHours}:${endMinutes.toString().padStart(2, '0')}`;
            
            this.roomItems.push({
              ...room,
              eventId: event.originalData.eventId,
              title: event.originalData.title,
              description: event.originalData.description,
              fileName: event.originalData.fileName,
              fileUrl: event.originalData.fileUrl,
              summary: event.originalData.summary,
              startAt: event.originalData.startAt,
              start_time: start_time,
              end_time: end_time,
              eventDay: eventDay,
              endAt: event.originalData.endAt,
              participants: event.originalData.participants,
              roomId: event.originalData.roomId,
              creatorId: event.originalData.creatorId,
            });
          });
        }
      });

      this.roomItems.forEach(room => {
        let finalReservatiorList = Array(this.timeList.length).fill('');

        room.schedules.forEach(schedule => {
          // const reservatorName = this.findUserNameById(schedule.creator.id);
          const reservatorName = schedule.participants[0].name;
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
              // id: schedule.creator.id,
              id: schedule.participants[0].id,
              name: reservatorName,
              nickName: nickName,
              title: schedule.title,
              startHours,
              startMinutes,
              endHours,
              endMinutes,
              eventTitle,
              state: '',
            };
          finalReservatiorList = finalReservatiorList.map((slot, index) => slot.name === reservatorName ? slot : newReservatorList[index] ? scheduleInfo : slot);
          
        });

        room.schedulesList = finalReservatiorList.slice(0, -1);

        const { hours: startHours, minutes: startMinutes } = this.getHours(room.startAt);
        const { hours: endHours, minutes: endMinutes } = this.getHours(room.endAt);

        const reservatorName = this.findUserNameById(room.creatorId);
        room.currentSchedulesList = this.transferReservatorList(startHours, startMinutes, endHours, endMinutes, reservatorName).slice(0, -1);

        for (let i = 0; i < room.schedulesList.length; i++) {
          const scheduleUnit = room.schedulesList[i];
          const currentScheduleUnit = room.currentSchedulesList[i];

          if (scheduleUnit && currentScheduleUnit) {
            room.schedulesList[i].state = 'current';
          } else if (scheduleUnit && !currentScheduleUnit) {
            room.schedulesList[i].state = 'occupied';
          } 
        }
        
      });
    console.log('roomItems:', this.roomItems);
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
      const month = date.getMonth() + 1; // getMonth() returns 0-11, so we add 1
      const day = date.getDate();
      const month_day = `${month}-${day}`;
      const hours = date.getHours();
      const minutes = date.getMinutes();

      return { month_day, hours, minutes}
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
    getDaytime() {
      const today = new Date();
      const year = today.getFullYear();
      const month = String(today.getMonth() + 1).padStart(2, '0'); // 月份從0開始，需要+1，並且補零
      const day = String(today.getDate()).padStart(2, '0'); // 日期補零
      return `${year}-${month}-${day}`;
    },
  },
  mounted() {
    // this.$refs.commWithGql.queryUsers();
    this.updateAllRooms();
    this.$refs.commWithGql.queryUsers().then(() => {
      this.formInfo.userId = [this.users[0].id];
      this.updateVariables.userIDs = [this.users[0].id];
    }).catch(error => {
      console.error("Failed to fetch users:", error);
    });
  }
}
</script>
