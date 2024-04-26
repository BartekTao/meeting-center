<template>
  <comm-with-gql @get-all-rooms="getAllRooms" ref="commWithGql"></comm-with-gql>
  <div class="recent-listing" id="items">
    <div class="container">
      <div class="row">
        <div class="col-lg-12">
          <div class="section-heading">
            <button @click="openForm('none')">Open Form</button>
          </div>
        </div>
        <div class="col-lg-12">
          <div class="">
            <div class="item">
              <div class="row">
                <div class="col-lg-12" v-for="item in test_items" :key="item.id">
                  <div class="listing-item">
                    <div class="left-image">
                      <a href="#"><img :src="image_url" :alt="item.name"></a>
                    </div>
                    <div class="right-content align-self-center">
                      <!-- capacity equipment id roomId rules -->
                      <a href="#"><h4>會議名稱：{{ item.roomId }}</h4></a>
                      <ul class="info">
                        <li>人數限制：{{ item.capacity }}</li>  
                        <li>有大桌子：{{ item.equipment.includes('big table') ? '是' : '否' }}</li>
                        <li>有投影機：{{ item.equipment.includes('projector') ? '是' : '否' }}</li>
                        <li>可否進食：{{ item.rules.includes('no food') ? '否' : '是' }}</li>
                        <li>可否喝水：{{ item.rules.includes('no drinks') ? '否' : '是' }}</li>
                      </ul><br>
                      <div class="flex-container">
                        <div class="main-white-button">
                          <a @click="openForm(item)"><img :src="docImage" alt="Edit">編輯</a>
                        </div>
                        <div class="main-white-button">
                          <a @click="deleteItem('6629c2edd7d285f521a5d787')"><img :src="deleteImage" alt="Delete">刪除</a>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
  
  <script>
  import CommWithGql from '@/components/CommWithGql.vue'
  
  export default {
    components: {
      CommWithGql
    },
    mounted() {
      this.$refs.commWithGql.getAllRooms();
    },
    data() {
      return {
        docImage: require('@/assets/images/google-docs.png'),
        deleteImage: require('@/assets/images/delete.png'),
        image_url: require('../assets/images/listing-01.jpg'),
        test_items: [],
      }
    },
    props: ['items'],
    emits: ['open-form', 'delete-item'],
    methods: {
      openForm(item) {
        this.$emit('open-form', item);
      },
      deleteItem(targetIndex) {
        this.$refs.commWithGql.deleteRoom(targetIndex);
      },
      getAllRooms(rooms) {
        this.test_items = rooms
      }
    },
  };
  </script>
  
  