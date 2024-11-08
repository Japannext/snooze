<script setup lang="ts">

// Functions / Libraries
import { ref, h, computed, onMounted, defineProps, defineModel } from 'vue'
import { useMessage, useLoadingBar } from 'naive-ui'
import axios from 'axios'

// Components
import { NModal, NDatePicker, NSpace, NButton, NCard, NTabs, NTabPane, NForm, NFormItem, NFormItemGi, NSelect, NInput, NGrid } from 'naive-ui'
import SGroupTag from '@/components/SGroupTag.vue'

// Types
import type { AxiosResponse } from 'axios'
import type { Snooze, Group } from '@/api/types'
import type { VNodeChild } from 'vue'
import type { SelectOption } from 'naive-ui'

// Variables
const show = defineModel('show', {type: Boolean, default: false})
const loading = useLoadingBar()
const selectLoading = ref<boolean>(false)
const message = useMessage()
const item = ref<Snooze>({})
const selectedGroups = ref<Group[]>([])
const groupOptions = ref<SelectOption>([])

// Props
defineProps<{
  show: boolean;
}>()


onMounted(() => {
  getGroupOptions("")
})

async function searchGroup(query: string) {
  await getGroupOptions(query)
}

function getGroupOptions(query: string): Promise {
  var params = {}
  if (query != "") {
    params.search = `*${query}*`
  }
  console.log(`getGroupOptions(${JSON.stringify(params)})`)
  return axios.get<Group>("/api/groups", {params: params})
    .then((resp: AxiosResponse) => {
      if (resp.data) {
        groupOptions.value = resp.data.items.map(group => {
          var labels = Object.keys(group.labels).map(key => {
            return `${key}=${group.labels[key]}`
          }).join(", ")
          return {
            value: group,
            label: `[${group.name}] ${labels}`,
          }
        })
      } else {
        console.log("Groups not found")
      }
      selectLoading.value = false
    })
    .catch(err => {
      message.error(`failed to load logs: ${err}`)
      selectLoading.value = false
    })
}

function create() {
  loading.start()
  axios.post("/api/snooze", item.value)
    .then(resp => {
      show.value = false
      loading.finish()
    })
    .catch(err => {
      message.error(`failed to snooze: ${err}`)
      loading.error()
    })
}

function cancel() {
  show.value = false
  item.value = {}
}

</script>

<template>
  <n-modal
    v-model:show="show"
    style="width: 800px"
  >
    <n-card title="Add snooze">
      <n-form size="small">
        <n-grid>
          <n-form-item-gi label="By group" size="small" :span="24">
            <n-select
              v-model:value="item.groups"
              multiple
              remote
              clearable
              filterable
              :loading="selectLoading"
              :options="groupOptions"
              @search="searchGroup"
            />
          </n-form-item-gi>
          <n-form-item-gi label="From" size="small" :span="12">
            <n-date-picker v-model:value="item.startAt" />
          </n-form-item-gi>
          <n-form-item-gi label="To" size="small" :span="12">
            <n-date-picker v-model:value="item.expireAt" />
          </n-form-item-gi>
          <n-form-item-gi label="Reason" size="small" :span="24">
            <n-input v-model:value="item.reason" type="textarea" />
          </n-form-item-gi>
        </n-grid>
      </n-form>
      <template #footer>
        <pre>{{ item }}</pre>
      </template>
      <template #action>
        <n-space justify="end">
          <n-button @click="cancel">Cancel</n-button>
          <n-button type="warning" @click="create">Snooze</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>
