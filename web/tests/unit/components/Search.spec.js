import { mount, shallowMount } from '@vue/test-utils'
import { render, fireEvent } from '@testing-library/vue'
import { config } from '@vue/test-utils'

import Search from '@/components/Search.vue'
import CoreuiVue from '@coreui/vue'

config.global.plugins = [CoreuiVue]

function vueHandler(vue) {
  vue.use(CoreuiVue)
  return {}
}

describe('Search.vue', () => {
  it('renders', () => {
    const wrapper = mount(Search)
  })

  it('clear', async () => {
    const { getByText, getByPlaceholderText, emitted } = render(Search, {props: {modelValue: "mytext"}}, vueHandler)
    const searchBar = getByPlaceholderText('Search')
    const button = getByText('Clear')
    await fireEvent.click(button)
    expect(emitted()).toHaveProperty('clear')
  })
})
