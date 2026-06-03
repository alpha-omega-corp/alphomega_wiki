export default defineEventHandler(async (event) => {
    const body = await readBody(event)
    console.log(body)
    return $fetch('http://localhost:1000/page', {
        method: 'POST',
        body,
    })
})