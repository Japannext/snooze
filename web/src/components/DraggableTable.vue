<template>
  <div class="animated fadeIn">
    <CForm @submit.prevent="" class="pt-0 px-0 pb-0">
      <Search @search="search" v-model="search_value" @clear="search_clear" ref='search'>
        <template #search_buttons>
          <!-- Slots for placing additional buttons in the header of the table -->
          <template v-if="Array.isArray(selected) && selected.length > 0">
            <slot name="selected_buttons"></slot>
            <CButton
              color="danger"
              @click="modal_show('delete', selected)"
            >Delete selection {{selected.map(i => i.item['name'])}}</CButton>
          </template>
          <slot name="head_buttons"></slot>
          <CButton color="success" @click="modal_show('new')">New</CButton>
          <CButton @click="refresh_table(true)" color="secondary" style="border-bottom-right-radius: 0"><i class="la la-refresh la-lg"></i></CButton>
        </template>
      </Search>
    </CForm>
    <div class="border" style='font-weight:bold'>
      <div class="d-flex">
        <div style="width: auto" class="align-middle p-1 singleline">
          <i class="la la-bars la-lg pe-1" style="visibility:hidden"></i>
          <input type="checkbox" class="pointer mx-1 me-1" :checked="selected.length == items.length" @change="toggle_check_all">
        </div>
        <div v-for="field in fields" :class="field.tdClass" :style="field.tdStyle">
          {{ capitalizeFirstLetter(field.label != undefined ? field.label : field.key) }}
        </div>
      </div>
    </div>
    <div class="border striped-bg p-2" style='font-weight:bold' v-if="!items.length">
      <slot name="no-items-view">
        <div class="text-center my-0">
          <h3 v-if="!is_busy" class="mb-0">
              No items
            <i
              width="20"
              class="la la-ban text-danger mb-2"
            ></i>
          </h3>
          <h3 class="mb-0" v-else>
            Loading...
          </h3>
        </div>
      </slot>
    </div>
    <Draggable ref="draggable" :flatData="items" idKey="nid" parentIdKey="pid" @drop-change="drag_end" triggerClass="can-drag">
      <template v-slot="{ node, tree }">
        <div :class="['border', {'striped-bg':node.striped}]" @mouseover="node.hover = true" @mouseleave="node.hover = false" @contextmenu="contextMenu(node.item, $event)" v-contextmenu:contextmenu>
          <div class="d-flex">
            <div style="width: auto" class="align-middle p-1 singleline m-auto">
              <i :class="['la la-bars la-lg can-drag pe-1', tree.dragging ? 'grabbing' : 'grab']"></i>
              <input type="checkbox" class="pointer ms-1 me-1" :checked="dig(node, '_checked')" @change="check(node)">
            </div>
            <div v-for="field in fields" :class="field.tdClass" :style="field.tdStyle">
              <Condition v-if="field.key == 'condition'" :data="dig(node.item, 'condition')" />
              <Modification v-else-if="field.key == 'modifications'" :data="dig(node.item, 'modifications')" />
              <span v-else>{{ dig(node.item, field.key) }}</span>
            </div>
            <div :class="['float-right', 'position-relative', {'d-none': !node.hover}]">
              <div style="position: absolute; right: 0px; top: 50%; transform: translateY(-50%)">
                <CButtonGroup role="group">
                  <CButton color="secondary" size="sm" @click="toggleDetails(node.item, $event)">
                    <i v-if="Boolean(dig(node.item, '_showDetails'))" class="la la-angle-up la-lg"></i>
                    <i v-else class="la la-angle-down la-lg"></i>
                  </CButton>
                  <CButton size="sm" @click="modal_show('new', [{'parent': dig(node.item, 'uid')}])" color="success" v-c-tooltip="{content: 'Add child'}"><i class="la la-plus la-lg"></i></CButton>
                  <CButton size="sm" @click="modal_show('edit', [node.item])" color="primary" v-c-tooltip="{content: 'Edit'}"><i class="la la-pencil-alt la-lg"></i></CButton>
                  <CButton size="sm" @click="modal_show('delete', [node.item])" color="danger" v-c-tooltip="{content: 'Delete'}"><i class="la la-trash la-lg"></i></CButton>
                </CButtonGroup>
              </div>
            </div>
          </div>
        </div>
        <CCard v-if="Boolean(dig(node.item, '_showDetails'))">
          <CRow class="m-0">
            <CCol class="p-2">
              <slot name="info" v-bind="node" />
              <Info :myobject="node.item" :excluded_fields="info_excluded_fields" />
            </CCol>
            <slot name="details_side" v-bind="node" />
          </CRow>
          <CButton size="sm" @click="toggleDetails(node.item, $event)"><i class="la la-angle-up la-lg"></i></CButton>
        </CCard>
      </template>
    </Draggable>
    <div class="d-flex align-items-center pt-2" v-if="!no_paging && nb_rows > per_page">
      <div class="me-2">
        <SPagination
          v-model:activePage="current_page"
          :pages="Math.ceil(nb_rows / per_page)"
          ulClass="m-0"
        />
      </div>
      <div>
        <CRow class="align-items-center gx-0">
          <CCol xs="auto px-1">
            <CFormLabel for="perPageSelect" class="col-form-label col-form-label-sm">Per page</CFormLabel>
          </CCol>
          <CCol xs="auto px-1">
            <CFormSelect
              v-model="per_page"
              :value="per_page"
              id="perPageSelect"
              size="sm"
            >
              <option v-for="opts in page_options" :value="opts">{{ opts }}</option>
            </CFormSelect>
          </CCol>
        </CRow>
      </div>
    </div>

    <CModal
      ref="modal"
      :visible="show_modal"
      @close="modal_clear"
      alignment="center"
      size="xl"
      backdrop="static"
    >
      <CModalHeader :class="`bg-${modal_bg_variant}`">
        <CModalTitle v-if="modal_type != 'delete'" :class="`text-${modal_text_variant}`">{{ modal_title }}</CModalTitle>
        <CModalTitle v-else-if="modal_data.length > 1" :class="`text-${modal_text_variant}`">Delete {{ modal_data.length }} items</CModalTitle>
        <CModalTitle v-else :class="`text-${modal_text_variant}`">Delete this item</CModalTitle>
      </CModalHeader>
      <CModalBody>
        <CForm v-if="modal_type != 'delete'" @submit.stop.prevent="check_form" novalidate ref="form">
          <Form v-model="modal_data" :metadata="form" :footer_metadata="form_footer"/>
        </CForm>
        <p v-else-if="modal_data.length > 1">This operation cannot be undone. Are you sure?</p>
        <p v-else>{{ modal_data }}</p>
      </CModalBody>
      <CModalFooter>
        <CButton @click="modal_clear" color="secondary">Cancel</CButton>
        <CButton @click="modal_submit" :color="modal_bg_variant">OK</CButton>
      </CModalFooter>
    </CModal>

    <v-contextmenu ref="contextmenu" @show="store_selection">
      <v-contextmenu-item @click="copy_browser" v-if="selectedText">
        <i class="la la-copy la-lg"></i> Copy
      </v-contextmenu-item>
      <v-contextmenu-item @click="context_search" v-if="selectedText">
        <i class="la la-search la-lg"></i> Search
      </v-contextmenu-item>
      <v-contextmenu-submenu title="To Clipboard">
        <template v-slot:title><i class="la la-clipboard la-lg"></i> To Clipboard</template>
        <v-contextmenu-item @click="copy_clipboard(itemCopy, fields, $event)" method="yaml">
          As YAML
        </v-contextmenu-item>
        <v-contextmenu-item @click="copy_clipboard(itemCopy, fields, $event)" method="yaml" full="true">
          As YAML (Full)
        </v-contextmenu-item>
        <v-contextmenu-divider />
        <v-contextmenu-item @click="copy_clipboard(itemCopy, fields, $event)" method="json">
          As JSON
        </v-contextmenu-item>
        <v-contextmenu-item @click="copy_clipboard(itemCopy, fields, $event)" method="json" full="true">
          As JSON (Full)
        </v-contextmenu-item>
        <v-contextmenu-divider />
        <v-contextmenu-item v-for="field in fields.filter(field => field.key != 'button' && field.key != 'select')" :key="field.key" @click="copy_clipboard(itemCopy, fields, $event)" method="simple" :field="field.key">
          {{ capitalizeFirstLetter(field.key) }}
        </v-contextmenu-item>
      </v-contextmenu-submenu>
    </v-contextmenu>
  </div>
