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
                          <li>有大桌子：{{ item.has_projector ? '是' : '否' }}</li>
                          <li>有投影機：{{ item.has_big_table ? '是' : '否' }}</li>
                        </ul><br>
                        <div class="flex-container">
                          <div class="main-white-button">
                            <a @click="getItem(item.index)"><img :src="docImage" alt="Edit">編輯</a>
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
      <div id="form-popup" v-if="formDisplay">
        <form id="submitForm">
          <label for="name">會議名稱：</label>
          <input type="text" id="name" name="name" v-model="formInfo.name" required><br><br>
          
          <label for="name">圖片網址：</label>
          <input type="text" id="image_url" name="image_url" v-model="formInfo.image_url" required><br><br>
  
          <label for="people_limit">人數限制：</label>
          <input type="number" id="people_limit" name="people_limit" v-model.number="formInfo.people_limit" required><br><br>
      
          <label for="can_eat">可否進食：</label>
          <input type="checkbox" id="can_eat" name="can_eat" v-model="formInfo.can_eat"><br><br>
  
          <label for="can_drink">可否喝水：</label>
          <input type="checkbox" id="can_drink" name="can_drink" v-model="formInfo.can_drink"><br><br>
  
          <label for="has_big_table">有大桌子：</label>
          <input type="checkbox" id="has_big_table" name="has_big_table" v-model="formInfo.can_drink"><br><br>
  
          <label for="has_projector">有投影機：</label>
          <input type="checkbox" id="has_projector" name="has_projector" v-model="formInfo.can_drink"><br><br>
  
          <button type="submit" @click.prevent="submitForm">預約</button>
          <button type="button" @click="closeForm">取消</button>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  
  export default {
    data() {
      return {
        formDisplay: false,
        itemsIndex: 3,
        editIndex: 'none',
        docImage: require('@/assets/images/google-docs.png'),
        deleteImage: require('@/assets/images/delete.png'),
        initialFormInfo: {
          name: '',
          image_url: require('../assets/images/listing-01.jpg'),
          people_limit: 2,
          can_eat: true,
          can_drink: true,
          has_projector: true,
          has_big_table: true
        },
        formInfo: {
          name: '',
          image_url: require('../assets/images/listing-01.jpg'),
          people_limit: 2,
          can_eat: true,
          can_drink: true,
          has_projector: true,
          has_big_table: true
        },
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
        ]
      }
    },
    methods: {
      openForm() {
        this.formDisplay = true
      },
      closeForm() {
        this.formDisplay = false
        this.formReset()
      },
      submitForm() {
        if (this.editIndex === 'none') {
          this.addForm()
        } else {
          this.overwriteItem()
        }
        this.closeForm()
      },
      addForm() {
        this.itemsIndex += 1
        this.formInfo.index = this.itemsIndex.toString()
        this.items.push({...this.formInfo})
      },
      deleteItem(targetIndex) {
        const itemIndex = this.items.findIndex(item => item.index === targetIndex)
        if (itemIndex !== -1) {
          this.items.splice(itemIndex, 1)
        }
      },
      formReset() {
        this.formInfo = {...this.initialFormInfo}
        this.editIndex = 'none'
      },
      getItem(targetIndex) {
        const item = this.items.find(i => i.index === targetIndex)
        if (item) {
          this.formInfo = {...item}
          this.editIndex = targetIndex
        }
      },
      overwriteItem() {
        const itemIndex = this.items.findIndex(i => i.index === this.editIndex)
        if (itemIndex !== -1) {
          Object.assign(this.items[itemIndex], this.formInfo)
        }
      }
    }
  };
  </script>
  
  