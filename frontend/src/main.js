import { createApp } from 'vue'
// import { createApp, h } from 'vue'
import App from './App.vue'
import router from './router'
// import { createApolloProvider } from '@vue/apollo-option'
// import { ApolloClient, createHttpLink, InMemoryCache } from '@apollo/client/core'

createApp(App).use(router).mount('#app')

// const httpLink = createHttpLink({
//   uri: 'http://localhost:3020/graphql',
// })

// const cache = new InMemoryCache()
// const apolloClient = new ApolloClient({
//   link: httpLink,
//   cache,
// })
// console.log(apolloClient)

// const apolloProvider = createApolloProvider({
//   defaultClient: apolloClient,
// })

// const app = createApp({
//   render: () => h(App),
// })

// app.use(apolloProvider)

