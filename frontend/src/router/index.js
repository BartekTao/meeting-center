import { createRouter, createWebHashHistory } from 'vue-router'
import QueryPage from '../views/QueryPage.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: QueryPage
  },
  {
    path: '/reservation',
    name: '已預約空間',
    component: () => import('../views/ReservationPage.vue')
  },
  {
    path: '/manager',
    name: '後台管理',
    component: () => import('../views/ManagerPage.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
