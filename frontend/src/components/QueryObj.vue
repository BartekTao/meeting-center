<script>
import gql from 'graphql-tag';

export const BOOKS_QUERY = gql`
  query Query {
    books {
      author
      id
      publish
      name
    }
  }
`;

export default {
  name: 'QueryObj',
};
</script>
