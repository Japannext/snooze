<script setup lang="ts">
import { computed, withDefaults, defineEmits } from 'vue'
import {
  NForm, NFormItem,
  NGrid, NGi,
  NInput,
} from 'naive-ui'
import type { FormRules } from 'naive-ui'

import type { RoleV2, RoleV2Partial } from '@/api'
import { ValidationStatus } from '@/utils/validation'
import SSetInput from '@/components/inputs/SSetInput.vue'

interface Props {
  value: RoleV2,
  validation: ValidationStatus,
}

const props = withDefaults(defineProps<Props>(), {
  value: {},
  validationStatus: () => new Map(),
  feedback: () => new Map(),
})

const emit = defineEmits<{
  (e: "update:value", value: RoleV2): void,
}>()

const dataValue = computed<RoleV2>({
  get() { return props.value },
  set(v) { emit("update:value", v) },
})

const role: FormRules = {
  name: [{required: true}],
  condition: [{required: true}],
}

const scopeOptions = [
  {label: "rules_v2:ro", value: "rules_v2:ro"},
  {label: "rules_v2:rw", value: "rules_v2:rw"}
]

</script>

<template>
  <n-form :rules="rules" :model="dataValue">
    <n-grid :span="24" :x-gap="24">
      <n-gi :span="12">
        <n-form-item
          label="Name"
          path="role.name"
          :validation-status="validation.status.get('name')"
          :feedback="validation.feedback.get('name')"
        >
          <n-input
            v-model:value="dataValue.name"
            size="small"
            placeholder="Name"
          />
        </n-form-item>
      </n-gi>
      <n-gi :span="12">
        <n-form-item
          label="Description"
          path="role.description"
          :validation-status="validation.status.get('description')"
          :feedback="validation.feedback.get('description')"
        >
          <n-input
            v-model:value="dataValue.description"
            placeholder="Name"
            size="small"
            type="textarea"
          />
        </n-form-item>
      </n-gi>
      <n-gi :span="24">
        <n-form-item
          label="Permissions"
          path="role.scopes"
          :validation-status="validation.status.get('scopes')"
          :feedback="validation.feedback.get('scopes')"
        >
          <s-set-input
            v-model:value="dataValue.scopes"
            :options="scopeOptions"
            :size="small"
          />
        </n-form-item>
      </n-gi>
    </n-grid>
  </n-form>
</template>
