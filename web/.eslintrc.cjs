require("@rushstack/eslint-patch/modern-module-resolution")

module.exports = {
  root: true,
  'env': {
    browser: true,
    es2021: true,
    jest: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-recommended',
    '@vue/eslint-config-typescript/recommended',
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    'parser': '@typescript-eslint/parser',
    'js': 'espree',
    "ts": "@typescript-eslint/parser",
    'ecmaVersion': 2021,
    'extraFileExtensions': ['.vue'],
    'sourceType': 'module',
    ecmaFeatures: {jsx: true},
  },
  plugins: [
    'vue',
    '@typescript-eslint',
  ],
  rules: {
    'vue/max-attributes-per-line': ['error', {singleline: {max: 4}, multiline: {max: 1}}],
    'vue/singleline-html-element-content-newline': 'off',
    'vue/html-self-closing': 'off',
  },
}

/**
---
extends:
  - eslint:recommended
  - plugin:vue/vue3-recommended
env:
  browser: true
  es2021: true
  jest: true
extends:
parser: 'vue-eslint-parser'
parserOptions:
  parser: '@typescript-eslint/parser'
  extraFileExtensions: ['.vue']
  ecmaFeatures: {jsx: true}
  ecmaVersion: 2021
  sourceType: module
plugins:
  - vue
  - '@typescript-eslint'
rules:
  vue/max-attributes-per-line:
    - error
    - singleline: {max: 4}
      multiline: {max: 1}
  vue/singleline-html-element-content-newline: 'off'
  vue/html-self-closing: 'off'
**/
