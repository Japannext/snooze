const merge = require('deepmerge')
const defaultPreset = require('@vue/cli-plugin-unit-jest/presets/typescript-and-babel/jest-preset')

const yamlPreset = {
  transform: {"\\.yaml$": "jest-transform-yaml"},
  moduleFileExtensions: ["yaml"],
}

const customConfig = {
  globals: {
    "ts-jest": {
      compiler: "ttypescript",
    },
  },
}

let config = merge.all([defaultPreset, yamlPreset, customConfig])
config.transformIgnorePatterns = [
  '/node_modules/(?!v-code-diff).+\\.js$',
]

module.exports = config
