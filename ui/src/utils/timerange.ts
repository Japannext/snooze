// Return a default start-end of time range (default to last 24h)
export function defaultRangeMillis(): [number, number] {
  var end = Date.now()
  var start = end - (24 * 3600 * 1000)
  return [start, end]
}
