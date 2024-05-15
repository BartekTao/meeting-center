    <template>
      <!-- <div>
          <button @click="createRoom(test_room)">Fetch GraphQL Schema</button>
      </div> -->
      <div></div>
    </template>
    
    <script>
    import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
    import { setContext } from '@apollo/client/link/context';
    import gql from 'graphql-tag';
    
    export default {
      name: 'GraphQLTester',
      data() {
        return {
          rooms: [],
          pageInfo: {},
        };
      },
      created() {
        const httpLink = createHttpLink({
          uri: 'http://localhost:8080/query', 
        });

        const authLink = setContext((_, { headers }) => {
          return {
            headers: {
              ...headers,
              authorization: "Bearer ",
            }
          }
        });

        this.client = new ApolloClient({
          link: authLink.concat(httpLink),
          cache: new InMemoryCache(),
        });
      },
      methods: {
        generateRandomId() {
          const randomId = Math.floor(Math.random() * 900000000) + 100000000;
          return randomId.toString();
        },
        createRoom(roomInput) {

          const CREATE_ROOM_MUTATION = gql`
            mutation myCreate($my_input: UpsertRoomInput!) {
              upsertRoom(room: $my_input) {
                id
                roomId
                capacity
                equipment
                rules
                isDelete
              }
            }
          `;

          const variables = {
            my_input: roomInput
          };

          this.client.mutate({
            mutation: CREATE_ROOM_MUTATION,
            variables
          }).then(response => {
            console.log("Room created or updated successfully:", response.data);
            // this.$emit('createRoom', response.data);
          }).catch(error => {
            console.error("Error creating or updating room:", error);
          });
        },
        queryAllRooms() {

          const GET_ALL_ROOMS_QUERY = gql`
            query {
              paginatedRooms(first: 3, after: "") {
                edges {
                  node {
                    id
                    roomId
                    capacity
                    equipment
                    rules
                    isDelete
                  }
                }
                pageInfo {
                  hasNextPage
                  hasPreviousPage
                  startCursor
                  endCursor
                }
              }
            }
          `;

          this.client.query({
            query: GET_ALL_ROOMS_QUERY
          }).then(response => {

            this.rooms = response.data.paginatedRooms.edges.map(edge => edge.node);
            this.pageInfo = response.data.paginatedRooms.pageInfo;
            let new_rooms = this.rooms.map(room => {
              return {
                ...room,
                id: room.id === '' ? this.generateRandomId() : room.id
              };
            });

            // console.log("Rooms fetched successfully:", new_rooms);
            this.$emit('queryAllRooms', new_rooms);
            console.log(new_rooms);
          }).catch(error => {
            console.error("Failed to fetch rooms:", error);
          });
          
        },
        deleteRoom(roomId) {
          const DELETE_ROOM_MUTATION = gql`
            mutation DeleteRoom($id: ID!) {
              deleteRoom(id: $id) {
                id
              }
            }
          `;

          this.client.mutate({
            mutation: DELETE_ROOM_MUTATION,
            variables: {
              id: roomId
            }
          }).then(response => {
            console.log('Room deleted:', response.data.deleteRoom.id);
          }).catch(error => {
            console.error('Error deleting room:', error);
          });
        }
      }
    }
    </script>
    