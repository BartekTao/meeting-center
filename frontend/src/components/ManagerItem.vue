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
  <!-- <div id="form-popup" v-if="formDisplay" class="container mt-3">
      <form id="submitForm">
        <div class="row mb-2">
          <label for="name" class="col-sm-2 col-form-label">會議名稱：</label>
          <div class="col-sm-9">
            <input type="text" id="name" name="name" class="form-control" v-model="formInfo.name" required>
          </div>
        </div>

        <div class="row mb-2">
          <label for="image_url" class="col-sm-2 col-form-label">圖片網址：</label>
          <div class="col-sm-9">
            <input type="text" id="image_url" name="image_url" class="form-control" v-model="formInfo.image_url" required>
          </div>
        </div>

        <div class="row mb-2">
          <label for="people_limit" class="col-sm-2 col-form-label">人數限制：

          </label>
          <div class="col-sm-9">
            <input type="number" id="people_limit" name="people_limit" class="form-control" v-model.number="formInfo.people_limit" required>
          </div>
        </div>

        <div class="row mb-2">
          <label for="can_eat" class="col-sm-2 col-form-label">可否進食：</label>
          <div class="col-sm-9">
            <input type="checkbox" id="can_eat" name="can_eat" class="form-check-input mt-2" v-model="formInfo.can_eat">
          </div>
        </div>

        <div class="row mb-2">
          <label for="can_drink" class="col-sm-2 col-form-label">可否喝水：</label>
          <div class="col-sm-9">
            <input type="checkbox" id="can_drink" name="can_drink" class="form-check-input mt-2" v-model="formInfo.can_drink">
          </div>
        </div>

        <div class="row mb-2">
          <label for="has_big_table" class="col-sm-2 col-form-label">有大桌子：</label>
          <div class="col-sm-9">
            <input type="checkbox" id="has_big_table" name="has_big_table" class="form-check-input mt-2" v-model="formInfo.has_big_table">
          </div>
        </div>

        <div class="row mb-2">
          <label for="has_projector" class="col-sm-2 col-form-label">有投影機：</label>
          <div class="col-sm-9">
            <input type="checkbox" id="has_projector" name="has_projector" class="form-check-input mt-2" v-model="formInfo.has_projector">
          </div>
        </div>

        <div class="row mb-2">
          <div class="col-sm-12 d-flex justify-content-center">
            <button type="submit" class="btn btn-primary margin-right-2cm" @click.prevent="submitForm">編輯</button>
            <button type="button" class="btn btn-secondary" @click="closeForm">取消</button>
          </div>
        </div>
      </form>
  </div> -->
</template>
  
  <script>
  import CommWithGql from '@/components/CommWithGql.vue'
  
  export default {
    components: {
      CommWithGql
    },
    data() {
      return {
        formDisplay: false,
        editIndex: 'none',
        docImage: require('@/assets/images/google-docs.png'),
        deleteImage: require('@/assets/images/delete.png'),
        image_url: require('../assets/images/listing-01.jpg'),
        initialFormInfo: {
          name: '',
          image_url: require('../assets/images/listing-01.jpg'),
          people_limit: 2,
          can_eat: true,
          can_drink: true,
          has_projector: true,
          has_big_table: true
        },
        test_items: [],
        formInfo: {
          name: '',
          image_url: require('../assets/images/listing-01.jpg'),
          people_limit: 2,
          can_eat: true,
          can_drink: true,
          has_projector: true,
          has_big_table: true
        },
      }
    },
    props: ['items'],
    emits: ['open-form', 'delete-item', 'overwrite-item'],
    methods: {
      openForm(item) {
        this.$emit('open-form', item);
      },
      submitForm() {
        if (this.editIndex === 'none') {
          this.addForm()
        } else {
          this.overwriteItem()
        }
        this.closeForm();
        this.formReset()
      },
      addForm() {
        this.itemsIndex += 1;
        this.formInfo.index = this.itemsIndex.toString();
        this.$emit('addItem', {...this.formInfo});
      },
      deleteItem(targetIndex) {
        this.$refs.commWithGql.deleteRoom(targetIndex);
      },
      formReset() {
        this.formInfo = {...this.initialFormInfo};
        this.editIndex = 'none';
      },
      EditRoom(targetIndex) {
        const item = this.items.find(i => i.index === targetIndex);
        if (item) {
          this.formInfo = {...item};
          this.openForm();
          this.editIndex = targetIndex;
            }
          },
      overwriteItem() {
        this.$emit('overwriteItem', this.editIndex, {...this.formInfo});
      },
      getAllRooms(rooms) {
        this.test_items = rooms
      }
    },
  };
  </script>
  
  