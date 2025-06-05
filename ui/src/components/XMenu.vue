<script setup lang="ts">
import type { MenuItem } from "primevue/menuitem";
import { Icon } from "@vicons/utils";
import logo from "@/assets/logo.png";
import {
  Dashboard,
  Megaphone,
  List,
  AlertSnooze20Regular,
  Notifications,
  PhoneAlt,
  Pricetag,
  TaskSettings,
} from "@/icons";

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
    route: "pipeline",
  },
];
</script>

<template>
  <div>
    <Menubar :model="items">
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
          <span :class="item.icon" />
          <span>{{ item.label }}</span>
          <span v-if="hasSubmenu" class="pi pi-fw pi-angle-down" />
        </a>
      </template>
      <template #end> test... </template>
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
