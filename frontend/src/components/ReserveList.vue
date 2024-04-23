<!-- ReserveList.vue -->
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
                  <ReserveBlock v-for="item in items" :key="item.name" :item="item" @openForm="openForm"/>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import ReserveBlock from './ReserveBlock.vue';
  
  export default {
    name: 'ReserveList',
    components: {
      ReserveBlock
    },
    props: ['items', 'openForm'],
  }
  </script>
  