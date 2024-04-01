
import { NIcon } from 'naive-ui'
import { h, Component } from 'vue'

export function renderIcon (icon: Component, color?: string) {
  var args: object|null
  if (color !== undefined) {
    args = {color: color}
  } else {
    args = null
  }
  return () => h(NIcon, args, { default: () => h(icon) })
}
