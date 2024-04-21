<template>
  <js-preloader></js-preloader>
  <background-header></background-header>

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

import BackgroundHeader from '@/components/BackgroundHeader.vue';
import JsPreloader from '@/components/JsPreloader.vue';

// import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
// import { BOOKS_QUERY } from '@/components/QueryObj.vue';

export default {
  name: 'App',
  // data () {
  //   return {
  //     first_data: '',
  //   };
  // },
  components: {
    BackgroundHeader,
    JsPreloader,
  },
  // mounted() {
  //   this.fetchBooks();
  // },
  // methods: {
  //   fetchBooks() {
  //     const httpLink = createHttpLink({
  //       uri: 'http://localhost:4000/' // 替换为您的 GraphQL API URI
  //     });

  //     const client = new ApolloClient({
  //       link: httpLink,
  //       cache: new InMemoryCache()
  //     });

  //     client.query({
  //       query: BOOKS_QUERY
  //     }).then(result => {
  //       console.log(result.data);
  //     }).catch(error => {
  //       console.error("Error fetching the books:", error);
  //     });
  //   }
  // }
}


</script>