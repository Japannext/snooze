/// <reference types="@rsbuild/core/types" />

interface ImportMetaEnv {
  readonly VITE_BASE_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import Vue from 'vue'
  export default Vue
}
