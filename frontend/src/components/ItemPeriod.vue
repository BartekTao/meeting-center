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
            {{ unit.nickName }}
            </div>
        </div>
        </div>
    </div>
</template>
  
 
  <script>
  // import { inject } from 'vue';
  export default {
    emits: ['showDiv', 'hideDiv'],
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
    methods: {
      updateSharedData(unit, index, event) {
        if (this.periodName === "下午：") {
          index += 6;
        }
        this.$emit('showDiv', { unit, index, event });
      },
      hideDiv() {
        this.$emit('hideDiv');
      }
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
