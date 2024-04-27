<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <ReserveList :openForm="openForm" :bookingAction="bookingAction" :editAction="editAction" :deleteAction="deleteAction"/>
    <ReserveForm :formDisplay="formDisplay" :roomInfo="roomInfo" :userName="userName" @close-form="closeForm"/>

    <!-- <reserved-child-page></reserved-child-page> -->
  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';
// import ReservedChildPage from '@/components/ReservedChildPage.vue';

export default {
  components: {
    HeadPage,
    ReserveList,
    ReserveForm
    // ReservedChildPage
  },
  data() {
    return {
      bookingAction: false,
      editAction: true,
      deleteAction: true,
      formDisplay: false,
      roomInfo: {
          reservatorList: [],
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
  }
}
</script>
