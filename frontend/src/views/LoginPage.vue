<template>
    <div class="login-page">
      <h1>登入</h1>
      <button @click="loginWithGoogle">Google 登入</button>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    name: 'LoginPage',
    methods: {
    async loginWithGoogle() {
            try {
                window.location.href = 'http://localhost:8080/auth/google/login';
            } catch (error) {
                console.error('登入失敗:', error);
            }
        }
    },
  }
  </script>
  
  <style scoped>
  .login-page {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }
  button {
    padding: 10px 20px;
    font-size: 16px;
  }
  </style>