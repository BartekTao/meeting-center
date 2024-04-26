<template>
  <div class="recent-listing" id="items">
    <div class="container">
      <div class="row">
        <div class="col-lg-12">
          <div class="section-heading">
            <button @click="openForm">新增會議室</button>
          </div>
        </div>
        <div class="col-lg-12">
          <div class="">
            <div class="item">
              <div class="row">
                <div class="col-lg-12" v-for="item in items" :key="item.index">
                  <div class="listing-item">
                    <div class="left-image">
                      <a href="#"><img :src="item.image_url" :alt="item.name"></a>
                    </div>
                    <div class="right-content align-self-center">
                      <a href="#"><h4>會議名稱：{{ item.name }}</h4></a>
                      <ul class="info">
                        <li>人數限制：{{ item.people_limit }}</li>
                        <li>可否進食：{{ item.can_eat ? '是' : '否' }}</li>
                        <li>可否喝水：{{ item.can_drink ? '是' : '否' }}</li>
                        <li>有大桌子：{{ item.has_big_table ? '是' : '否' }}</li>
                        <li>有投影機：{{ item.has_projector ? '是' : '否' }}</li>
                      </ul><br>
                      <div class="flex-container">
                        <div class="main-white-button">
                          <a @click="editItem(item.index)"><img :src="editImage" alt="Edit">編輯</a>
                        </div>
                        <div class="main-white-button">
                          <a @click="deleteItem(item.index)"><img :src="deleteImage" alt="Delete">刪除</a>
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
export default {
  props: ['items'],
  emits: ['editItem', 'deleteItem'],
  data() {
    return {
      editImage: require('@/assets/images/edit.png'),
      deleteImage: require('@/assets/images/delete.png')
    };
  },
  methods: {
    editItem(index) {
      this.$emit('editItem', index);
    },
    deleteItem(index) {
      this.$emit('deleteItem', index);
    },
    openForm() {
      this.$emit('openForm');
    }
  }
};
</script>
