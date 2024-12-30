<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useMessage, useLoadingBar } from 'naive-ui'
import { useRouteQuery } from '@vueuse/router'

import { NAnchor, NAnchorLink, NCard, NTabs, NTabPane, NSpace, NThing, NList, NListItem, NScrollbar } from 'naive-ui'

import { getProcessConfig, type ProcessConfig } from '@/api'

const loading = useLoadingBar()
const message = useMessage()
const section = useRouteQuery("section", "overall")

const processConfig = ref<ProcessConfig>({
  transforms: [],
  groupings: [],
  profiles: [],
  silences: [],
  ratelimits: [],
  notifications: [],
  defaultDestinations: [],
})

function refresh() {
  loading.start()
  getProcessConfig()
    .then((cfg) => {
      processConfig.value = cfg
      loading.finish()
    })
    .catch((err) => {
      message.error(`error loading process config: ${err}`)
      loading.error()
    })
}

const anchors = computed(() => {
  return []
})

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-tabs v-model:value="section" type="card" placement="left">
    <n-tab-pane name="overall" tab="Overall">
      TODO
    </n-tab-pane>

    <n-tab-pane name="transforms" tab="Transforms">
      <pre>{{ processConfig.transforms }}</pre>
    </n-tab-pane>

    <n-tab-pane name="groupings" tab="Groupings">
      <pre>{{ processConfig.groupings }}</pre>
    </n-tab-pane>

    <n-tab-pane name="profiles" tab="Profiles">
      <n-space>
        <n-card
          v-for="profile in processConfig.profiles"
          :key="profile.name"
          :title="profile.name"
        >
          <n-list hoverable clickable>
            <n-list-item v-for="pattern in profile.patterns" :key="pattern.regex">
              <n-thing :title="pattern.name">
                <template v-if="pattern.description" #description>
                  {{ pattern.description }}
                </template>
                <pre>{{ pattern.regex }}</pre>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-space>
    </n-tab-pane>

    <n-tab-pane name="silences" tab="Silences">
      <pre>{{ processConfig.silences }}</pre>
    </n-tab-pane>

    <n-tab-pane name="ratelimits" tab="Ratelimits">
      <pre>{{ processConfig.ratelimits }}</pre>
    </n-tab-pane>

    <n-tab-pane name="notifications" tab="Notifications">
      <pre>{{ processConfig.notifications }}</pre>
    </n-tab-pane>
  </n-tabs>
  <n-anchor affix>
    <n-anchor-link
      v-for="anchor in anchors"
      :key="anchor"
      :title="anchor"
      :href="`#/process-config?section=${section}#${anchor}.vue`"
    />
  </n-anchor>
</template>
