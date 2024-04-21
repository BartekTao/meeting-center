<template>
    <div class="page-heading">
      <div class="container">
        <div class="row">
          <div class="col-lg-8">
            <div class="top-text header-text">
              <h2>{{ pageContent }}</h2>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>
  
  <script>
  export default {
    name: 'HeadPage',
    props: {
      pageContent: {
        type: String,
        required: true,
        default: 'Default Heading'
      }
    }
  }
  </script>
  
  