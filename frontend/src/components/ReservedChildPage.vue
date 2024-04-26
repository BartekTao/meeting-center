<template>
    <div class="recent-listing" id="items">
      <div class="container">
        <div class="row">
          <div class="col-lg-12">
            <div class="section-heading">
              <h2>已預約空間</h2>
            </div>
          </div>
          <div class="col-lg-12">
            <div class="">
              <div class="item">
                <div class="row">
                  <div class="col-lg-12" v-for="item in items" :key="item.name">
                    <div class="listing-item">
                      <div class="left-image">
                        <a href="#"><img :src="item.image_url" :alt="item.name"></a>
                      </div>
                      <div class="right-content align-self-center">
                        <a href="#"><h4>會議名稱：{{ item.name }}</h4></a>
                        <ul class="info">
                          <li>人數限制：{{ item.people_limit }}</li>
                          <li>可否飲食：{{ item.can_eat ? '是' : '否' }}</li>
                        </ul><br>
                        <div class="flex-container">
                          <div class="main-white-button">
                              <a class="openFormBtn" @click="openForm"><img :src="docImage" alt="Edit">編輯</a>
                          </div>
                          <div class="main-white-button">
                              <a class="openFormBtn"><img :src="deleteImage" alt="Delete">刪除</a>
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
          
          <label for="email">參與人信箱：</label>
          <input type="email" id="email" name="email" v-model="formInfo.email" required><br><br>
          
          <label for="period">選擇時段：</label>
          <select id="period" name="period" v-model="formInfo.period" required>
            <option value="13:00">13:00</option>
            <option value="14:00">14:00</option>
            <option value="15:00">15:00</option>
          </select><br><br>
          
          <label for="content">會議內容：</label><br>
          <textarea id="content" name="content" maxlength="200" v-model="formInfo.content"></textarea><br><br>
          
          <label for="file">上傳檔案：</label>
          <input type="file" id="file" name="file"><br><br>
          
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
        formInfo: {
          name: '',
          email: '',
          period: '',
          content: '',
          file: ''
        },
        docImage: require('@/assets/images/google-docs.png'),
        deleteImage: require('@/assets/images/delete.png'),
        items: [
          {
            name: '001',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 12,
            can_eat: false
          },
          {
            name: '002',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 24,
            can_eat: false
          },
          {
            name: '003',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 32,
            can_eat: true
          },
        ]
      }
    },
    methods: {
      openForm() {
        this.formDisplay = true;
      },
      closeForm() {
        this.formDisplay = false;
      },
      submitForm() {
        this.formDisplay = false;
      }
    }
  };
  </script>
  