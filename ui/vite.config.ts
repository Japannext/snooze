import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from 'unplugin-vue-components/vite';
import {PrimeVueResolver} from '@primevue/auto-import-resolver';
import * as path from "path";
import tailwindcss from "@tailwindcss/vite"

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: [{ find: "@", replacement: path.resolve(__dirname, "src") }],
  },
  plugins: [
    vue(),
    Components({
      resolvers: [
        PrimeVueResolver(),
      ],
    }),    
  ],
  server: {
    allowedHosts: true,
  }
});