</template>

<script>

import dig from 'object-dig'
import '@he-tree/vue3/dist/he-tree-vue3.css'
import { Draggable } from '@he-tree/vue3'
import { get_data, add_items, update_items, delete_items, capitalizeFirstLetter, to_clipboard, copy_clipboard } from '@/utils/api'
import Form from '@/components/Form.vue'
import Search from '@/components/Search.vue'
import Condition from '@/components/Condition.vue'
import Modification from '@/components/Modification.vue'
import Field from '@/components/Field.vue'
import DateTime from '@/components/DateTime.vue'
import Info from '@/components/Info.vue'
import ColorBadge from '@/components/ColorBadge.vue'
import SPagination from '@/components/SPagination.vue'

export default {
  components: {
    Draggable,
    Form,
    Search,
    Condition,
    Modification,
    Field,
    DateTime,
    Info,
    ColorBadge,
    SPagination,
  },
  props: {
    endpoint_prop: {
      type: String,
      required: true,
    },
    default_search_prop: {type: String, default: ''},
    page_options_prop: {type: Array, default: () => ['20', '50', '100']},
    info_excluded_fields: {type: Array, default: () => []},
  },
  data () {
    return {
      busy_interval: null,
      is_busy: false,
      settings: {},
      loaded: false,
      endpoint: this.endpoint_prop,
      search_value: '',
      selected: [],
      form: {},
      form_footer: {},
      default_fields: this.fields_prop,
      fields: this.fields_prop,
      default_search: this.default_search_prop,
      per_page: this.page_options_prop[0],
      page_options: this.page_options_prop,
      nb_rows: 0,
      current_page: 1,
      items: [],
      show_modal: false,
      selectedText: '',
      itemCopy: {},
      modal_title: '',
      modal_message: null,
      modal_type: '',
      modal_bg_variant: '',
      modal_text_variant: '',
      modal_data: {},
      to_clipboard: to_clipboard,
      copy_clipboard: copy_clipboard,
      capitalizeFirstLetter: capitalizeFirstLetter,
      add_items: add_items,
      update_items: update_items,
      delete_items: delete_items,
      dig: dig,
    }
  },
  mounted () {
    this.settings = JSON.parse(localStorage.getItem(this.endpoint+'_json') || '{}')
    get_data('settings/?c='+encodeURIComponent(`web/${this.endpoint}`)+'&checksum='+(this.settings.checksum || ""), null, {}, this.load_table)
  },
  methods: {
    stripe () {
      var nodes = this.items
      nodes.forEach( (node, index) => {
        node.striped = index%2 == 0
      })
    },
    load_table (response) {
      if (response.data) {
        if (response.data.count > 0) {
          this.settings = response.data
          localStorage.setItem(this.endpoint+'_json', JSON.stringify(response.data))
        }
        var data = this.settings.data[0]
        this.form = dig(data, 'form')
        this.form_footer = dig(data, 'form_footer')
        this.endpoint = dig(data, 'endpoint') || this.endpoint
        this.fields = dig(data, 'fields')
        this.fields.forEach(field => {
          field.tdClass = (field.tdClass || []).concat(['p-1', 'border-start', 'd-flex', 'align-items-center'])
          //field.tdStyle = Object.assign(field.tdStyle || {}, {'border-left': '1px solid'})
        })
        this.default_fields = this.fields
        this.reload()
      }
    },
    reload () {
      var search = this.default_search
      if (this.$route.query.perpage !== undefined) {
        this.per_page = this.$route.query.perpage
      }
      if (this.$route.query.pagenb !== undefined) {
        this.current_page = parseInt(this.$route.query.pagenb)
      }
      if (this.$route.query.s !== undefined) {
        search = decodeURIComponent(this.$route.query.s)
      }
      this.search_value = search
      if (this.$refs.search) {
        this.$refs.search.datavalue = search
      }
      this.refresh_table()
    },
    refresh_table(feedback = false) {
      this.uncheck_all()
      this.set_busy(true)
      var query = []
      var options = {
        perpage: this.per_page,
        pagenb: this.current_page,
        asc: false,
      }
      if (this.search_value) {
        if (this.search_value[0] == '[') {
          var search_json = JSON.parse(this.search_value)
          if (search_json) {
            query = join_queries([query, search_json])
          }
        } else {
          options["ql"] = this.search_value
        }
      }
      options["orderby"] = 'name'
      get_data(this.endpoint, query, options, feedback == true ? this.feedback_then_update : this.update_table, null)
    },
    feedback_then_update(response) {
      this.$root.show_alert()
      this.update_table(response)
    },
    update_table(response) {
      this.set_busy(false)
      if (response.data) {
        this.items = []
        this.nb_rows = response.data.count
        var rows = response.data.data || []
        var parent
        var item
        rows.forEach((row, index) => {
          if ( this.items.every(x => x['uid'] != row['uid']) ) {
            item = {'item': row}
            this.items.push(item)
            item.nid = index
            if (index > 0 && row.parent == rows[index-1].uid) {
              row.pid = parent
            } else {
              parent = row.nid
            }
          }
        })
        this.stripe()
      }
      if (!this.loaded) {
        this.loaded = true
      }
    },
    uncheck_all() {
      this.selected = []
      this.$refs.draggable.nodes.forEach(node => {
        node._checked = false
      })
    },
    check_all() {
      this.selected = []
      this.$refs.draggable.nodes.forEach(node => {
        node._checked = true
        this.selected.push(node)
      })
    },
    toggle_check_all () {
      if(this.selected.length == this.items.length) {
        this.uncheck_all()
      } else {
        this.check_all()
      }
    },
    check (node) {
      var found = this.selected.indexOf(node)
      if (found >= 0) {
        this.selected.splice(found, 1)
        node._checked = false
      } else {
        this.selected.push(node)
        node._checked = true
      }
    },
    set_busy(busy) {
      if (this.busy_interval) {
        clearInterval(this.busy_interval)
        this.busy_interval = null
      }
      if (busy) {
        this.busy_interval = setInterval(() => {this.is_busy = true}, 500);
      } else {
        this.is_busy = false
      }
    },
    add_history() {
      const query = { s: (this.search_value || this.default_search || ''), perpage: this.per_page, pagenb: this.current_page }
      if (this.$route.query.s != query.s || this.$route.query.perpage != query.perpage || this.$route.query.pagenb != query.pagenb) {
        this.$router.push({ query: query })
      }
    },
    toggleDetails(row, event) {
      event.stopPropagation()
      row._showDetails = !row._showDetails
    },
    search(query) {
      this.add_history()
    },
    search_clear() {
      this.search_value = this.default_search
      this.add_history()
    },
    modal_clear() {
      this.modal_data = {}
      this.modal_title = ''
      this.modal_message = null
      this.modal_type = ''
      this.modal_bg_variant = ''
      this.modal_text_variant = ''
      this.show_modal = false
      Array.from(document.getElementsByClassName('modal')).forEach(el => el.style.display = "none")
      Array.from(document.getElementsByClassName('modal-backdrop')).forEach(el => el.style.display = "none")
    },
    modal_show (type = '', items = [{}]) {
      this.modal_type = type
      switch (this.modal_type) {
        case 'edit':
          this.modal_title = 'Edit'
          this.modal_bg_variant = 'info'
          this.modal_text_variant = 'white'
          this.modal_data = JSON.parse(JSON.stringify(items[0]))
          break
        case 'delete':
          this.modal_title = 'Delete'
          this.modal_bg_variant = 'danger'
          this.modal_text_variant = 'white'
          this.modal_data = items
          break
        default:
          this.modal_title = 'New'
          this.modal_bg_variant = 'success'
          this.modal_text_variant = 'white'
          this.modal_data = items[0]
      }
      this.show_modal = true
    },
    check_form (node) {
      return (node.$el.getElementsByClassName('form-control is-invalid').length + node.$el.getElementsByClassName('has-error').length) == 0
    },
    modal_submit (bvModalEvt, endpoint = this.endpoint) {
      bvModalEvt.preventDefault()
      if (this.$refs.form && !this.check_form(this.$refs.form)) {
        this.$root.text_alert('Form is invalid', 'danger')
        return
      }
      this.set_busy(true)
      this.$nextTick(() => {
        this.modal_clear()
      })
      if (!Array.isArray(this.modal_data)) {
        this.modal_data = [this.modal_data]
      }
      switch (this.modal_type) {
        case 'edit':
          this.update_items(this.endpoint, this.modal_data, this.submit_callback)
          break
        case 'delete':
          this.delete_items(this.endpoint, this.modal_data, this.submit_callback)
          break
        default:
          this.add_items(this.endpoint, this.modal_data, this.submit_callback)
      }
    },
    submit_callback (response) {
      this.set_busy(false)
      this.refresh_table()
    },
    contextMenu (item, e) {
      this.itemCopy = item
      this.$refs.contextmenu.hide()
      this.$refs.contextmenu.show({top: e.pageY, left: e.pageX})
    },
    store_selection() {
      this.selectedText = window.getSelection().toString()
    },
    copy_browser(event) {
      this.to_clipboard(this.selectedText)
    },
    context_search(event) {
      if (this.selectedText != '') {
        this.search_value = this.selectedText
        this.$refs.search.datavalue = this.selectedText
        this.search(this.selectedText)
      }
    },
    drag_end () {
      this.stripe()
    },
  },
  watch: {
    current_page: function() {
      this.add_history()
    },
    per_page: function() {
      this.add_history()
    },
    $route() {
      if (this.loaded && this.$route.path == `/${this.endpoint}`) {
        this.$nextTick(this.reload);
      }
    }
  },
}
</script>
