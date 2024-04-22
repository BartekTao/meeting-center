<template>
    <div class="recent-listing" id="items">
      <div class="container">
        <div class="row">
          <div class="col-lg-12">
            <div class="section-heading">
              <h2>查詢空間</h2>
            </div>
          </div>
          <div class="col-lg-12">
            <div class="">
              <div class="item">
                <div class="row">
                  <EventItem v-for="item in items" :key="item.name" :item="item" @openForm="openForm"/>
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
          
          <label for="email">參與人信箱(如有多個，請在間隔留下,符號)：</label><br>
          <input type="text" id="email" name="email" v-model="formInfo.email" required><br><br>
  
          <label>開始時間：</label>
          <select v-model="formInfo.start_time" required>
              <option :value="item" v-for="item in time_period"  :key="item">{{ item }}</option>
          </select><br><br>
  
          <label>結束時間：</label>
          <select v-model="formInfo.end_time" required>
              <option :value="item" v-for="item in time_period"  :key="item">{{ item }}</option>
          </select><br><br>
  
          <label>會議內容：</label><br>
          <textarea id="content" name="content" maxlength="200" v-model="formInfo.content"></textarea><br><br>
          
          <label>上傳檔案：</label>
          <input type="file" id="file" name="file"><br><br>
          
          <button type="submit" @click.prevent="submitForm">預約</button>
          <button type="button" @click="closeForm">取消</button>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  import EventItem from './EventItem.vue';
  export default {
    name: 'QueryItems',
    data() {
      return {
        time_period: ['9:00', '9:30', '10:00', '10:30', '11:00', '11:30',
          '12:00', '12:30', '13:00', '13:30', '14:00', '14:30', '15:00',
          '15:30', '16:00', '16:30', '17:00', '17:30', '18:00'],
        formDisplay: false,
        tapImage: require('@/assets/images/tap.png'),
        formInfo: {
          name: 'Ray',
          email: 'example@gmail.com',
          start_time: '10:00',
          end_time: '12:00',
          content: 'test content',
          file: ''
        },
        items: [
          {
            name: '001',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 12,
            can_eat: false,
            reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
          {
            name: '002',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 24,
            can_eat: false,
            reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
          {
            name: '003',
            image_url: require('../assets/images/listing-01.jpg'),
            people_limit: 32,
            can_eat: true,
            reservatorList: ['', '', 'Ivan', '', '', 'Ray', '', '', '', 'Ivan', 'Kevin', 'Ray', 'John', '', '', '', 'Ray', ''],
          },
        ]
      };
    },
    components: {
      EventItem
    },
    methods: {
      openForm() {
        this.formDisplay = true;
      },
      closeForm() {
        this.formDisplay = false;
      },
      submitForm() {
        alert(JSON.stringify(this.formInfo, null, 2)); 
        this.closeForm();
      }
    }
  };
  </script>
  
  