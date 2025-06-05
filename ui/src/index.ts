import axios from "axios";
import { createApp } from "vue";
import { createI18n } from "vue-i18n";
import { useLocalStorage } from "@vueuse/core";
import PrimeVue from "primevue/config";
import ToastService from "primevue/toastservice"
import DialogService from "primevue/dialogservice"
import Aura from "@primeuix/themes/aura"

import { router } from "@/router";
//import { messages } from '@/i18n'
import "@/style.scss";
import 'primeicons/primeicons.css'

import App from "@/App.vue";

/*
import { useLocale } from '@/stores'
const locale = useLocalStorage('locale', 'en')
const i18n = createI18n({
  legacy: false,
  locale: locale.value || 'en',
  fallbackLocale: 'en',
  messages,
})
*/

if (import.meta.env.BASE_URL) {
  axios.defaults.baseURL = import.meta.env.BASE_URL;
}
axios.defaults.headers.post["Content-Type"] = "application/json";
axios.defaults.withCredentials = true;

const app = createApp(App);
app.use(router);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
  },
})
app.use(ToastService)
app.use(DialogService)
// app.use(i18n)
app.mount("#app");
