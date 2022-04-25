import React, { useState } from 'react'
import { Modal, Button, Form } from 'react-bootstrap';

import './Borrow.css';

const Borrow = ({ show, book, closeBorrow, confirmBorrow }) => {

  return (
    <Modal show={show} onHide={closeBorrow}>
      <Modal.Header>
        <Modal.Title>{book.Title}</Modal.Title>
      </Modal.Header>

      <Modal.Body style={{margin:'auto',textAlign:'center'}}>
        <img src={book.Img_url} style={{maxWidth:'300px'}}/>
      </Modal.Body>

      <Modal.Footer>
        <Button variant="primary" style={{width:'100%'}} onClick={() => confirmBorrow()}>Borrow</Button>
        <Button variant="secondary" style={{width:'100%'}} onClick={closeBorrow}>Cancel</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default Borrow