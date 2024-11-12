import { h, VNodeChild } from 'vue'
import {createRouter, createWebHashHistory, RouterLink } from 'vue-router'
import type {RouteRecordRaw} from 'vue-router'
import type { MenuOption } from 'naive-ui'

import { Dashboard, Megaphone, List, AlertSnooze20Regular, Notifications, PhoneAlt, Wrench } from '@/icons'

import { renderIcon } from '@/utils/render'

export const menuRoutes: RouteRecordRaw[] = [
  {name: 'dashboard',     path: '/dashboard',     component: () => import('@/views/XDashboard.vue') },
  {name: 'alerts',        path: '/alerts',        component: () => import('@/views/XAlerts.vue')},
  {name: 'logs',          path: '/logs',          component: () => import('@/views/XLogs.vue')},
  {name: 'snooze',        path: '/snooze',        component: () => import('@/views/XSnoozes.vue')},
  {name: 'notifications', path: '/notifications', component: () => import('@/views/XNotifications.vue')},
  {name: 'escalations',   path: '/escalations',   component: () => import('@/views/XEscalations.vue')},
  {name: 'admin',         path: 'admin',          component: () => import('@/views/XAdmin.vue')}
]

export const menuOptions: MenuOption[] = [
  {key: 'dashboard',     label: renderLink("dashboard", "Dashboard"),         icon: renderIcon(Dashboard)},
  {key: 'alerts',        label: renderLink("alerts", "Alerts"),               icon: renderIcon(Megaphone)},
  {key: 'logs',          label: renderLink("logs", "Logs"),                   icon: renderIcon(List)},
  {key: 'snooze',        label: renderLink("snooze", "Snooze"),               icon: renderIcon(AlertSnooze20Regular)},
  {key: 'notifications', label: renderLink("notifications", "Notifications"), icon: renderIcon(Notifications)},
  {key: 'escalations',   label: renderLink("escalations", "Escalations"),     icon: renderIcon(PhoneAlt)},
  {key: 'admin',         label: renderLink("admin", "Admin"),                 icon: renderIcon(Wrench)},
]

function renderLink(name: string, label: string): () => VNodeChild {
  return () => h(RouterLink, {to: {name: name}}, {default: () => label})
}

const routes: RouteRecordRaw[] = [
  {
    name: 'menu',
    path: '/',
    component: () => import('@/layouts/SDefaultLayout.vue'),
    children: menuRoutes,
  },
  {
    name: 'centered',
    path: '/',
    component: () => import('@/layouts/SCenteredLayout.vue'),
    children: [
      {name: 'login', path: '/login', component: () => import('@/views/SLogin.vue')},
    ],
  },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
})
