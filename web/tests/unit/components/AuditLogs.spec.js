import { mount, config } from '@vue/test-utils'
import CoreuiVue from '@coreui/vue'

import AuditLogs from '@/components/AuditLogs.vue'

config.global.plugins = [CoreuiVue]

describe('AuditLogs.vue', () => {

  test('renders', () => {
    const rule1 = {
      uid: 'c300351d-d0f1-45fe-b535-cdd005e1e0ef',
      name: 'Rule 1',
      condition: ['=', 'a', 'x'],
      modifications: [
        ['SET', 'b', 'y'],
      ],
    }
    const props = {collection: 'rule', object: rule1}
    const wrapper = mount(AuditLogs, {props: props})
  })

})
