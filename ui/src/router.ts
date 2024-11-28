import { h, VNodeChild } from 'vue'
import {createRouter, createWebHashHistory, RouterLink } from 'vue-router'
import type {RouteRecordRaw} from 'vue-router'
import type { MenuOption } from 'naive-ui'

import { Dashboard, Megaphone, List, AlertSnooze20Regular, Notifications, TaskSettings, PhoneAlt, Pricetag, Wrench } from '@/icons'

import { renderIcon } from '@/utils/render'

export const menuRoutes: RouteRecordRaw[] = [
  {name: 'dashboard',     path: '/dashboard',     component: () => import('@/views/XDashboard.vue') },
  {name: 'alerts',        path: '/alerts',        component: () => import('@/views/XAlerts.vue')},
  {name: 'logs',          path: '/logs',          component: () => import('@/views/XLogs.vue')},
  {name: 'snooze',        path: '/snooze',        component: () => import('@/views/XSnoozes.vue')},
  {name: 'notifications', path: '/notifications', component: () => import('@/views/XNotifications.vue')},
  {name: 'escalations',   path: '/escalations',   component: () => import('@/views/XEscalations.vue')},
  {name: 'tags',          path: '/tags',          component: () => import('@/views/XTags.vue')},
  {name: 'process-config', path: '/process-config', component: () => import('@/views/XProcessConfig.vue')},
]

export const menuOptions: MenuOption[] = [
  {key: 'dashboard',     label: renderLink("dashboard", "Dashboard"),         icon: renderIcon(Dashboard)},
  {key: 'alerts',        label: renderLink("alerts", "Alerts"),               icon: renderIcon(Megaphone)},
  {key: 'logs',          label: renderLink("logs", "Logs"),                   icon: renderIcon(List)},
  {key: 'snooze',        label: renderLink("snooze", "Snooze"),               icon: renderIcon(AlertSnooze20Regular)},
  {key: 'notifications', label: renderLink("notifications", "Notifications"), icon: renderIcon(Notifications)},
  {key: 'escalations',   label: renderLink("escalations", "Escalations"),     icon: renderIcon(PhoneAlt)},
  {key: "tags",          label: renderLink("tags", "Tags"),                   icon: renderIcon(Pricetag)},
  {key: "process-config", label: renderLink("process-config", "Process Config"), icon: renderIcon(TaskSettings)},
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
      {name: 'login', path: '/login', component: () => import('@/views/auth/XLogin.vue')},
      {name: 'oidc-callback', path: '/oidc/callback', component: () => import('@/views/auth/XOidcCallback.vue') },
      {name: 'oidc-silent-callback', path: '/oidc/silent-callback', component: () => import('@/views/auth/XOidcSilentCallback.vue') },
    ],
  },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
})
