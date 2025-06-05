<script setup lang="ts">
import { ref, defineModel, defineEmits } from "vue";
import { useMessage } from "naive-ui";
import { createSnooze, type Snooze, type Grouping } from "@/api";
import { DateTime } from "luxon";

const grouping = ref<Grouping | undefined>(undefined);

import { InformationCircleOutline } from "@/icons";
import { AutoComplete, InputGroupAddon } from "primevue";

const emit = defineEmits(["success"]);

const showAdditionalCondition = ref(false);
const show = defineModel("show", { type: Boolean, default: false });
const loading = ref(false);
const message = useMessage();
const item = ref<Snooze>(defaultValue());
const labels = ref<any>({});

function defaultValue(): Snooze {
  var now = DateTime.now();
  return {
    reason: "",
    startsAt: now.toMillis(),
    endsAt: now.plus({ hour: 1 }).toMillis(),
    groups: [],
    tags: [],
    username: "",
  };
}

function submit() {
  loading.value = true;
  console.log("submit()");
  createSnooze(item)
    .then((resp) => {
      emit("success");
      console.log("emit(`success`)");
    })
    .catch((err) => {
      console.log(`error: ${err}`);
      message.error(`failed to snooze: ${err}`);
    })
    .finally(() => {
      console.log("finally");
      loading.value = false;
      exit();
    });
}

function exit() {
  show.value = false;
  item.value = defaultValue();
}
</script>

<template>
  <Dialog
    v-model:visible="show"
    modal
    header="Snooze"
    :style="{ width: '50rem' }"
  >
    <div class="card grid grid-cols-1 mdLgrid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <XGroupSelect v-model:group="item.groups[0]" />
      </div>

      <div class="flex flex-col gap-1">
        <ToggleSwitch
          v-model="showAdditionalCondition"
          aria-labeledby="showAdditionalCondition"
        />
        <label id="showAdditionalCondition">Additional Condition</label>
      </div>

      <div class="flex flex-col gap-1">
        <XTimePeriod v-model:start="item.startsAt" v-model:end="item.endsAt" />
      </div>

      <div class="flex flex-col gap-1">
        <IftaLabel>
          <TextArea inputId="reason" v-model="item.reason" />
          <label for="reason">Reason</label>
        </IftaLabel>
      </div>
    </div>
    <template #footer>
      <Button label="Cancel" @click="exit" severity="secondary" />
      <Button type="submit" label="Snooze" :loading="loading" severity="warn" />
    </template>
  </Dialog>
</template>
