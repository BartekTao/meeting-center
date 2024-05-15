import { createRouter, createWebHashHistory } from 'vue-router'
import ReservePage from '../views/ReservePage.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: ReservePage
  },
  {
    path: '/reserved-page',
    name: '已預約空間',
    component: () => import('../views/ReservedPage.vue')
  },
  {
    path: '/manager',
    name: '後台管理',
    component: () => import('../views/RoomPage.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
