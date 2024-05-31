<template>
  <comm-with-gql ref="commWithGql"></comm-with-gql>
  <div class="recent-listing" id="items">
    <div class="container">
      <div class="row">
        <div class="col-lg-12">
          <div class="section-heading">
            <button @click="openForm('none')">新建會議室</button>
          </div>
        </div>
        <div class="col-lg-12">
          <div class="">
            <div class="item">
              <div class="row">
                <div class="col-lg-12" v-for="item in testItems" :key="item.id">
                  <div class="listing-item">
                    <div class="left-image">
                      <a href="#"><img :src="image_url" :alt="item.name"></a>
                    </div>
                    <div class="right-content align-self-center">
                      <!-- capacity equipment id roomId rules -->
                      <a href="#"><h4>會議名稱：{{ item.name }}</h4></a>
                      <ul class="info">
                        <li>人數限制：{{ item.capacity }}</li>  
                        <li>有大桌子：{{ item.equipments.includes('TABLE') ? '是' : '否' }}</li>
                        <li>有投影機：{{ item.equipments.includes('PROJECTOR') ? '是' : '否' }}</li>
                        <li>可否進食：{{ item.rules.includes('NO_FOOD') ? '否' : '是' }}</li>
                        <li>可否喝水：{{ item.rules.includes('NO_DRINK') ? '否' : '是' }}</li>
                      </ul><br>
                      <div class="flex-container">
                        <div class="main-white-button">
                          <a @click="openForm(item)"><img :src="docImage" alt="Edit">編輯</a>
                        </div>
                        <div class="main-white-button">
                          <a @click="deleteRoom(item.id)"><img :src="deleteImage" alt="Delete">刪除</a>
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
    data() {
      return {
        docImage: require('@/assets/images/google-docs.png'),
        deleteImage: require('@/assets/images/delete.png'),
        image_url: require('../assets/images/listing-01.jpg'),
      }
    },
    emits: ['open-form', 'delete-item', 'update-all-rooms'],
    props: ['testItems'],
    methods: {
      openForm(item) {
        this.$emit('open-form', item);
      },
      deleteRoom(targetIndex) {
        this.$refs.commWithGql.deleteRoom(targetIndex)
        this.$emit('update-all-rooms');
      },
    },
  };
  </script>
  
  