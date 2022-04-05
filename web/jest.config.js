const deepmerge = require('deepmerge')
const defaultPreset = require('@vue/cli-plugin-unit-jest/presets/typescript-and-babel/jest-preset')

const yamlPreset = {
  transform: {"\\.yaml$": "jest-transform-yaml"},
  moduleFileExtensions: ["yaml"],
}

const config = {
  //transformIgnorePatterns: [
  //],
  global: {
    config: {
      compilerOptions: {
        isCustomElement: tag => tag.startsWith('c'),
      }
    }
  }
}

module.exports = deepmerge(
  defaultPreset,
  yamlPreset,
  config,
)
