//import { colors } from '@/objects/Field.yaml'
import { getStyle } from '@coreui/utils/src'

const brandPrimary = getStyle('--primary') || '#304ffe'
const brandSecondary = getStyle('--secondary') || '#c4cfd4'
const brandSuccess = getStyle('--success') || '#00c853'
const brandInfo = getStyle('--info') || '#2196f3'
const brandWarning = getStyle('--warning') || '#ffc107'
const brandDanger = getStyle('--danger') || '#f44336'
const brandLight = getStyle('--light') || '#fafafa'
const brandDark = getStyle('--dark') || '#212121'
const brandTertiary = getStyle('--tertiary') || '#aa00ff'
const brandQuaternary = getStyle('--quaternary') || '#ff6d00'

type ThemeColor = 'primary'|'secondary'|'success'|'info'|'warning'|'danger'|'light'|'dark'|'tertiary'|'quaternary'

export const PALETTE = new Map<ThemeColor, string>([
  ['primary', brandPrimary],
  ['secondary', brandSecondary],
  ['success', brandSuccess],
  ['info', brandInfo],
  ['warning', brandWarning],
  ['danger', brandDanger],
  ['light', brandLight],
  ['dark', brandDark],
  ['tertiary', brandTertiary],
  ['quaternary', brandQuaternary],
])

/**
export let theme_colors = {
	primary: brandPrimary,
	secondary: brandSecondary,
	success: brandSuccess,
	info: brandInfo,
	warning: brandWarning,
	danger: brandDanger,
	light: brandLight,
	dark: brandDark,
	tertiary: brandTertiary,
	quaternary: brandQuaternary,
}
**/

/** Convert a hex code color to a rgba code
**/
/**
export function hexToRgba(color: HexColor, opacity = 100): RgbaColor {
  if (typeof color === 'undefined') {
    throw new TypeError('Hex color is not defined')
  }
  const hex = color.match(/^#(?:[0-9a-f]{3}){1,2}$/i)
  if (!hex) {
    throw new Error(`${color} is not a valid hex color`)
  }
  let r
  let g
  let b
  if (color.length === 7) {
    r = parseInt(color.slice(1, 3), 16)
    g = parseInt(color.slice(3, 5), 16)
    b = parseInt(color.slice(5, 7), 16)
  } else {
    r = parseInt(color.slice(1, 2).repeat(2), 16)
    g = parseInt(color.slice(2, 3).repeat(2), 16)
    b = parseInt(color.slice(3, 4).repeat(2), 16)
  }

  return `rgba(${r}, ${g}, ${b}, ${opacity / 100})`
}
**/

/** Generate the CSS style definition for a given color in hexadecimal code.
 * @param {string} hexcolor
 & @returns {string} The style to apply to HTML component (works with CBadge).
**/
/**
export function genColor(hexcolor) {
  if (hexcolor) {
    var color = hexToRgba(hexcolor)
    var fontcolor = (color.r*0.299 + color.g*0.587 + color.b*0.114) > 186 ? '#4f5d73' : '#ffffff'
    return 'background-color: ' + hexcolor + ' !important; border-color: ' + hexcolor +' !important; color: ' + fontcolor + ' !important'
  } else {
    return ''
  }
}
**/

/**
export function gen_color_outline(hexcolor, borderWidth = 1) {
  if (hexcolor) {
    return 'background-color: #fff !important; border-style: solid; border-color: ' + hexcolor + ' !important; color: #3c4b64 !important; border-width: ' + borderWidth + 'px !important'
  } else {
    return ''
  }
}
**/

/**
export function get_color(field) {
  if(field in colors) {
    return colors[field]
  } else {
    return 'info'
  }
}
**/

/** Generate a palette of colors
**/
/**
export function genPalette(n, m=7) {
 var palette = []
 var reference = [brandInfo, brandDanger, brandSuccess, brandWarning, brandPrimary, brandTertiary, brandQuaternary, brandLight, brandDark]
 var ref_len = Math.min(reference.length, m)
 for ( var i = 0; i < n; i++) {
   palette.push(reference[i%ref_len])
 }
 return palette
}
**/
