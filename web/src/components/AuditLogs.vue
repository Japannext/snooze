<template>
  <CCol v-if="curatedLogs.length > 0" class="p-2">
    <CCard no-body>
      <CCardHeader class="text-center" style="font-weight:bold">Audit logs</CCardHeader>
      <CCardBody class="p-2">
        <template v-if="authorized">
          <CRow v-for="(auditLog, index) in curatedLogs" :key="auditLog.uid" class="m-0">
            <CCol class="p-0">
              <CCard class="mb-2">
                <CRow class="m-2">
                  <!-- Icon -->
                  <CCol xs="auto" class="p-0">
                    <div :class="`bg-${auditLog.color} text-white rounded p-2`">
                      <i :class="auditLog.icon"></i>
                    </div>
                  </CCol>

                  <!-- Text -->
                  <CCol xs="auto">
                    <!-- First line -->
                    <CRow class="d-flex">
                      <CCol xs="auto">
                        <span :class="`fw-bold text-${auditLog.color}`">{{ auditLog.name }}</span> by
                        <CBadge :color="auditLog.methodColor">{{ auditLog.method }}</CBadge> /
                        <span class="fw-bold" style="font-size: 1.0rem">{{ auditLog.username }}</span>
                        @<DateTime class="fst-italic muted" :date="auditLog.timestamp" show_secs="true" />
                      </CCol>
                    </CRow>

                    <!-- Second line -->
                    <CRow class="d-flex">
                      <CCol>
                        Modified fields:
                        <template v-for="(entry, i) in auditLog.quickSummary" :key="entry.name">
                          <span :class="`fw-bold text-${entry.color}`">{{ entry.symbol }}{{ entry.name }}</span>
                          <span v-if="i < auditLog.quickSummary.length-1">, </span>
                        </template>
                        <span v-if="auditLog.summaryCount > summaryMax">, ...</span>
                      </CCol>
                    </CRow>
                  </CCol>

                  <!-- Button -->
                  <CCol class="d-flex justify-content-end align-middle p-0">
                    <CButton color="warning" class="ml-auto mr-0" @click="modalShow(index)">
                      Inspect
                      <i class="la la-lg la-search-plus"></i>
                    </CButton>
                  </CCol>
                </CRow>
              </CCard>
            </CCol>
          </CRow>

          <!-- Pagination -->
          <div>
            <CButtonToolbar role="group">
              <div class="d-flex ms-auto me-2 align-items-center">
                <div class="me-2">
                  <SPagination
                    v-model:activePage="currentPage"
                    :pages="Math.ceil(numberOfRows / perPage)"
                    ul-class="m-0"
                  />
                </div>
                <div>
                  <CRow class="align-items-center gx-0">
                    <CCol xs="auto px-1">
                      <CFormLabel for="perPageSelect" class="col-form-label col-form-label-sm">Per page</CFormLabel>
                    </CCol>
                    <CCol xs="auto px-1">
                      <CFormSelect
                        id="perPageSelect"
                        v-model="perPage"
                        size="sm"
                      >
                        <option v-for="opt in pageOptions" :key="opt" :value="opt">{{ opt }}</option>
                      </CFormSelect>
                    </CCol>
                  </CRow>
                </div>
              </div>
              <CButtonGroup role="group">
                <CButton v-c-tooltip="{content: 'Refresh'}" size="sm" color="secondary" @click="refresh">
                  <i class="la la-refresh la-lg"></i>
                </CButton>
              </CButtonGroup>
            </CButtonToolbar>
          </div>
        </template>
        <template v-else>
          You don't have permission to see audit logs.<br />
          You need one of the following permissions: {{ AUTHORIZED_PERMISSIONS.join(', ') }}
        </template>
      </CCardBody>
    </CCard>
    <AuditModal
      ref="diffModal"
      :collection="collection"
      :object-id="uid"
      :audit-logs="curatedLogs"
    />
  </CCol>
</template>

