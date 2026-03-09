import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/machines' },
    { path: '/machines', name: 'machines', component: () => import('@/views/MachinesView.vue') },
    { path: '/locations', name: 'locations', component: () => import('@/views/LocationsView.vue') },
    { path: '/map', name: 'map', component: () => import('@/views/MapView.vue') },
    { path: '/statistics', name: 'statistics', component: () => import('@/views/StatisticsView.vue') },
  ],
})

export default router
