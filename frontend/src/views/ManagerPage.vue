<template>
    <div>
      <head-page pageContent="後台管理"></head-page>
      <manager-item :items="items" @open-form="openForm" @overwriteItem="handleOverwriteItem"></manager-item>
      <room-edit-form :formDisplay="formDisplay" :formInfo="formInfo" @close-form="closeForm" @update-info="handleUpdate"></room-edit-form>
      <comm-with-gql ref="commWithGql"></comm-with-gql>
    </div>
</template>
  
  <script>
  import HeadPage from '@/components/HeadPage.vue';
  import ManagerItem from '@/components/ManagerItem.vue';
  import CommWithGql from '@/components/CommWithGql.vue'
  import RoomEditForm from '@/components/RoomEditForm.vue';
  
  export default {
    components: {
      HeadPage,
      ManagerItem,
      CommWithGql,
      RoomEditForm
    },
    data() {
      return {
        formDisplay: false,
        editIndex: 'none',
        itemsIndex: 0,
        formInfo: {},
        items: [
          { 
            index: '1',
            name: '001',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 12,
            can_eat: false,
            can_drink: false,
            has_projector: true,
            has_big_table: false
          },
          {
            index: '2',
            name: '002',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 24,
            can_eat: false,
            can_drink: false,
            has_projector: true,
            has_big_table: false
          },
          {
            index: '3',
            name: '003',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 32,
            can_eat: true,
            can_drink: false,
            has_projector: true,
            has_big_table: false
          }
        ],
        init_room: {
            roomId: "test",
            capacity: 10,
            equipment: [],
            rules: ["no food", "no drinks" ]
        }
      };
    },  
    methods: {
      handleOverwriteItem(index, newItem) {
        const idx = this.items.findIndex(item => item.index === index);
        if (idx !== -1) {
          Object.assign(this.items[idx], newItem);
        }
      },
      openForm(item) {
        this.formDisplay = true;
        
        if (item === 'none') {
          this.formInfo = { ...this.init_room };
        } else {
          this.formInfo = {
            id: '6629c2edd7d285f521a5d787',
            roomId: this.init_room.roomId,
            capacity: this.init_room.capacity,
            equipment: this.init_room.equipment,
            rules: this.init_room.rules
          };
        }
        // console.log(this.formInfo)
      },
      closeForm() {
        this.formDisplay = false;
      },
      handleUpdate({ field, value }) {
        this.$set(this.formInfo, field, value);
      }
    },
    mounted() {
      this.itemsIndex = this.items.length;
    }
  }
  </script>
  