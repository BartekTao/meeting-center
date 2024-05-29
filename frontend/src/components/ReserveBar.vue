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
                    <button class="main-button" @click.prevent="updateAllRooms"><i class="fa fa-search"></i> 查詢</button>
                  </fieldset>
                </div>
              </div>
            </form>
          </div>
          <div class="col-lg-10 offset-lg-1">
            <ul class="categories">
              <li><a><span class="icon" :class="{ checked: PROJECTOR }" @click="changeClickSpan('PROJECTOR')"><img :src="projectorImage"></span> 有投影機</a></li>
              <li><a><span class="icon" :class="{ checked: NO_FOOD }" @click="changeClickSpan('NO_FOOD')"><img :src="eatImage"></span> 可否飲食 </a></li>
              <li><a><span class="icon" :class="{ checked: TABLE }" @click="changeClickSpan('TABLE')"><img :src="whiteBoardImage"></span> 有無桌子</a></li>
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
        selectedStartPeriod: '9:00',
        selectedEndPeriod: '18:00',
        selectedNumberOfPeople: 5,
        currentDate: this.getCurrentDate(), 
        selectedTimePeriod: '1',
        projectorImage: require('@/assets/images/projector.png'),
        eatImage: require('@/assets/images/fast-food.png'),
        whiteBoardImage: require('@/assets/images/whiteboard.png'),
        NO_FOOD: false,
        PROJECTOR: false,
        TABLE: false,
        first_data: null,
        variables : {
          dayTime: '',
          startAt: 1625077800,
          endAt: 1625081400,
          rules: [],
          equipments: [],
          first: 20,
          after: null
        },
      };
    },
    methods: {
      updateAllRooms() {
        this.variables.rules = [];
        this.variables.equipments = [];

        const dayTime = this.currentDate;
        
        const startTime = dayTime+'-'+this.selectedStartPeriod + ':00';
        const startTimeStamp = this.transferToTimestamp(startTime);

        const endTime = dayTime+'-'+this.selectedEndPeriod + ':00';
        const endTimeStamp = this.transferToTimestamp(endTime);

        if (this.NO_FOOD && !this.variables.rules.includes('NO_FOOD')) {
          this.variables.rules.push('NO_FOOD');}
        if (this.PROJECTOR && !this.variables.equipments.includes('PROJECTOR')) {
          this.variables.equipments.push('PROJECTOR');}
        if (this.TABLE && !this.variables.equipments.includes('TABLE')) {
          this.variables.equipments.push('TABLE');}

        this.variables.dayTime = dayTime;
        this.variables.startAt = startTimeStamp;
        this.variables.endAt = endTimeStamp;
        this.$emit('updateAllRooms', this.variables);
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
      },
      transferToTimestamp(time) {
        const formattedTime_ = time.replace(/(\d{4})-(\d{2})-(\d{2})-(\d{1,2}):(\d{2}):(\d{2})/, '$1-$2-$3T$4:$5:$6');
        const formattedTime = formattedTime_.replace(/T(\d):/, 'T0$1:');
        const date = new Date(formattedTime);
        return date.getTime();
      }
    },
    mounted() {
      this.time_period = this.$names;
    },
    emits: ['updateAllRooms'],
  };
  </script>
  
  