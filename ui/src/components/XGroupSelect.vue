<script setup lang="ts">
import { ref, defineModel, computed, onMounted } from "vue";
import {
  getPipelineGroupings,
  getAutocompleteLog,
  type Grouping,
  type Group,
  type GetAutocompleteLogParams,
} from "@/api";

const loading = ref(false);

const props = defineProps<{}>();

const group = defineModel("group", {
  type: Object,
  default: { name: "", hash: "", labels: {} },
});
const grouping = ref<Grouping | undefined>(undefined);
const groupingOptions = ref<Grouping[]>();

/*
const emit = defineEmits<{
  (e: "update:group", value: Group): void;
}>();
*/

function updateGrouping(newValue: Grouping) {
  grouping.value = newValue;
  group.value = newValue.name;
}

function refresh() {
  console.log(`XGroupSelect.refresh()`);
  loading.value = true;
  getPipelineGroupings()
    .then((list) => {
      groupingOptions.value = list;
      console.log(`Groupings: ${JSON.stringify(groupingOptions.value)}`);
    })
    .finally(() => {
      loading.value = false;
    });
}

onMounted(() => {
  refresh();
});

function autocomplete(fieldName: string) {
  var value = group.value.labels[fieldName];
  if (groupBy) var params = {};
  getAutocompleteLog(params).then((suggestions) => {});
}
</script>

<template>
  <InputGroup>
    <IftaLabel>
      <Select
        v-if="!loading"
        inputId="grouping"
        v-model="grouping"
        :options="groupingOptions"
        :loading="loading"
        optionLabel="name"
        size="small"
        style="width: 100"
        @update:modelValue="updateGrouping"
      />
      <label for="grouping">Grouping</label>
    </IftaLabel>
    <template v-if="grouping">
      <InputGroupAddon size="small">#</InputGroupAddon>
      <template v-for="groupBy in grouping.groupBy">
        <IftaLabel>
          <AutoComplete
            v-if="groupBy.kind == 'label'"
            :inputId="`${groupBy.kind}-${groupBy.key}`"
            v-model="selector.matchLabels[groupBy.key]"
            size="small"
            @complete="autocomplete(groupBy)"
          />
          <AutoComplete
            v-if="groupBy.kind == 'identity'"
            :inputId="`${groupBy.kind}-${groupBy.key}`"
            v-model="selector.matchIdentity[groupBy.key]"
            size="small"
            @complete="autocomplete(groupBy)"
          />
          <label :for="field.name">{{ field.name }}</label>
        </IftaLabel>
      </template>
    </template>
  </InputGroup>
</template>
