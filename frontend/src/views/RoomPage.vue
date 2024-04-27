<template>
    <div>
      <head-page pageContent="後台管理"></head-page>
      <room-list @open-form="openForm"></room-list>
      <room-edit-form :formDisplay="formDisplay" :roomInfo="roomInfo" @close-form="closeForm" @update-info="handleUpdate"></room-edit-form>
      <comm-with-gql ref="commWithGql"></comm-with-gql>
    </div>
</template>
  
  <script>
  import HeadPage from '@/components/HeadPage.vue';
  import RoomList from '@/components/RoomList.vue';
  import CommWithGql from '@/components/CommWithGql.vue'
  import RoomEditForm from '@/components/RoomEditForm.vue';
  
  export default {
    components: {
      HeadPage,
      RoomList,
      CommWithGql,
      RoomEditForm
    },
    data() {
      return {
        formDisplay: false,
        roomInfo: {},
        init_room: {
            roomId: "test",
            capacity: 10,
            equipment: [],
            rules: ["no food", "no drinks" ]
        }
      };
    },  
    methods: {
      openForm(item) {
        this.formDisplay = true;
        
        if (item === 'none') {
          this.roomInfo = { ...this.init_room };
        } else {
          this.roomInfo = {
            id: '6629c2edd7d285f521a5d787',
            roomId: item.roomId,
            capacity: item.capacity,
            equipment: item.equipment,
            rules: item.rules
          };
        }
      },
      closeForm() {
        this.formDisplay = false;
      },
      handleUpdate({ field, value }) {
        this.$set(this.roomInfo, field, value);
      },
    },
  }
  </script>
  