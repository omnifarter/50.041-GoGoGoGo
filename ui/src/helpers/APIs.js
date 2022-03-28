const BACKEND_URL = "http://localhost:8080"

export const getAllBooks = async () => {
    const response = await fetch(`${BACKEND_URL}/all`,
    {
        method:"GET"
    })
    if (response.ok) {
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const borrowBook = async (bookId,userId) => {
    const response = await fetch(`${BACKEND_URL}/user/borrow`,{
        method:"POST",
        body:{
            userId,
            bookId
        }
    })
    if (response.ok) {
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}
