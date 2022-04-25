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
        <Table hover>
          <tbody>
          <tr>
            <th>#</th>
            <th>Book Title</th>
            <th>Return Book?</th>
          </tr>
          </tbody>
          <tbody>
            {myBooks && myBooks.map((book, index) => 
            book &&
              <tr key={index}>
                <td>{index+1}</td>
                <td>{book.Title}</td>
                <td><Button variant="outline-primary" onClick={()=>onReturn(book.Id)}>Return</Button></td>
              </tr>
            )}
          </tbody>
        </Table>
      </Modal.Body>
    </Modal>
  )
}

export default MyBooks