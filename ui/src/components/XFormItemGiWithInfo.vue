<script setup lang="ts">
import { useSlots, withDefaults, defineProps } from "vue";
import { NFormItemGi, NButton, NIcon } from "naive-ui";

import { InformationCircleOutline } from "@/icons";

interface Props {
  label: string;
  help: string;
  size: string;
  span: number;
  required: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  label: "",
  help: "",
  size: "small",
  span: 24,
  required: false,
});

const slots = useSlots();
</script>

<template>
  <n-form-item-gi
    :size="props.size"
    :span="props.span"
    :required="props.required"
  >
    <template #label>
      {{ props.label }}
      <n-button text style="font-size: 24px">
        <n-icon :component="InformationCircleOutline" />
      </n-button>
    </template>
    <template #default>
      <slot name="default" v-bind="" />
    </template>
  </n-form-item-gi>
</template>
