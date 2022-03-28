import React from 'react'
import { Button } from 'react-bootstrap'

import './Book.css';

const Book = ({ book, openBook }) => {
  return (
    <div className='Book-item'>
      <img className='Book-cover' src={book.image} alt={`Book Cover of ${book.title}`} />
      <h6 className="Book-title">{book.title}</h6>
      <Button variant="primary" onClick={() => openBook(book)}>Borrow</Button>
    </div>
      
  )
}

export default Book