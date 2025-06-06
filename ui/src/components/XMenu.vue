<script setup lang="ts">
import { ref } from "vue";
import type { MenuItem } from "primevue/menuitem";
import { Icon } from "@vicons/utils";
import logo from "@/assets/logo.png";
import {
  AlertSnooze20Regular,
  ClipboardList,
  DoNotDisturbOnTotalSilenceFilled,
  Dashboard,
  DataThresholdingFilled,
  LayerGroup,
  List,
  Megaphone,
  Notifications,
  PhoneAlt,
  Pricetag,
  RouterVoice,
  SailingRound,
  TaskSettings,
  Table,
  VolumeMute,
  TransformFilled,
} from "@/icons";

const state = ref({});

const items: MenuItem[] = [
  {
    label: "Dashboard",
    icon: Dashboard,
    route: "/dashboard",
  },
  {
    label: "Alerts",
    icon: Megaphone,
    route: "/alerts",
  },
  {
    label: "Logs",
    icon: List,
    route: "/logs",
  },
  {
    label: "Snooze",
    icon: AlertSnooze20Regular,
    route: "/snoozes",
  },
  {
    label: "Notifications",
    icon: Notifications,
    route: "/notifictions",
  },
  {
    label: "Escalations",
    icon: PhoneAlt,
    route: "/escalations",
  },
  {
    label: "Tags",
    icon: Pricetag,
    route: "/tags",
  },
  {
    label: "Pipeline",
    icon: TaskSettings,
    items: [
      {
        label: "Transforms",
        icon: TransformFilled,
        route: "/pipeline/transforms",
      },
      {
        label: "Mappings",
        icon: Table,
        route: "/pipeline/mappings",
      },
      {
        label: "Groupings",
        icon: LayerGroup,
        route: "/pipeline/groupings",
      },
      {
        label: "Profiles",
        icon: ClipboardList,
        route: "/pipeline/profiles",
      },
      {
        label: "Services",
        icon: SailingRound,
        route: "/pipeline/services",
      },
      {
        label: "Silences",
        icon: VolumeMute,
        route: "/pipeline/silences",
      },
      {
        label: "Ratelimits",
        icon: DataThresholdingFilled,
        route: "/pipeline/ratelimits",
      },
      {
        label: "Notifications",
        icon: RouterVoice,
        route: "/pipeline/notifications",
      },
    ],
  },
];
</script>

<template>
  <div>
    <Menubar :model="items" :state="state">
      <template #start>
        <img :src="logo" :height="40" />
      </template>
      <template #item="{ item, props, hasSubmenu }">
        <router-link
          v-if="item.route"
          v-slot="{ href, navigate }"
          :to="item.route"
          custom
        >
          <a v-ripple :href="href" v-bind="props.action" @click="navigate">
            <Icon>
              <component :is="item.icon" />
            </Icon>
            <span>{{ item.label }}</span>
          </a>
        </router-link>
        <a
          v-else
          v-ripple
          :href="item.url"
          :target="item.target"
          v-bind="props.action"
        >
          <Icon>
            <component :is="item.icon" />
          </Icon>
          <span>{{ item.label }}</span>
          <span v-if="hasSubmenu" class="pi pi-fw pi-angle-down" />
        </a>
      </template>
      <template #end> {{ state }}</template>
    </Menubar>
  </div>
</template>

<style lang="scss">
.layout-topbar {
  position: fixed;
  top: 0;
  left: 0;
  display: flex;
}
</style>
