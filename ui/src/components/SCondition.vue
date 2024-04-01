<script setup lang="ts">

import { h, VNodeChild } from 'vue'
import { NText, NTag, NIcon } from 'naive-ui'
import { Refresh } from '@vicons/ionicons5'

import { Condition } from '@/api'

const props = defineProps<{
  value: Condition,
}>()

function insertInBetween<T>(array: Array<T>, item: T): Array<T> {
  var res: Array<T> = []
  array.forEach((element) => {
    res.push(element)
    res.push(item)
  })
  res.pop()
  return res
}

function printField(field: string[]): VNodeChild {
  return h('span', insertInBetween(
    field.map((f) => h('span', f)),
    h(NText, {strong: true}, '.'),
  ))
}

function getKind(kind: string): VNodeChild {
  switch(kind) {
    case 'and':
      return h('span', ' & ')
    case 'or':
      return h('span', ' | ')
    case 'not':
      return h('span', 'NOT')
    case "=": case ">": case "<":
      return h('span', ` ${kind} `)
    case "!=":
      return h('span', ' ≠ ')
    case ">=":
      return h('span', ' ≥ ')
    case "<=":
      return h('span', ' ≤ ')
    case "matches":
      return h('span', ' ~ ')
    case "exists":
      return h('span', ' ?')
    default:
      return h('span', " ??? ")
  }
}

function printCondition(item: Condition): VNodeChild {
  switch(item.kind) {
    case 'always_true':
      return h(NTag, null, {
        icon: () => h(NIcon, {component: Refresh}),
        default: () => "Always true",
      })
    case 'and': case 'or':
      return h('span', insertInBetween(
        item.conditions.map((c) => printCondition(c)),
        getKind(item.kind),
      ))
    case 'not':
      return h('span', [
        getKind(item.kind),
        printCondition(item.condition),
      ])
    // operation case
    case '=': case '!=':
    case '>': case '<': case '>=': case '<=':
    case 'matches':
      return h(NTag, {}, [
        printField(item.field),
        getKind(item.kind),
        h(NText, JSON.stringify(item.value))
      ])
    case 'exists':
      return h(NTag, {}, [
        printField(item.field),
        getKind(item.kind),
      ])
    default:
      return h('span', " ??? ")
  }
}

</script>

<template>
  <component
    :is="printCondition(value)"
    v-if="value !== undefined && value.kind !== undefined"
  />
  <div v-else>Error</div>
</template>
