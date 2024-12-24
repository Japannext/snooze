import { defineConfig } from '@rsbuild/core'
import { pluginVue } from '@rsbuild/plugin-vue'
import { pluginTypeCheck } from "@rsbuild/plugin-type-check"

export default defineConfig({
  plugins: [pluginVue(), pluginTypeCheck()],
  source: {
    entry: {
      index: './src/main.ts',
    },
  },
});
