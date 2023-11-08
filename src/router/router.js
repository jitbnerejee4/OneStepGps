import { createRouter, createWebHistory } from 'vue-router';

import HomeComponent from '@/components/HomeComponent.vue'
import LandingComponent from '@/components/LandingComponent.vue'
import ForgetPasswordComponent from '@/components/ForgetPasswordComponent'
import ResetPasswordComponent from '@/components/ResetPasswordComponent'

const routes = [
    { path: '/home', name: 'home', component: HomeComponent },
    { path: '/', name: 'landing', component: LandingComponent },
    {path: '/forgot-password', name: 'forgot-password', component: ForgetPasswordComponent },
    {path: '/reset-password', name: 'reset-password', component: ResetPasswordComponent}
    
  ];

const router = createRouter({
    history: createWebHistory(),
    routes,
});
router.beforeEach((to, from, next) => {
  const publicPages = ['/', '/forgot-password', '/reset-password'];
  const authRequired = !publicPages.includes(to.path);
  const loggedIn = localStorage.getItem('token');

  if (authRequired && !loggedIn) {
      return next('/');
  }

  next();
});
  
export default router;