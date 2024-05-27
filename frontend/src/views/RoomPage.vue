<template>
    <div>
      <head-page pageContent="後台管理"></head-page>
      <room-list @open-form="openForm" @update-all-rooms="updateAllRooms" :testItems="testItems"></room-list>
      <room-edit-form :formDisplay="formDisplay" :roomInfo="roomInfo" @close-form="closeForm" @update-info="handleUpdate" @update-all-rooms="updateAllRooms"></room-edit-form>
      <!-- <comm-with-gql ref="commWithGql"></comm-with-gql> -->
      <comm-with-gql @query-all-rooms="queryAllRooms" ref="commWithGql"></comm-with-gql>
      <js-preloader ref="jsPreloader"></js-preloader>
    </div>
</template>
  
  <script>
  import HeadPage from '@/components/HeadPage.vue';
  import RoomList from '@/components/RoomList.vue';
  import CommWithGql from '@/components/CommWithGql.vue'
  import RoomEditForm from '@/components/RoomEditForm.vue';
  import JsPreloader from '@/components/JsPreloader.vue';
  
  export default {
    components: {
      HeadPage,
      RoomList,
      CommWithGql,
      RoomEditForm,
      JsPreloader
    },
    mounted() {
      this.updateAllRooms();

      // this.$nextTick(() => {
      //   this.updateAllRooms();
      // });
    },
    data() {
      return {
        formDisplay: false,
        roomInfo: {},
        init_room: {
            name: "test",
            capacity: 10,
            equipments: [],
            rules: ["NO_FOOD", "NO_DRINK" ]
        },
        testItems: [],
      };
    },  
    methods: {
      openForm(item) {
        this.formDisplay = true;
        
        if (item === 'none') {
          this.roomInfo = { ...this.init_room };
        } else {
          this.roomInfo = {
            id: item.id,
            name: item.name,
            capacity: item.capacity,
            equipments: item.equipments,
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
      queryAllRooms(rooms) {
        this.testItems = rooms
      },
      updateAllRooms() {
        this.loadPreLoader(1000).then(() => {
          this.$refs.commWithGql.queryAllRooms();
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
      }
    },
  }
  </script>
  