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
                    <button class="main-button" @click.prevent="queryData"><i class="fa fa-search"></i> 查詢</button>
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
  import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
  import { BOOKS_QUERY } from './QueryObj.vue';

  export default {
    name: 'QueryBar',
    data() {
      return {
        time_period: ['9:00', '9:30', '10:00', '10:30', '11:00', '11:30', '12:00', '12:30', '13:00', '13:30', '14:00', '14:30', '15:00', '15:30', '16:00', '16:30', '17:00', '17:30', '18:00'],
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
      queryData() {
        alert(JSON.stringify({
          number_of_people: this.selectedNumberOfPeople,
          date: this.currentDate,
          start_period: this.selectedStartPeriod,
          end_period: this.selectedEndPeriod,
          time_period: this.selectedTimePeriod,
          must_eat: this.mustEat,
          has_projector: this.hasProjector,
          has_whiteboard: this.hasWhiteBoard,
        }, null, 2));
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
      fetchBooks() {
        const httpLink = createHttpLink({
          uri: 'http://localhost:4000/' // 替换为您的 GraphQL API URI
        });

        const client = new ApolloClient({
          link: httpLink,
          cache: new InMemoryCache()
        });

        client.query({
          query: BOOKS_QUERY
        }).then(result => {
          console.log(result.data);
          if (result.data && result.data.books && result.data.books.length > 0) {
            this.first_data = result.data.books[0];  // 更新 first_data 为查询结果的第一条数据
          }
        }).catch(error => {
          console.error("Error fetching the books:", error);
        });
      }
    },
    mounted() {
      this.fetchBooks();
    },
  };
  </script>
  
  