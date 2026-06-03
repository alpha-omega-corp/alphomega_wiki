export default defineEventHandler(async (event) => {
    const id = getRouterParam(event, 'id')
    return $fetch(`http://localhost:1000/page/${id}`, {
        method: 'DELETE',
    })
})