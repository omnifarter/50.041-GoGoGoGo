import { Button, Form, Col, Row, Container } from "react-bootstrap";
import { useEffect, useState } from "react";
import {getBook} from '../../../helpers/APIs'
function GetBook({setResponse}) {

    const [bookId, setBookId] = useState()

    const getBookTest = async () => {
        const data = await getBook(bookId)
        setResponse(data)
    }

    const bookIdHandler = (e) => {
        setBookId(e.target.value);
    };

    return (                        
        <>
            <h3 className="Library-title">Get Book</h3>
            <div>
                <Form>
                    <Form.Group className="mb-6">
                        <Form.Label>Book ID</Form.Label>
                        <Form.Control type="number" placeholder="Enter Book ID"  onChange={bookIdHandler} required/>
                    </Form.Group>
                <br />
                <Button onClick={getBookTest} style={{width:'100%'}}>Get Book Value</Button>
                </Form>
            </div>
        </>
);
}

export default GetBook;