import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { XDefaultLayout, XCenteredLayout } from "./layouts";

const XCallback = () => import("./views/auth/XCallback.vue");

const routes: RouteRecordRaw[] = [
  {
    name: "menu",
    path: "/",
    component: XDefaultLayout,
    children: [
      {
        name: "dashboard",
        path: "/dashboard",
        component: () => import("./views/XDashboard.vue"),
      },
      {
        name: "alerts",
        path: "/alerts",
        component: () => import("./views/XAlerts.vue"),
      },
      {
        name: "logs",
        path: "/logs",
        component: () => import("./views/XLogs.vue"),
      },
      {
        name: "snoozes",
        path: "/snoozes",
        component: () => import("./views/XSnoozes.vue"),
      },
      {
        name: "notifications",
        path: "/notifications",
        component: () => import("./views/XNotifications.vue"),
      },
      {
        name: "escalations",
        path: "/escalations",
        component: () => import("./views/XEscalations.vue"),
      },
      {
        name: "tags",
        path: "/tags",
        component: () => import("./views/XTags.vue"),
      },
      {
        name: "pipeline",
        path: "/pipeline",
        children: [
          {
            name: "mappings",
            path: "mappings",
            component: () => import("./views/XPipelineMappings.vue"),
          },
          {
            name: "transforms",
            path: "transforms",
            component: () => import("./views/XPipelineTransforms.vue"),
          },
          {
            name: "silences",
            path: "silences",
            component: () => import("./views/XPipelineSilences.vue")
          },
          {
            name: "groupings",
            path: "groupings",
            component: () => import("./views/XPipelineGroupings.vue"),
          },
          {
            name: "notifications",
            path: "notifications",
            component: () => import("./views/XPipelineNotifications.vue"),
          },
          {
            name: "ratelimits",
            path: "ratelimits",
            component: () => import("./views/XPipelineRatelimits.vue"),
          },
          {
            name: "profiles",
            path: "profiles",
            component: () => import("./views/XPipelineProfiles.vue"),
          },
          {
            name: "Services",
            path: "services",
            component: () => import("./views/XPipelineServices.vue"),
          },
        ],
      },
    ],
  },
  {
    name: "centered",
    path: "/",
    component: XCenteredLayout,
    children: [
      {
        name: "login",
        path: "/login",
        component: () => import("./views/auth/XLogin.vue"),
      },
      { name: "callback", path: "/callback/:provider", component: XCallback },
    ],
  },
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
});
