import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import { useLocalStorage } from '@vueuse/core'

import './style.css'
import { router } from '@/router'
import { messages } from '@/i18n'

import App from '@/App.vue'
// import { useLocale } from '@/stores'

const locale = useLocalStorage('locale', 'en')

const i18n = createI18n({
  legacy: false,
  locale: locale.value || 'en',
  fallbackLocale: 'en',
  messages,
})

/* if (import.meta.env.VITE_BASE_URL) {
  OpenAPI.BASE = import.meta.env.VITE_BASE_URL
} */

const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')
