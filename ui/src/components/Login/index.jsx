import React, { useState } from 'react'
import { Modal, Button, Form } from 'react-bootstrap';

const Borrow = ({ show, onSetUser }) => {
  const [id, setId] = useState(null)

  return (
    <Modal show={show}>
      <Modal.Header>
        <Modal.Title>Login</Modal.Title>
      </Modal.Header>

      <Modal.Body>
        <Form.Group className="mb-3" controlId="library-ID">
          <Form.Label>To continue, please log in.</Form.Label>
          <Form.Control placeholder="User ID" onChange={(e) => setId(parseInt(e.target.value))} />
        </Form.Group>

      </Modal.Body>

      <Modal.Footer>
        <Button variant="primary" onClick={() => onSetUser(id)}>Confirm</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default Borrow