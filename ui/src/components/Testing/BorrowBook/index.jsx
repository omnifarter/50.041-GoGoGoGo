import { Button, Form, Col, Row, Container } from "react-bootstrap";
import { useEffect, useState } from "react";
import {borrowBook} from '../../../helpers/APIs'

function BorrowBook({setResponse}) {

    const [bookId, setBookId] = useState()
    const [userId, setUserId] = useState()

    const borrowBookTest = async () => {
        let data = await borrowBook(bookId,userId)
        setResponse(data)
    }

    const userIdHandler = (e) => {
        setUserId(e.target.value);
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
                        <Form.Label>User ID</Form.Label>
                        <Form.Control type="number" placeholder="Enter User ID"  onChange={userIdHandler} required/>
                    </Form.Group>
                    <br />
                    <Button onClick={borrowBookTest}>Borrow Book</Button>
                </Form>
            </div>
        </>
);
}

export default BorrowBook;