<template>
    <div class="period-container">
        <div class="period-name">{{ periodName }}</div>
        <div class="progress-stacked" :style="{ width: infoProgressWidth + 'px' }">
        <div v-for="(unit, index) in reservatorList" :key="index" 
            class="progress" 
            :style="{ width: unitWidth + '%' }"
            @mouseover="showText(unit)"
            @mouseleave="clearText">
            <div :class="['progress-bar', unit ? 'bg-info' : '']">
            {{ unit }}
            </div>
        </div>
        </div>
    </div>
</template>
  
  <script>
  export default {
    props: {
      periodName: String,
      reservatorList: Array,
      infoProgressWidth: Number
    },
    computed: {
      unitWidth() {
        return 100 / this.reservatorList.length;
      }
    },
    methods: {
        showText(unit) {
        this.$emit('update-show-reservator', unit);
        },
        clearText() {
        this.$emit('update-show-reservator', '');
        }
    }
  }
  </script>
  
<style>

.period-container {
    display: flex;          
    align-items: center;    
}

</style>
