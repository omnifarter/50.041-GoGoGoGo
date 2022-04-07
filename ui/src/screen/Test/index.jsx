import { Button, Form, Col } from "react-bootstrap";
import { useEffect, useState } from "react";

function Test() {
    const [noOfNodes, setNoOfNodes] = useState();

    const addBook = () => { };
    const removeBook = () => { };
    const addNode = () => { };
    const removeNode = () => { };

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
            <Col >
                <div classname="App-header">
                    <Form>
                        <Form.Group className="mb-6">
                            <Form.Label>Book Title</Form.Label>
                            <Form.Control type="text" placeholder="Book Title" />
                        </Form.Group>
                        <Form.Group className="mb-6">
                            <Form.Label>Book Image URL</Form.Label>
                            <Form.Control type="text" placeholder="www.bookimageurl.com" />
                        </Form.Group>
                        <br />
                        <Button variant="info" onClick={() => addBook()}>
                            Add Book
                        </Button>{" "}
                    </Form>
                </div>
            </Col>
            <Col>
                <div classname="App-header">
                    <Form>
                        <Form.Group className="mb-6">
                            <Form.Label>Book ID</Form.Label>
                            <Form.Control type="text" placeholder="1" />
                        </Form.Group>
                        <br />
                        <Button variant="info" onClick={() => removeBook()}>
                            Remove Book
                        </Button>{" "}
                    </Form>
                </div>
            </Col>

        </div>
    );
}

export default Test;
