export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return $fetch('http://localhost:1000/pages', {
    method: 'POST',
    body,
  })
})