import React from 'react'
import { Modal, Button, Table } from 'react-bootstrap';
// import './MyBooks.css';

const MyBooks = ({show, closeMyBooks, myBooks,onReturn}) => {
  return (
    <Modal show={show} onHide={closeMyBooks}>
      <Modal.Header closeButton>
        <Modal.Title>My Books</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <div>
              {myBooks && myBooks.length!==0 ? myBooks.map((book, index) => 
              book &&
                <div style={{display:'flex', justifyContent:"space-between", margin:'12px 4px', alignItems:'baseline'}} key={index}>
                  <p style={{fontWeight:'bold'}}>{book.Title}</p>
                  <Button variant="outline-primary" onClick={()=>onReturn(book.Id)}>Return</Button>
                </div>
              )
              :
              <div style={{padding:'12px'}}>
                Looks like you haven't borrowed any books yet!
              </div>
            }
        </div>
      </Modal.Body>
    </Modal>
  )
}

export default MyBooks