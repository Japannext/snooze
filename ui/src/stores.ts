// All the local storages, wrapped or not with pinia

import { useLocalStorage } from '@vueuse/core'

import { defineRefStore } from '@/utils/pinia-store'

export const useTimeRelative = defineRefStore('time-relative', () => ({
  timeRelative: useLocalStorage('time-relative', true)
}))

export const useDevMode = defineRefStore('dev-mode', () => ({
  devMode: useLocalStorage('dev-mode', false),
}))

export const useLocale = defineRefStore('locale', () => ({
  locale: useLocalStorage('locale', 'en'),
}))

export const useTheme = defineRefStore('theme', () => ({
  theme: useLocalStorage('theme', 'light'),
}))

export const useSider = defineRefStore('sider', () => {
  return {
    collapsed: useLocalStorage('sider-collapsed', false),
  }
})

export const xSnoozeToken = useLocalStorage('x-snooze-token', "")
