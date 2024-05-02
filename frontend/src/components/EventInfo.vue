<template>
    <div class="card text-bg-info mb-3 white-text" :style="showDivStyle">
      <div class="card-header">詳細資訊</div>
      <div class="card-body">
        <h5 class="card-title">預約人：{{ reservator }}</h5>
        <h5 class="card-title">開始時間：{{ start_time }}</h5>
        <h5 class="card-title">結束時間：{{ end_time }}</h5>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        showDivStyle: {
          display: 'none',
          position: 'absolute',
          maxWidth: '18rem',
          zIndex: 1000,
          left: '0px',
          top: '0px'
        },
        time_period: [],  // 这个数组需要从父组件传入或其他方式获取
        reservator: '',
        start_time: '',
        end_time: ''
      };
    },
    methods: {
      showDiv(data) {
        this.showDivStyle.display = 'block';
        this.showDivStyle.left = data.event.pageX + 'px';
        this.showDivStyle.top = data.event.pageY + 'px';
        this.reservator = data.unit;
        this.start_time = this.time_period[data.index];
        this.end_time = this.time_period[data.index + 1];
      },
      hideDiv() {
        this.showDivStyle.display = 'none';
      }
    }
  }
  </script>
  
  <style scoped>
  .text-bg-info {
    color: white;
  }
  </style>
  