import { Button, Form, Col, Row, Container } from "react-bootstrap";
import { useEffect, useState } from "react";
import {getBook,getAllBooks,borrowBook} from '../../helpers/APIs'
import ReactJson from 'react-json-view'

function Test() {
    const [noOfNodes, setNoOfNodes] = useState();
    const [response, setResponse] = useState()
    const [bookId, setBookId] = useState()
    const [bookTitleAdd, setBookTitleAdd] = useState('');
    const [bookURLAdd, setBookURLAdd] = useState('');
    const [userId, setUserId] = useState()


    const addBook = () => { }; //TODO: implement
    const removeBook = () => { }; //TODO: implement

    const getBookTest = async () => {
        const data = await getBook(bookId)
        setResponse(data)
    }

    const borrowBookTest = async () => {
        await borrowBook(bookId,userId)
    }

    const getAllBooksTest = async () => {
        const data = await getAllBooks()
        setResponse(data)
    }
    const addNode = () => { };
    const removeNode = () => { };


    const bookTitleHandler = (e) => {
        setBookTitleAdd(e.target.value);
    };

    const userIdHandler = (e) => {
        setUserId(e.target.value);
    };

    const bookURLHandler = (e) => {
        setBookURLAdd(e.target.value);
    };
    const addBookHandler = (e) => {
        e.preventDefault();
        setBookTitleAdd('');
        setBookURLAdd('');
        console.log(bookTitleAdd+','+bookURLAdd);
        return alert(bookTitleAdd + ',' + bookURLAdd)
    }

    const bookIdHandler = (e) => {
        setBookId(e.target.value);
    };
    const removeBookHandler = (e) => {
        e.preventDefault();
        return alert("Book Id: " + bookId)
    };

    return (
        <div>
            <header className="App-header">
                <h1 className="Library-title">GoGoGoGo - Test Page</h1>
                <div>No. of Nodes: {noOfNodes}</div>
            </header>
            <br />

            <div className="App-header">

                <Button variant="success" onClick={() => addNode()}>
                    Add Node
                </Button>{" "}
                <Button variant="danger" onClick={() => removeNode()}>
                    Remove Node
                </Button>{" "}
            </div>
            <Container>
                <Row>
                    <Col>
                        <h3 className="Library-title">Add Book!</h3>
                        <div className="App-header">
                            <Form onSubmit={addBookHandler}>
                                <Form.Group className="mb-6">
                                    <Form.Label>Book Title</Form.Label>
                                    <Form.Control type="text" placeholder="Book Title" value={bookTitleAdd} onChange={bookTitleHandler} required/>
                                </Form.Group>
                                <Form.Group className="mb-6">
                                    <Form.Label>Book Image URL</Form.Label>
                                    <Form.Control type="text" placeholder="www.bookimageurl.com" value={bookURLAdd} onChange={bookURLHandler} required/>
                                </Form.Group>
                                <br />
                                <Button variant="info" type="submit">
                                    Add Book
                                </Button>{" "}
                            </Form>
                        </div>
                    </Col>
                    <Col>
                        <h3 className="Library-title">Remove Book!</h3>
                        <div>
                            <Form onSubmit={removeBookHandler}>
                                <Form.Group className="mb-6">
                                    <Form.Label>Book ID</Form.Label>
                                    <Form.Control type="number" placeholder="1"  onChange={bookIdHandler} required/>
                                </Form.Group>
                                <br />
                                <Form.Group className="mb-6">
                                    <Form.Label>User ID</Form.Label>
                                    <Form.Control type="number" placeholder="1"  onChange={userIdHandler} required/>
                                </Form.Group>
                                <br />
                                <Button variant="info" type="submit">
                                    Remove Book
                                </Button>{" "}
                                <Button onClick={getBookTest}>Get Book Value</Button>
                                 <Button onClick={getAllBooksTest}>Get all books value</Button>
                                 <Button onClick={borrowBookTest}>Borrow Book</Button>
                            </Form>
                        </div>
                    </Col>
                </Row>
            </Container>

            <ReactJson src={response} />

        </div>
    );
}

export default Test;
