<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <ReserveList @showDiv="showDiv" @hideDiv="hideDiv" :openForm="openForm" :openCommentForm="openCommentForm" @update-form="updateAllRooms" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction" :roomItems="roomItems"/>
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
// import ReservedChildPage from '@/components/ReservedChildPage.vue';

export default {
  components: {
    HeadPage,
    ReserveList,
    ReserveForm,
    EventInfo,
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
      console.log('item:', item);
      this.formDisplay = true;
      this.formInfo.roomName = item.name;
      this.roomId = item.id;
      this.formInfo.roomId = item.id;
      this.schedulesList = item.schedulesList;

      this.formInfo.title = item.title;
      this.formInfo.description = item.description;
      this.formInfo.summary = item.summary;

      this.formInfo.fileName = item.fileName;
      this.formInfo.fileUrl = item.fileUrl;

      const { hours: startHours, minutes: startMinutes } = this.getHours(item.startAt);
      const { hours: endHours, minutes: endMinutes } = this.getHours(item.endAt);

      this.formInfo.start_time = `${startHours}:${startMinutes.toString().padStart(2, '0')}`;
      this.formInfo.end_time = `${endHours}:${endMinutes.toString().padStart(2, '0')}`;

      const namesArray = item.participants.map(participant => participant.name);
      this.formInfo.namesString = namesArray.join(', ');
      this.formInfo.eventId = item.eventId;
      

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
    },
    closeCommentForm() {
      this.commentDisplay = false;
    },
    queryUsers(users) {
      this.users = users
    },
    updateAllRooms() {
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
            this.roomItems.push({
              ...room,
              eventId: event.originalData.eventId,
              title: event.originalData.title,
              description: event.originalData.description,
              fileName: event.originalData.fileName,
              fileUrl: event.originalData.fileUrl,
              summary: event.originalData.summary,
              startAt: event.originalData.startAt,
              endAt: event.originalData.endAt,
              participants: event.originalData.participants,
              roomId: event.originalData.roomId
            });
          });
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
    this.updateAllRooms();
  }
}
</script>
