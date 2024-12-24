// All the local storages, wrapped or not with pinia

import { ref, reactive } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { getTags, type Tag } from '@/api'

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

const TAG_EXPIRY = 60 * 1000 // 1 minute in milliseconds

export function useTagMap() {
  const tagMap = ref<Map<string, Tag>>(new Map<string, Tag>())
  const expiry = ref<number>(0)
  return reactive({
    tagMap: tagMap,
    expiry: expiry,
    get(key: string): Tag|undefined {
      var now = new Date()
      if (now.getMilliseconds() >= expiry.value) {
        updateTagMap()
      }
      return tagMap.value.get(key)
    },
    mget(keys: string[]): Array<Tag|undefined> {
      var now = new Date()
      if (now.getMilliseconds() >= expiry.value) {
        updateTagMap()
      }
      return keys.map(key => tagMap.value.get(key))
    },
    updateTagMap: async() => {
      var params = {}
      var list = await getTags(params)
      var newMap = new Map<string, Tag>()
      list.items.forEach(tag => {
        newMap.set(tag.name, tag)
      })
      tagMap.value = newMap
      var now = new Date()
      expiry.value = now.getMilliseconds() + TAG_EXPIRY
    },
  })
}
