import { useState } from "react";
import { Button } from "react-bootstrap";
import { getAllBooks, borrowBook } from "../../helpers/APIs"
import ReactJson from 'react-json-view'

function Test() {
    
    const [response, setResponse] = useState({})
    const getAllBooksTest = async () => {
        const data = await getAllBooks()
        setResponse(data)
    }

    const borrowBookTest = async () => {
        const data = await borrowBook(0,0)
        setResponse(data)
    }
    
    return (
        <div>
            <Button onClick={getAllBooksTest}>Get all books</Button>
            <Button onClick={borrowBookTest}>Borrow Book 0</Button>
            <ReactJson src={response} />
        </div>
    );
}

export default Test;