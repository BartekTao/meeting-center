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
          users: [],
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
              authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxNjk5NTIwNiwibmFtZSI6Ikl2YW4gTGVlIiwic3ViIjoiNjY0NWVjZTEzNmUyYTBmMDM1OTYxYmRkIn0.ZV9ZaG8CCS9u7o2V920ch0B2vrz4VxQ68Bk2qJq0rT8",
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
        createEvent(eventInput) {

          const CREATE_EVENT_MUTATION = gql`
            mutation createEvent($input: UpsertEventInput!) {
              upsertEvent(input: $input) {
                id
                title
                description
                startAt
                endAt
                roomReservation {
                  room {
                    id
                    name
                    capacity
                    equipments
                    rules
                    isDelete
                  }
                  status
                }
                participants {
                  id
                  name
                  email
                }
                notes
                remindAt
                creator {
                  id
                  name
                  email
                }
                isDelete
              }
            }
          `;


          const variables = {
            input: eventInput
          };

          return this.client.mutate({
            mutation: CREATE_EVENT_MUTATION,
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
          
          this.calculateStartOfDay(variables.startAt);
          this.calculateEndOfDay(variables.endAt);
          
          this.client.query({
            query: GET_AVAILABLE_ROOMS,
            variables
          }).then(response => {
            this.rooms = response.data.paginatedAvailableRooms.edges.map(edge => edge.node);
            this.ids = this.rooms.map(room => room.id);
            
          this.queryRoomSchedules(this.ids, this.startOfDayTimestamp, this.endOfDayTimestamp)
            .then(response => {
              const edges = response.data.paginatedRoomSchedules.edges;
              this.rooms = edges.map(edge => {
                return {
                  ...edge.node.room,
                  schedules: edge.node.schedules
                };
              });
              this.$emit('fetchAvailableRooms', this.rooms);
              this.pageInfo = response.data.paginatedRoomSchedules.pageInfo;
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
            return response
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
        },
        queryUsers() {

          const GET_PAGINATED_USERS = gql`
            query getPaginatedUsers($first: Int = 20, $after: String) {
              paginatedUsers(first: $first, after: $after) {
                edges {
                  node {
                    id
                    name
                  }
                }
              }
            }
          `;

          const variables = {
            first: 20,
            after: null
          };

          this.client.query({
            query: GET_PAGINATED_USERS,
            variables
          }).then(response => {
            this.users = response.data.paginatedUsers.edges.map(edge => edge.node);
            this.$emit('queryUsers', this.users);
          }).catch(error => {
            console.error('Error fetching users:', error);
          });
        }
      },
    }
    </script>
    