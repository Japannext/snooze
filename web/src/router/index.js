import Vue from 'vue'
import Router from 'vue-router'

// Containers
const TheContainer = () => import('@/containers/TheContainer')
// Components
const Record = () => import('@/views/Record')
const Aggregate = () => import('@/views/Aggregate')
const Status = () => import('@/views/Status')
const Rule = () => import('@/views/Rule')
const AggregateRule = () => import('@/views/AggregateRule')
const Snooze = () => import('@/views/Snooze')
const Notification = () => import('@/views/Notification')
const Action = () => import('@/views/Action')
const Widget = () => import('@/views/Widget')
const Patlite = () => import('@/views/Patlite')
const User = () => import('@/views/User')
const Role = () => import('@/views/Role')
const Settings = () => import('@/views/Settings')
const Profile = () => import('@/views/Profile')
const Login = () => import('@/views/Login')

Vue.use(Router)

export default new Router({
  mode: 'hash', // https://router.vuejs.org/api/#mode
  linkActiveClass: 'active',
  scrollBehavior: () => ({ y: 0 }),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/',
      redirect: '/record',
      name: 'Home',
      component: TheContainer,
      children: [
        {
          path: 'record',
          name: 'Records',
          component: Record
        },
        {
          path: 'aggregate',
          name: 'Aggregates',
          component: Aggregate
        },
        {
          path: 'status',
          name: 'Status',
          component: Status
        },
        {
          path: 'rule',
          name: 'Rules',
          component: Rule
        },
        {
          path: 'aggregaterule',
          name: 'Aggregate Rules',
          component: AggregateRule
        },
        {
          path: 'snooze',
          meta: {label: 'Snooze'},
          component: Snooze,
        },
        {
          path: 'snooze/*',
          meta: {label: 'Snooze'},
          component: Snooze,
        },
        {
          path: '/notification',
          name: 'Notifications',
          component: Notification,
        },
        {
          path: '/action',
          name: 'Actions',
          component: Action,
        },
        {
          path: '/widget',
          name: 'Widgets',
          component: Widget,
        },
        {
          path: '/patlite',
          name: 'Patlite',
          component: Patlite,
        },
        {
          path: '/user',
          name: 'Users',
          component: User,
        },
        {
          path: '/role',
          name: 'Roles',
          component: Role,
        },
        {
          path: '/settings',
          name: 'Settings',
          component: Settings,
        },
        {
          path: '/profile',
          name: 'Profile',
          component: Profile,
        },
      ]
    },
  ]
})
