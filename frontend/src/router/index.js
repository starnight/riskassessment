import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'login',
    component: () => import('../views/LoginView.vue')
  },
  {
    path: '/registration.html',
    name: 'registration',
    component: () => import('../views/RegistrationView.vue')
  },
  {
    path: '/assets.html',
    name: 'assets',
    component: () => import('../views/AssetsView.vue')
  },
  {
    path: '/riskassessment.html',
    name: 'riskassessment',
    component: () => import('../views/RiskAssessmentView.vue')
  },
  {
    path: '/scopes.html',
    name: 'scope',
    component: () => import('../views/ScopesView.vue')
  },
  {
    path: '/users.html',
    name: 'users',
    component: () => import('../views/UsersView.vue')
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
