<template>
    <div class="main-banner" id="query">
      <div class="container">
        <div class="row">
          <div class="col-lg-12">
            <div class="top-text header-text">
              <h2>會議預約</h2>
            </div>
          </div>
          <div class="col-lg-12">
            <form id="search-form" name="gs" method="submit" role="search" action="#">
              <div class="row">
                <div class="col-lg-3 align-self-center">
                  <fieldset>
                    <input type="number" id="number_of_people" name="number_of_people" v-model="selectedNumberOfPeople" min="1" placeholder="最多人數" aria-label="Number of People" required>
                  </fieldset>
                </div>
                <div class="col-lg-3 align-self-center">
                  <fieldset>
                    <input type="date" name="date" class="form-control" v-model="currentDate" required>
                  </fieldset>
                </div>        
                <div class="col-lg-3 align-self-center">
                  <fieldset>
                    <select name="start_period" class="form-select" v-model="selectedStartPeriod" placeholder="開始時段" aria-label="Number of People" id="start_period" required>
                      <option value="" disabled>開始時間</option>
                      <option :value="item" v-for="item in time_period" :key="item">{{ item }}</option>
                    </select>
                  </fieldset>
                </div>
                <div class="col-lg-3 align-self-center">
                  <fieldset>
                    <select name="end_period" class="form-select" v-model="selectedEndPeriod" placeholder="開始時段" aria-label="Number of People" id="end_period" required>
                      <option value="" disabled>結束時間</option>
                      <option :value="item" v-for="item in time_period" :key="item">{{ item }}</option>
                    </select>
                  </fieldset>
                </div>
                <div class="col-lg-3 align-self-center">
                  <fieldset>
                    <input type="number" name="time_period" v-model="selectedTimePeriod" placeholder="連續時間(小時)" id="time_period" required min="0" max="9" step="0.5">
                  </fieldset>
                </div>
                <div class="col-lg-3">                        
                  <fieldset>
                    <button class="main-button" @click.prevent="fetchAvailableRooms"><i class="fa fa-search"></i> 查詢</button>
                  </fieldset>
                </div>
              </div>
            </form>
          </div>
          <div class="col-lg-10 offset-lg-1">
            <ul class="categories">
              <li><a><span class="icon" :class="{ checked: mustEat }" @click="changeClickSpan('mustEat')"><img :src="projectorImage"></span> 有投影機</a></li>
              <li><a><span class="icon" :class="{ checked: hasProjector }" @click="changeClickSpan('hasProjector')"><img :src="eatImage"></span> 可否飲食 </a></li>
              <li><a><span class="icon" :class="{ checked: hasWhiteBoard }" @click="changeClickSpan('hasWhiteBoard')"><img :src="whiteBoardImage"></span> 有無白板</a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>

  export default {
    name: 'ReserveBar',
    data() {
      return {
        time_period: [],
        selectedNumberOfPeople: 5,
        currentDate: this.getCurrentDate(), 
        selectedStartPeriod: '9:00',
        selectedEndPeriod: '18:00',
        selectedTimePeriod: '1',
        projectorImage: require('@/assets/images/projector.png'),
        eatImage: require('@/assets/images/fast-food.png'),
        whiteBoardImage: require('@/assets/images/whiteboard.png'),
        mustEat: false,
        hasProjector: false,
        hasWhiteBoard: false,
        first_data: null
      };
    },
    methods: {
      fetchAvailableRooms() {
      },
      changeClickSpan(key) {
        this[key] = !this[key];
      },
      getCurrentDate() {
        const today = new Date();
        const dd = String(today.getDate()).padStart(2, '0');
        const mm = String(today.getMonth() + 1).padStart(2, '0');
        const yyyy = today.getFullYear();
        return `${yyyy}-${mm}-${dd}`;
      }
    },
    mounted() {
      this.time_period = this.$names;
    },
  };
  </script>
  
  