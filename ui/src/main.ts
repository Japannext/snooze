import axios from 'axios'
import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import { createPinia } from 'pinia'
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
const pinia = createPinia()

if (import.meta.env.VITE_BASE_URL) {
  axios.defaults.baseURL = import.meta.env.VITE_BASE_URL
}
axios.defaults.headers.post['Content-Type'] = 'application/json'

const app = createApp(App)
app.use(router)
app.use(i18n)
app.use(pinia)
app.mount('#app')
