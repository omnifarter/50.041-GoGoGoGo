import React, { useState } from 'react'
import { Modal, Button, Form } from 'react-bootstrap';

import './Borrow.css';

const Borrow = ({ show, book, closeBorrow, confirmBorrow }) => {
  const [id, setId] = useState(-1)

  return (
    <Modal show={show} onHide={closeBorrow}>
      <Modal.Header>
        <Modal.Title>Borrow Book?</Modal.Title>
      </Modal.Header>

      <Modal.Body>
        <p>{`Book Title: ${book.title}`}</p>
        <Form.Group className="mb-3" controlId="library-ID">
          <Form.Label>Enter your Library ID</Form.Label>
          <Form.Control placeholder="Library ID" onChange={(e) => setId(parseInt(e.target.value, 10))} />
        </Form.Group>

      </Modal.Body>

      <Modal.Footer>
        <Button variant="secondary" onClick={closeBorrow}>Cancel</Button>
        <Button variant="primary" onClick={() => confirmBorrow(id)}>Confirm</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default Borrow