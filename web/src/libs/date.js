export const formatDate = (date) => {
  const y = date.getFullYear()
  const m = date.getMonth()
  const d = date.getDate()
  const h = date.getHours()
  const i = date.getMinutes()
  const s = date.getSeconds()

  return y + '-' + m + '-' + d + ' ' + h + ':' + i + ':' + s
}
