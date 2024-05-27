    <template>
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
          filteredRooms: []
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
              authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxNjgwMTUxOCwibmFtZSI6Ikl2YW4gTGVlIiwic3ViIjoiNjY0NWVjZTEzNmUyYTBmMDM1OTYxYmRkIn0.-gCelsRRgt8Da11WcioKAHfe-IqxHXfD5FJkoMyHZE8",
            }
          }
        });

        this.client = new ApolloClient({
          link: authLink.concat(httpLink),
          cache: new InMemoryCache(),
        });
      },
      methods: {
        createRoom(roomInput) {

          const CREATE_ROOM_MUTATION = gql`
            mutation myCreate($my_input: UpsertRoomInput!) {
              upsertRoom(room: $my_input) {
                id
                name
                capacity
                equipments
                rules
                isDelete
              }
            }
          `;

          const variables = {
            my_input: roomInput
          };

          return this.client.mutate({
            mutation: CREATE_ROOM_MUTATION,
            variables
          }).then(response => {
            console.log("Room created or updated successfully:", response.data);
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
                    name
                    capacity
                    equipments
                    rules
                    isDelete
                  }
                }
                pageInfo {
                  startCursor
                  endCursor
                }
              }
            }
          `;

          this.client.query({
            query: GET_ALL_ROOMS_QUERY,
            fetchPolicy: 'no-cache'
          }).then(response => {
            this.rooms = response.data.paginatedRooms.edges.map(edge => edge.node);
            this.pageInfo = response.data.paginatedRooms.pageInfo;
            this.$emit('queryAllRooms', this.rooms);
          }).catch(error => {
            console.error("Failed to fetch rooms:", error);
          });
          
        },
        deleteRoom(roomId) {
          const DELETE_ROOM_MUTATION = gql`
            mutation deleteRoom($id: ID!) {
              deleteRoom(id: $id) {
                id
              }
            }
          `;
          // console.log(roomId)
          return this.client.mutate({
            mutation: DELETE_ROOM_MUTATION,
            variables: {
              id: roomId
            }
          }).then(response => {
            console.log('Room deleted:', response.data.deleteRoom.id);
          }).catch(error => {
            console.error('Error deleting room:', error);
          });
        },
        fetchAvailableRooms() {
          const GET_AVAILABLE_ROOMS = gql`
            query getAvailableRooms($startAt: Int64!, $endAt: Int64!, $rules: [Rule!], $equipments: [Equipment!], $first: Int = 20, $after: String) {
              paginatedAvailableRooms(startAt: $startAt, endAt: $endAt, rules: $rules, equipments: $equipments, first: $first, after: $after) {
                edges {
                  node {
                    id
                    name
                    capacity
                    equipments
                    rules
                    isDelete
                  }
                  cursor
                }
                pageInfo {
                  endCursor
                }
              }
            }
          `;

          const variables = {
            startAt: 1625077800,
            endAt: 1625081400,
            rules: [],
            equipments: [],
            first: 20,
            after: null
          };

          this.client.query({
            query: GET_AVAILABLE_ROOMS,
            variables
          }).then(response => {
            this.rooms = response.data.paginatedAvailableRooms.edges.map(edge => edge.node);
            this.pageInfo = response.data.paginatedAvailableRooms.pageInfo;
            this.filteredRooms = this.rooms.filter(room => !room.isDelete);
            console.log('Available rooms:', this.filteredRooms);
            this.$emit('queryAllRooms', this.rooms);
          }).catch(error => {
            console.error('Error fetching available rooms:', error);
          });
        }
      }
    }
    </script>
    