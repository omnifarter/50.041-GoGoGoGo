// hardcoded data for display

import { useEffect, useState } from "react"
import { Button } from "react-bootstrap"
import Book from "../../components/Book"
import Borrow from "../../components/Borrow"
import MyBooks from "../../components/MyBooks"
import { borrowBook, getAllBooks } from "../../helpers/APIs"

// TODO: remove once backend is connected
// const books = [{
//     title: "Animal Farm",
//     id: 0,
//     image: "https://d1w7fb2mkkr3kw.cloudfront.net/assets/images/book/lrg/9780/1410/9780141036137.jpg",
//     available: true,
//   }, {
//     title: "1984",
//     id: 1,
//     image: "https://kbimages1-a.akamaihd.net/5d088fbe-c36c-4a03-9317-0755143820c7/353/569/90/False/iFcx1981QTSb5BLWINslVA.jpg",
//     available: true,
//   }, {
//     title: "Macbeth",
//     id: 3,
//     image: "https://m.media-amazon.com/images/I/411rwBu7c4L.jpg",
//     available: true,
//   }]
  
  const emptyBook = {
    title: "",
    id: -1,
    image: "",
    available: true,
  }
  
function Home() {
    const [borrow, setBorrow] = useState(false)
    const [myBooks, setMyBooks] = useState(false)
    const [selectedBook, setSelectedBook] = useState(emptyBook)
    const [books, setBooks] = useState([])

    const fetchAllBooks = async () => {
        const booksFetched = await getAllBooks()
        setBooks(booksFetched)
    }
    // open the borrow modal
    const openBook = (book) => {
      console.log('open modal...')
      setSelectedBook(book)
      setBorrow(true)
    }
  
    // close the borrow modal
    const closeBorrow = () => {
      console.log('close modal...')
      setSelectedBook(emptyBook)
      setBorrow(false)
    }
  
    // open the myBooks modal
    const openMyBooks = () => {
      console.log('open my books modal...')
      setMyBooks(true)
    }
  
    // close the myBooks modal
    const closeMyBooks = () => {
      setMyBooks(false)
    }
  
    // parameter: client id
    const confirmBorrow = async (id) => {
      await borrowBook(selectedBook.id,id)
      setBorrow(false)
    }
  
    useEffect(() => {
      fetchAllBooks()
    }, [])
    useEffect(() => {
      fetchAllBooks()
    }, [borrow])
  
    return (
      <div className="App">
        <header className="App-header">
          <h1 className='Library-title'>GoGoGoGo - Digital Library</h1>
          <Button variant="info" onClick={() => openMyBooks()}>View My Books</Button>
        </header>
  
        <div className="Books-library">
          <h4>View All Available Books</h4>
          <div className='All-books'>  
            {books.map((b) => b.available && <Book key={b.id} book={b} openBook={openBook} />)}
          </div>
        </div>
        
        <Borrow show={borrow} book={selectedBook} closeBorrow={closeBorrow} confirmBorrow={confirmBorrow} />
        <MyBooks show={myBooks} closeMyBooks={closeMyBooks} myBooks={books} />
      </div>
    );
  }

export default Home;