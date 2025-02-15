import { createRouter, createWebHistory } from 'vue-router';

import LogIn from '../components/LogIn.vue';
import Home from '../components/Home.vue';
import Play from '../components/PlayGame.vue';
import RunningGame from '../components/RunningGame.vue';
import Scores from '../components/GameScores.vue';
import SignUp from "../components/SignUp.vue";
import RunningGameScores from '../components/RunningGameScores.vue';



const routes = [
  {
    path: '/', // Default route
    name: 'LogIn',
    component: LogIn, // This should render your Connection
  },
  {
    path: '/home',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: "/signup",
    name: 'SignUp',
    component: SignUp
  },
  {
    path: '/runninggame',
    name: 'RunningGame',
    component: RunningGame,
    meta: { requiresAuth: true }
  },
  {
    path: '/runninggamescores',
    name: 'RunningGameScores',
    component: RunningGameScores,
    meta: { requiresAuth: true },
  },
  {
    path: '/playgame',
    name: 'PlayGame',
    component: Play,
    meta: { requiresAuth: true }
  },
  {
    path: '/gamescores',
    name: 'GameScores',
    component: Scores,
    meta: { requiresAuth: true }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;