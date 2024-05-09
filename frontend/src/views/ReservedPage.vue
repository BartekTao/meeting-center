<template>
  <div>
    <head-page pageContent="已預約空間"></head-page>
    <ReserveList :openForm="openForm" :openCommentForm="openCommentForm" :bookingAction="bookingAction" :editAction="editAction" :editCommentAction="editCommentAction" :deleteAction="deleteAction"/>
    <ReserveForm @close-form="closeForm" :formDisplay="formDisplay" :roomInfo="roomInfo" :userName="userName"/>
    <CommentForm @close-comment-form="closeCommentForm" :commentDisplay="commentDisplay" :roomInfo="roomInfo" :userName="userName"/>

  </div>
</template>

<script>
import HeadPage from '@/components/HeadPage.vue';
import ReserveList from '@/components/ReserveList.vue';
import ReserveForm from '@/components/ReserveForm.vue';
import CommentForm from '@/components/CommentForm.vue';
// import ReservedChildPage from '@/components/ReservedChildPage.vue';

export default {
  components: {
    HeadPage,
    ReserveList,
    ReserveForm,
    CommentForm
    // ReservedChildPage
  },
  data() {
    return {
      bookingAction: false,
      editAction: true,
      editCommentAction: true,
      deleteAction: true,

      formDisplay: false,
      commentDisplay: false,

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
  }
}
</script>
