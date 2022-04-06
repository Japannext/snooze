import VueLogger from 'vuejs-logger'

const config = {
  isEnabled: true,
  logLevel: 'debug',
}

export const log = new VueLogger(config)
