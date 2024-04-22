    <template>
    <div>
        <button @click="createRoom">Create Room</button>
    </div>
    </template>

    <script>
    import gql from 'graphql-tag';
    // import client from '@/apollo-client'; 

    import { ApolloClient, createHttpLink, InMemoryCache } from '@apollo/client/core'
    import { setContext } from '@apollo/client/link/context';

    const httpLink = createHttpLink({
    uri: 'http://localhost:4000/graphql', 
    });

    const authLink = setContext((_, { headers }) => {
    return {
    headers: {
    ...headers,
    authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxMzc4NzkzM30.4TXPDGZjEb-fwEoluJIdD0KP5f_iqBlAAhbI-4Tb1-c",
    }
    }
    });

    const client = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
    });

    export default {
    name: 'CreateRoom',
    methods: {
        createRoom() {
        const mutation = gql`
            mutation myCreate($myinput: UpsertRoomInput!) {
            upsertRoom(room: $myinput) {
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
            myinput: {
            roomId: "xxx7777",
            capacity: 15,
            equipment: ["projector", "big table"],
            rules: ["no food", "no drinks"]
            }
        };

        client.mutate({
            mutation,
            variables
        }).then(result => {
            console.log('Room Created: ', result.data.upsertRoom);
        }).catch(error => {
            console.error('Error creating room: ', error);
        });
        }
        }
    }


    const TEST_QUERY = gql`
    {
        __schema {
        queryType {
            name
        }
        }
    }
    `;

    client.query({
    query: TEST_QUERY
    })
    .then(response => {
    console.log("連線成功：", response.data);
    })
    .catch(error => {
    console.error("連線失敗：", error);
    });

    </script>
