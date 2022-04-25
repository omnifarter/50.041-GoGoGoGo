import React from 'react'
import { Button } from 'react-bootstrap'

import './Book.css';

const Book = ({ book, openBook, ...others }) => {
  return (
    <div className="card" style={{padding:'12px 16px', maxWidth:'600px', height:'240px', cursor:'pointer'}} onClick={() => openBook(book)}>
      <div className='Book-item'>
        <img className='Book-cover' src={book.Img_url} alt={`Book Cover of ${book.Title}`} />
        <div style={{marginLeft:'12px'}}>
          <p style={{color:'grey',marginBottom:'4px'}}>Title</p>
          <h6 className="Book-title">{book.Title}</h6>
        </div>
      </div>
    </div>
  )
}

export default Book