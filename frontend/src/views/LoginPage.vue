<template>
  <div class="page">
      <div class="container">
          <h2 class="form-title">Meeting Center</h2>
          <form @submit.prevent="handleSubmit">
              <div class="mb-3">
                  <label for="exampleInputEmail1" class="form-label">User name</label>
                  <input type="text" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp">
              </div>
              <div class="mb-3">
                  <label for="exampleInputPassword1" class="form-label">Password</label>
                  <input type="password" class="form-control" id="exampleInputPassword1">
              </div>
              <div class="button-container">
                  <button type="submit" class="btn btn-primary">Submit</button>
              </div>
          </form>
      </div>
  </div>
</template>

<script>
export default {
  methods: {
      handleSubmit() {
          console.log('Form submitted');
          this.$router.push('/');
      }
  }
}
</script>

<style scoped>
.page {
  display: grid;
  place-items: center;
  height: 100vh;
  background: linear-gradient(to bottom, #f8f9fa 0%, #00c6ff 70%, #0072ff 100%);
}

.container {
  max-width: 500px;
  width: 100%;
  padding: 2rem;
  border: 1px solid #ced4da;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background-color: #f8f9fa;
}

.form-title {
  text-align: center;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
  color: #333;
}

.form-label {
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.form-control {
  border-radius: 5px;
  margin-bottom: 1rem;
  padding: 0.75rem;
}

.form-text {
  color: #6c757d;
  margin-bottom: 1rem;
}

.form-check-label {
  margin-left: 0.5rem;
}

.button-container {
  display: flex;
  justify-content: center;
}

.btn {
  padding: 0.5rem 1.5rem;
  border-radius: 5px;
}
</style>