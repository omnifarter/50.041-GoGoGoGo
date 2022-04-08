const BACKEND_URL = "http://localhost:8080"

export const getAllBooks = async () => {
    const response = await fetch(`${BACKEND_URL}/all`,
    {
        method:"GET",
        // mode: 'cors',
    })
    if (response.ok) {
        const res =  await response.json()
        console.log(res)
        return res
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const getBook = async (bookId) => {
    const response = await fetch(`${BACKEND_URL}/books/get?` + new URLSearchParams({
        bookId
    }),
    {
        method:"GET",
        // mode: 'cors',

    })
    if (response.ok) {
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const borrowBook = async (bookId,userId) => {
    const response = await fetch(`${BACKEND_URL}/user/borrow`,{
        method:"PUT",
        body:JSON.stringify({
            "userId":parseInt(userId),
            "bookId":parseInt(bookId)
        })
    })
    if (response.ok) {
        console.log("Response OK")
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const getNodes = async () => {
    const response = await fetch(`${BACKEND_URL}/nodes/all`,
    {
        method:"GET",
        // mode: 'cors',

    })
    if (response.ok) {
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const addNode = async () => {
    const response = await fetch(`${BACKEND_URL}/nodes/add`,
    {
        method:"GET",
        // mode: 'cors',

    })
    if (response.ok) {
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const removeNode = async () => {
    const response = await fetch(`${BACKEND_URL}/nodes/kill`,
    {
        method:"GET",
        // mode: 'cors',

    })
    if (response.ok) {
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}

export const addBook = async (bookTitle,bookImgUrl) => {
    const response = await fetch(`${BACKEND_URL}/books/add`,{
        method:"PUT",
        body:JSON.stringify({
            "Title":bookTitle,
            "Img_url":bookImgUrl,
        })
    })
    if (response.ok) {
        console.log("Response OK")
        return await response.json()
    } else {
        throw new Error(`Error ${response.status}`)
    }
}