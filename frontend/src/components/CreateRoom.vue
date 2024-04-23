<template>
    <div>
        <button @click="fetchSchema">Fetch GraphQL Schema</button>
    </div>
    </template>
    
    <script>
    import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
    import { setContext } from '@apollo/client/link/context';
    import gql from 'graphql-tag';
    
    export default {
      name: 'GraphQLTester',
      methods: {
        fetchSchema() {
          const httpLink = createHttpLink({
            uri: 'http://localhost:8080/query', 
          });
    
          const authLink = setContext((_, { headers }) => {
            return {
              headers: {
                ...headers,
                authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxMzg4MTQzOH0.sm2X27FYWEE4l6ph2lssQSpywtJQBfbRz7iS-6g_Xhw",
              }
            }
          });
    
          const client = new ApolloClient({
            link: authLink.concat(httpLink),
            cache: new InMemoryCache(),
          });
    
          const GET_SCHEMA_QUERY = gql`
            {
              __schema {
                queryType {
                  name
                }
              }
            }
          `;
    
          client.query({
            query: GET_SCHEMA_QUERY
          }).then(response => {
            console.log("Schema fetched successfully:", response.data);
          }).catch(error => {
            console.error("Failed to fetch schema:", error);
          });
          
        }
      }
    }
    </script>
    