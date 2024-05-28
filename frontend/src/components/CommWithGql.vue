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
          filteredRooms: [],
          startOfDayTimestamp: null,
          endOfDayTimestamp: null,
          ids: [],
          edges: [],
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
              authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxNjkwODY4OCwibmFtZSI6Ikl2YW4gTGVlIiwic3ViIjoiNjY0NWVjZTEzNmUyYTBmMDM1OTYxYmRkIn0.fQ_yHLfcu3U4RfN6HkS5zBeHg1G8JnSfZA_ajmsv7NM",
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
        fetchAvailableRooms(variables) {
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
          
          // console.log(variables.startAt);
          this.calculateStartOfDay(variables.startAt);
          // console.log(this.startOfDayTimestamp);
          this.calculateEndOfDay(variables.endAt);
          console.log(variables);
          this.client.query({
            query: GET_AVAILABLE_ROOMS,
            variables
          }).then(response => {
            this.rooms = response.data.paginatedAvailableRooms.edges.map(edge => edge.node);
            this.$emit('fetchAvailableRooms', this.rooms);
            this.ids = this.rooms.map(room => room.id);
            // const schedules = this.queryRoomSchedules(this.ids, this.startOfDayTimestamp, this.endOfDayTimestamp)
            // console.log(this.rooms);
            
          this.queryRoomSchedules(this.ids, this.startOfDayTimestamp, this.endOfDayTimestamp)
            .then(response => {
              console.log(response)
              const edges = response.data.paginatedRoomSchedules.edges;
              this.rooms = edges.map(edge => {
                return {
                  ...edge.node.room,
                  schedules: edge.node.schedules
                };
              });
              this.pageInfo = response.data.paginatedRoomSchedules.pageInfo;
              // console.log("Room schedules set in state:", response.data.paginatedRoomSchedules.edges[0].node.schedules);
              // console.log("Room schedules set in state:", this.rooms);
            })
            .catch(error => {
              this.error = error;
              console.error("Error in queryRoomSchedules:", error);
            });

          }).catch(error => {
            console.error('Error fetching available rooms:', error);
          });
        },   
        queryRoomSchedules(ids, startAt, endAt) {

          const GET_ROOM_SCHEDULES = gql`
            query getRoomSchedules(
              $ids: [String!]!,
              $startAt: Int64!,
              $endAt: Int64!,
              $rules: [Rule!],
              $equipments: [Equipment!],
              $first: Int = 20,
              $after: String
            ) {
              paginatedRoomSchedules(
                ids: $ids,
                startAt: $startAt,
                endAt: $endAt,
                rules: $rules,
                equipments: $equipments,
                first: $first,
                after: $after
              ) {
                edges {
                  node {
                    room {
                      id
                      name
                      capacity
                      equipments
                      rules
                      isDelete
                    }
                    schedules {
                      startAt
                      endAt
                    }
                  }
                  cursor
                }
                pageInfo {
                  endCursor
                }
              }
            }
          `;

          const defaultVariables = {
            ids: ids, //["6655178de1dfe965fa4b1951"],
            startAt: startAt, // 1625077800,
            endAt:endAt, // 1625081400,
            rules: [],
            equipments: [],
            first: 20,
            after: null
          };

          return this.client.query({
            query: GET_ROOM_SCHEDULES,
            variables: defaultVariables,
          }).then(response => {
            // console.log(response); 
            // console.log(response.data.paginatedRoomSchedules.edges[0].node.schedules); 
            return response // .data.paginatedRoomSchedules.edges
          }).catch(error => {
            this.error = error;
            console.error("Failed to fetch room schedules:", error);
          });
        },
        calculateStartOfDay(timestamp) {
          const date = new Date(timestamp);
          const year = date.getFullYear();
          const month = date.getMonth();
          const day = date.getDate();

          const startOfDay = new Date(year, month, day, 0, 0, 0, 0);

          this.startOfDayTimestamp = startOfDay.getTime();
        },
        calculateEndOfDay(timestamp) {
          const date = new Date(timestamp);
          const year = date.getFullYear();
          const month = date.getMonth();
          const day = date.getDate();

          const endOfDay = new Date(year, month, day, 23, 59, 59, 999);

          this.endOfDayTimestamp = endOfDay.getTime();
        }
      },
    }
    </script>
    