<script setup lang="ts">
import { ref, computed, withDefaults, defineEmits } from 'vue'
import {
  NForm, NFormItem,
  NGrid, NGi,
  NInput,
} from 'naive-ui'
import type { FormRules } from 'naive-ui'

import type { RuleV2, RuleV2Partial } from '@/api'
import { ValidationStatus } from '@/utils/validation'
import { SConditionInput } from '@/components'

interface Props {
  value: RuleV2,
  validation: ValidationStatus,
}

const props = withDefaults(defineProps<Props>(), {
  value: {},
  validationStatus: () => new Map(),
  feedback: () => new Map(),
})

const emit = defineEmits<{
  (e: "update:value", value: RuleV2): void,
}>()

const dataValue = computed<RuleV2>({
  get() { return props.value },
  set(v) { emit("update:value", v) },
})

const rules: FormRules = {
  name: [{required: true}],
  condition: [{required: true}],
}

</script>

<template>
  <n-form :rules="rules" :model="dataValue">
    <n-grid :span="24" :x-gap="24">
      <n-gi :span="12">
        <n-form-item
          label="Name"
          path="rule.name"
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
          path="rule.description"
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
          label="Condition"
          path="rule.condition"
          :validation-status="validation.status.get('condition')"
          :feedback="validation.feedback.get('condition')"
        >
          <s-condition-input
            v-model:value="dataValue.condition"
            size="small"
          />
        </n-form-item>
      </n-gi>
    </n-grid>
  </n-form>
</template>
