<template>
    <div class="period-container">
        <div class="period-name">{{ periodName }}</div>
        <div class="progress-stacked" :style="{ width: infoProgressWidth + 'px', 'margin-left': marginLeft + 'px' }">
        <div v-for="(unit, index) in reservatorList" :key="index" 
            class="progress bordered" 
            :style="{ width: unitWidth + '%' }"
            @mouseover="updateSharedData(unit, index, $event)"
            @mouseleave="hideDiv()">
            <div :class="['unselectable', 'progress-bar', unit ? 'bg-info' : '']">
            {{ unit }}
            </div>
        </div>
        </div>
    </div>
</template>
  
  <script>
  // import { inject } from 'vue';
  export default {
    props: {
      periodName: String,
      reservatorList: Array,
      infoProgressWidth: Number,
      marginLeft: Number,
    },
    computed: {
      unitWidth() {
        return 100 / this.reservatorList.length;
      }
    },
    setup(props, { emit }) {

      function updateSharedData(unit, index, event) {
        if (props.periodName === "下午：") {
          index += 6;
        }
        emit('showDiv', { unit, index, event });
      }

      function hideDiv() {
        emit('hideDiv');
      }

      return {
        updateSharedData,
        hideDiv
      };
    }
  }
  </script>
  
<style>

.period-container {
    display: flex;          
    align-items: center;    
}
.occupied {
    background-color: #E0F2FE;  
}

</style>
