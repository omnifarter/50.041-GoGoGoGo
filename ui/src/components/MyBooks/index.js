import React from 'react'
import { Modal, Button, Table } from 'react-bootstrap';

// import './MyBooks.css';

const MyBooks = ({show, closeMyBooks, myBooks}) => {
  return (
    <Modal show={show} onHide={closeMyBooks}>
      <Modal.Header closeButton>
        <Modal.Title>My Books</Modal.Title>
      </Modal.Header>

      <Modal.Body>
        <Table hover>
          <tr>
            <th>#</th>
            <th>Book Title</th>
            <th>Return Book?</th>
          </tr>
          <tbody>
            {myBooks.map((book, index) => 
              <tr key={index}>
                <td>{index+1}</td>
                <td>{book.title}</td>
                <td><Button variant="outline-primary">Return</Button></td>
              </tr>
            )}
          </tbody>
        </Table>
      </Modal.Body>
    </Modal>
  )
}

export default MyBooks