<script lang="ts">
const ACTIONS = {
  added: {name: 'CREATE', icon: 'la la-plus la-2x', color: 'success', symbol: '+'},
  replaced: {name: 'REPLACED', icon: '', color: 'primary', symbol: '~'},
  updated: {name: 'UPDATE', icon: 'la la-edit la-2x', color: 'primary', symbol: '~'},
  rejected: {name: 'ERROR', icon: 'la la-bug la-2x', color: 'danger'},
  removed: {name: 'DELETE', icon: 'la la-cross la-2x', color: 'danger', symbol: '-'},
  undefined: {name: 'UNKNOWN', icon: 'la la-bug la-2x', color: 'danger', symbol: ''},
}

const AUTHORIZED_PERMISSIONS = ['rw_all', 'ro_all', 'rw_audit', 'ro_audit']

import { defineComponent, PropType } from 'vue'

import { colors } from '@/objects/Field.yaml'

import { api2 } from '@/api2'

import AuditModal from '@/components/AuditModal.vue'
import DateTime from '@/components/DateTime.vue'
import SPagination from '@/components/SPagination.vue'

import { AuditItem, AuditSummary, SummaryEntry, AuditMetadata, DatabaseItem } from '@/utils/types'

export default defineComponent({
  name: 'AuditLogs',
  components: {
    DateTime,
    SPagination,
    AuditModal,
  },
  props: {
    collection: {type: String, required: true},
    object: {type: Object as PropType<DatabaseItem>, required: true},
  },
  data () {
    return {
      ACTIONS: ACTIONS,
      COLORS: colors,
      AUTHORIZED_PERMISSIONS: AUTHORIZED_PERMISSIONS,
      auditEndpoint: api2.endpoint('audit'),
      objectId: this.object.uid,
      auditLogs: [],
      curatedLogs: [],
      summaryMax: 3,
      perPage: '5',
      pageOptions: ['5', '10', '20'],
      numberOfRows: 0,
      currentPage: 1,
      previousLog: {},
    }
  },
  watch: {
    currentPage: function() {
      this.refresh()
    },
    perPage: function() {
      this.refresh()
    },
  },
  mounted() {
    this.authorized = this.isAuthorized()
    this.refresh()
    //this.auto_interval = setInterval(this.refresh, 2000);
  },
  beforeUnmount() {
    if (this.auto_interval) {
      clearInterval(this.auto_interval)
    }
  },
  methods: {
    refresh() {
      let query = ['AND',
        ['=', 'object_id', this.object.uid],
        ['=', 'collection', this.collection],
      ]
      let options = {
        perpage: this.perPage,
        pagenb: this.currentPage,
        asc: false,
        orderby: 'timestamp',
      }
      this.auditEndpoint.find(query, options)
        .then((results: AuditItem[]) => {
          this.auditLogs = results
          this.curatedLogs = this.auditLogs.map(auditLog => this.appendMetadata(auditLog))
          this.numberOfRows = results.length
        })
    },
    isAuthorized(): boolean {
      var permissionString = localStorage.getItem('permissions')
      var permissions = (permissionString !== null) ? permissionString.split(',') : []

      var authorized = this.AUTHORIZED_PERMISSIONS.some(authorizedPermission => {
        return permissions.includes(authorizedPermission)
      })
      return authorized
    },
    appendMetadata(auditLog: AuditItem): AuditMetadata {
      var action = auditLog.action
      auditLog.name = this.ACTIONS[action].name
      auditLog.color = this.ACTIONS[action].color
      auditLog.icon = this.ACTIONS[action].icon
      auditLog.methodColor = this.COLORS[auditLog.method]
      auditLog.quickSummary = this.computeQuickSummary(auditLog.summary)
      auditLog.summaryCount = Object.keys(auditLog.summary).length
      return auditLog
    },
    computeQuickSummary(summary: AuditSummary): SummaryEntry[] {
      var quickSummary = []
      for (const [field, action] of Object.entries(summary)) {
        var symbol = this.ACTIONS[action].symbol
        var color = this.ACTIONS[action].color
        var entry = {symbol: symbol, color: color, name: field}
        quickSummary.push(entry)
      }
      return quickSummary.slice(0, this.summaryMax)
    },
    modalShow(index: number) {
      this.$refs.diffModal.show(index)
    },
  },
})
</script>
