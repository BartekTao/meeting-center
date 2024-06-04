<template>
  <!-- <background-header></background-header> -->

  <router-view></router-view>
</template>

<style lang="scss">
@import "bootstrap";
</style>

<style scoped>
@import '@/assets/css/fontawesome.css';
@import '@/assets/css/templatemo-plot-listing.css';
@import '@/assets/css/animated.css';
@import '@/assets/css/customization.css';
@import '@/assets/css/bootstrap.min.css';
</style>

<script>

// import BackgroundHeader from '@/components/BackgroundHeader.vue';

export default {
  name: 'App',
  data () {
    return {
      first_data: '',
    };
  },
  mounted() {
    console.log('Mounted hook is triggered');
    // Check if there is a token in the URL (after login redirect)
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');
    
    // If token exists, store it in localStorage
    if (token) {
      localStorage.setItem('token', token);
      
      // Remove token from URL to prevent token leakage
      window.history.replaceState({}, document.title, "/");
    }
  }
  // components: {
  //   BackgroundHeader,
  // },
}


</script>