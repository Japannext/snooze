import { h, VNodeChild } from 'vue'
import {createRouter, createWebHashHistory, RouterLink } from 'vue-router'
import type {RouteRecordRaw} from 'vue-router'
import type { MenuOption } from 'naive-ui'

import { BalanceScale, Sink, Users, AddressBook, Lightbulb, LuggageCart, Wrench, PhoneAlt } from '@vicons/fa'
import { MailUnread, Construct, ListSharp, Alert, List, Megaphone, Notifications } from '@vicons/ionicons5'
import { TrafficFlow, Dashboard } from '@vicons/carbon'
import { AlertOff20Filled, AlertSnooze20Regular } from '@vicons/fluent'

import { renderIcon } from '@/utils/render'

export const menuRoutes: RouteRecordRaw[] = [
  {name: 'dashboard', path: '/dashboard', component: () => import('@/views/SDashboard.vue') },
  {name: 'alerts', path: '/alerts', component: () => import('@/views/SAlertList.vue')},
  {name: 'logs', path: '/logs', component: () => import('@/views/SLogs.vue')},
  {name: 'snooze', path: '/snooze', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'notifications', path: '/notifications', component: () => import('@/views/SNotifications.vue')},
  {name: 'escalations', path: '/escalations', component: () => import('@/views/SEscalations.vue')},
  {name: 'admin', path: 'admin', component: () => import('@/views/SAdmin.vue')}
]

export const menuOptions: MenuOption[] = [
  {key: 'dashboard', label: renderLink("dashboard", "Dashboard"), icon: renderIcon(Dashboard)},
  {key: 'alerts', label: renderLink("alerts", "Alerts"), icon: renderIcon(Megaphone)},
  {key: 'logs', label: renderLink("logs", "Logs"), icon: renderIcon(List)},
  {key: 'snooze', label: renderLink("snooze", "Snooze"), icon: renderIcon(AlertSnooze20Regular)},
  {key: 'notifications', label: renderLink("notifications", "Notifications"), icon: renderIcon(Notifications)},
  {key: 'escalations', label: renderLink("escalations", "Escalations"), icon: renderIcon(PhoneAlt)},
  {key: 'admin', label: renderLink("admin", "Admin"), icon: renderIcon(Wrench)},
]

function renderLink(name: string, label: string): () => VNodeChild {
  return () => h(RouterLink, {to: {name: name}}, () => label)
}

const routes: RouteRecordRaw[] = [
  {
    name: 'menu',
    path: '/',
    redirect: "logs",
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
