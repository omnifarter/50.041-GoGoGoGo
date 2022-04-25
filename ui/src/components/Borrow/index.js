import React, { useState } from 'react'
import { Modal, Button, Form } from 'react-bootstrap';

import './Borrow.css';

const Borrow = ({ show, book, closeBorrow, confirmBorrow }) => {

  return (
    <Modal show={show} onHide={closeBorrow}>
      <Modal.Header>
        <Modal.Title>Borrow Book?</Modal.Title>
      </Modal.Header>

      <Modal.Body style={{margin:'auto',textAlign:'center'}}>
        <p>{`Book Title: ${book.Title}`}</p>
        <img src={book.Img_url}/>
      </Modal.Body>

      <Modal.Footer>
        <Button variant="secondary" onClick={closeBorrow}>Cancel</Button>
        <Button variant="primary" onClick={() => confirmBorrow()}>Confirm</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default Borrow