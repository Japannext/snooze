// All the local storages, wrapped or not with pinia

import { ref, reactive } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { getTags, type Tag } from '@/api'
import { DateTime } from 'luxon'

export function useTimeRelative() {
  return useLocalStorage('time-relative', false)
}

export function useDevMode() {
  return useLocalStorage('dev-mode', false)
}

export function useLocale() {
  return useLocalStorage('locale', 'en')
}

export function useTheme() {
  return useLocalStorage('theme', 'light')
}

export var tagMap = new Map<string, Tag>()
setInterval(() => {
  var params = {
    pagination: {},
  }
  getTags(params)
    .then((list) => {
      var newMap = new Map<string, Tag>()
      list.items.forEach(tag => {
        newMap.set(tag.name, tag)
      })
      tagMap = newMap
    })
}, 120000)
