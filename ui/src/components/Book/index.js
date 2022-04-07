import React from 'react'
import { Button } from 'react-bootstrap'

import './Book.css';

const Book = ({ book, openBook }) => {
  return (
    <div className='Book-item'>
      <img className='Book-cover' src={book.Img_url} alt={`Book Cover of ${book.Title}`} />
      <h6 className="Book-title">{book.Title}</h6>
      <Button variant="primary" onClick={() => openBook(book)}>Borrow</Button>
    </div>
      
  )
}

export default Book