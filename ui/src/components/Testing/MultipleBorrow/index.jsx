import { Button, Form } from "react-bootstrap";
import { useState } from "react";
import { getBook, borrowBook} from '../../../helpers/APIs'

function MultipleBorrow({setResponse}) {

    const [bookId, setBookId] = useState()
    const [userId, setUserId] = useState()
    const [userId2, setUserId2] = useState()

    const borrowBookTest = async () => {
        borrowBook(bookId,userId)
        borrowBook(bookId,userId2)
        const data = await getBook(bookId)
        setResponse(data)
    }

    const userIdHandler = (e) => {
        setUserId(e.target.value);
    };

    const userId2Handler = (e) => {
        setUserId2(e.target.value);
    };

    const bookIdHandler = (e) => {
        setBookId(e.target.value);
    };

    return (                        
        <>
            <h3 className="Library-title">Borrow Book</h3>
            <div>
                <Form>
                    <Form.Group className="mb-6">
                        <Form.Label>Book ID</Form.Label>
                        <Form.Control type="number" placeholder="Enter Book ID"  onChange={bookIdHandler} required/>
                        <br />
                        <Form.Label>User ID 1</Form.Label>
                        <Form.Control type="number" placeholder="Enter User ID"  onChange={userIdHandler} required/>
                        <br />
                        <Form.Label>User ID 2</Form.Label>
                        <Form.Control type="number" placeholder="Enter User ID"  onChange={userId2Handler} required/>
                    </Form.Group>
                    <br />
                    <Button onClick={borrowBookTest}>Borrow Book</Button>
                </Form>
            </div>
        </>
);
}

export default MultipleBorrow;