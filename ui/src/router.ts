import { h, VNodeChild } from 'vue'
import {createRouter, createWebHashHistory, RouterLink } from 'vue-router'
import type {RouteRecordRaw} from 'vue-router'
import type { MenuOption } from 'naive-ui'

import { BalanceScale, Sink, Users, AddressBook, Lightbulb, LuggageCart } from '@vicons/fa'
import { MailUnread, Construct, ListSharp, Alert } from '@vicons/ionicons5'
import { TrafficFlow } from '@vicons/carbon'
import { AlertOff20Filled } from '@vicons/fluent'

import { renderIcon } from '@/utils/render'

const prodMenuRoutes: RouteRecordRaw[] = [
  {name: 'logs', path: '/logs', component: () => import('@/views/SLogList.vue')},
  {name: 'rules', path: '/rules', component: () => import('@/views/SRules.vue')},
  {name: 'grouping', path: '/grouping', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'throttles', path: '/throttles', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'snooze', path: '/snooze', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'notifications', path: '/notifications', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'sinks', path: '/sinks', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'users', path: '/users', component: () => import('@/views/SPlaceholder.vue')},
  {name: 'roles', path: '/roles', component: () => import('@/views/SRoleList.vue')},
]

// Additional routes for dev environment (for debugging components)
const devMenuRoutes: RouteRecordRaw[] = [
  {name: 'components', path: '/components', component: () => import('@/views/SComponents.vue')},
]

function renderLink(name: string, label: string): () => VNodeChild {
  return () => h(RouterLink, {to: {name: name}}, () => label)
}

export const menuRoutes = (import.meta.env.MODE == "development") ? [...prodMenuRoutes, ...devMenuRoutes] : prodMenuRoutes

const prodMenuOptions: MenuOption[] = [
  {
    key: 'logs',
    label: renderLink("logs", "Logs"),
    icon: renderIcon(Alert),
  },

  {type: 'divider', key: 'process'},
  {
    key: 'rules',
    label: renderLink("rules", "Rules"),
    icon: renderIcon(BalanceScale),
  },
  {
    key: "grouping",
    label: renderLink("grouping", "Grouping"),
    icon: renderIcon(LuggageCart),
  },
  {
    key: "throttles",
    label: renderLink("throttles", "Throttles"),
    icon: renderIcon(TrafficFlow),
  },
  {
    key: "snooze",
    label: renderLink("snooze", "Snooze"),
    icon: renderIcon(AlertOff20Filled),
  },

  {type: 'divider', key: 'notify'},
  {
    key: 'notifications',
    label: renderLink("notifications", "Notifications"),
    icon: renderIcon(MailUnread),
  },
  {
    key: 'sinks',
    label: renderLink("sinks", "Sinks"),
    icon: renderIcon(Sink),
  },

  {type: 'divider', key: 'admin'},
  {
    key: 'users',
    label: renderLink("users", "Users"),
    icon: renderIcon(Users),
  },
  {
    key: 'roles',
    label: renderLink("roles", "Roles"),
    icon: renderIcon(AddressBook),
  },

  {type: 'divider', key: 'settings'},
]



// Additional menus for dev environment (for debugging components)
const devMenuOptions: MenuOption[] = [
  {type: 'divider', key: 'dev'},
  {
    key: 'components',
    label: renderLink("components", "Components"),
    icon: renderIcon(Lightbulb),
  },
]

export const menuOptions = (import.meta.env.MODE == "development") ? [...prodMenuOptions, ...devMenuOptions] : prodMenuOptions

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
