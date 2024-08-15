<script setup lang="ts">

import { ref } from 'vue'

import {
  NCard, NSpace,
  NInputGroup,
} from 'naive-ui'



/*
import SKeyValue from '@/components/SKeyValue.vue'
import SMultiValue from '@/components/SMultiValue.vue'
import SField from '@/components/SField.vue'
import SOperation from '@/components/SOperation.vue'
import SCondition from '@/components/SCondition.vue'
import SConditionOperator from '@/components/SConditionOperator.vue'
import SConditionInput from '@/components/SConditionInput.vue'
import STimeConstraint from '@/components/STimeConstraint.vue'
import STimeConstraintInput from '@/components/STimeConstraintInput.vue'

import SAnyType from '@/components/inputs/SAnyType.vue'
*/

/*
const kvModel = ref<Map<string, string|number>>(new Map<string, string|number>())
const mvModel = ref<string|number|null>("")
const opModel = ref<string>("=")
const fieldModel = ref<string[]>(["resource"])
const coModel = ref<object>({kind: "=", field: ["resource"], value: {kind: "string", data: ""}})
const ciModel = ref<object>()
const tcModel = ref<object>({
  kind: "weekly",
  start: {weekday: "monday", time: "00:00:00"},
  end: {weekday: "friday", time: "23:59:59"},
})
const atModel = ref<object>({
  kind: "string",
  data: "",
})
*/

</script>

<template>
  <n-space>
    <!--
    <n-card title="s-key-value">
      <s-key-value v-model:value="kvModel" title="My Key-value" />
      <template #footer>
        <pre>{{ kvModel }}</pre>
      </template>
    </n-card>

    <n-card title="s-multi-value">
      <n-input-group>
        <s-multi-value v-model:value="mvModel" size="small" />
      </n-input-group>
      <template #footer>
        <pre>{{ JSON.stringify(mvModel) }}</pre>
      </template>
    </n-card>

    <n-card title="s-operation">
      <s-operation v-model:value="opModel" size="small" />
      <template #footer>
        <pre>{{ JSON.stringify(opModel) }}</pre>
      </template>
    </n-card>

    <n-card title="s-field">
      <n-input-group>
        <s-field v-model:value="fieldModel" size="small" />
      </n-input-group>
      <template #footer>
        <pre>{{ JSON.stringify(fieldModel) }}</pre>
      </template>
    </n-card>

    <n-card title="s-condition-operator">
      <n-input-group>
        <s-condition-operator v-model:value="coModel" size="small" />
      </n-input-group>
      <template #footer>
        <pre>{{ JSON.stringify(coModel) }}</pre>
      </template>
    </n-card>

    <n-card title="s-condition">
      <s-condition
        :value="ciModel"
      />
    </n-card>

    <n-card title="s-condition-input">
      <s-condition-input
        v-model:value="ciModel"
        size="small"
      />
      <template #footer>
        <pre>{{ JSON.stringify(ciModel, null, 2) }}</pre>
      </template>
    </n-card>
    <n-card title="s-time-constraint-input">
      <n-space vertical>
        <s-time-constraint-input
          v-model:value="tcModel"
          size="small"
        />
      </n-space>
      <template #footer>
        <pre>{{ JSON.stringify(tcModel, null, 2) }}</pre>
      </template>
    </n-card>

    <n-card title="s-time-constraint">
      <s-time-constraint :value="tcModel" />
    </n-card>

    <n-card title="s-any-type">
      <n-space item-style="display: flex;" align="center">
      <n-input-group>
        <s-any-type :value="atModel" />
      </n-input-group>
      </n-space>
      <template #footer>
        <pre>{{ JSON.stringify(atModel, null, 2) }}</pre>
      </template>
    </n-card>
    -->
  </n-space>

</template>
