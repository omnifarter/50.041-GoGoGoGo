import { Button } from "react-bootstrap";
import { getAllBooks, borrowBook, getBook } from "../../helpers/APIs"
function Test() {
    
    const getBook0 = async () => {
        const response = await getBook(0)

        console.log(response)
    }

    const getAllBooksTest = async () => {
        const response = await getAllBooks()
        
        console.log(response)
    }

    const borrowBookTest = async () => {
        const response = await borrowBook(0,0)
        console.log(response)
    }
    
    return (
        <div>
            <Button onClick={getAllBooksTest}>Get all books</Button>
            <Button onClick={borrowBookTest}>Borrow Book 0</Button>
            <Button onClick={getBook0}>Get Book 0</Button>
        </div>
    );
}

export default Test;