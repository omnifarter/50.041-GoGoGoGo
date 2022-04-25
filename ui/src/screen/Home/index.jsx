// hardcoded data for display
import { useEffect, useState } from "react"
import { Button } from "react-bootstrap"
import Book from "../../components/Book"
import Borrow from "../../components/Borrow"
import MyBooks from "../../components/MyBooks"
import { borrowBook, getAllBooks } from "../../helpers/APIs"
import Login from '../../components/Login'
  const emptyBook = {
    title: "",
    id: -1,
    image: "",
    available: true,
  }
  
function Home() {
    const [user, setUser] = useState(null)
    const [borrow, setBorrow] = useState(false)
    const [myBooks, setMyBooks] = useState(false)
    const [myBookList, setMyBookList] = useState(false)
    const [selectedBook, setSelectedBook] = useState(emptyBook)
    const [books, setBooks] = useState([])

    const fetchAllBooks = async () => {
        const booksFetched = await getAllBooks()
        setBooks(booksFetched.data)
        setMyBookList(booksFetched.data.filter(b=>{if(b.Borrowed && b.UserId==user) return b}))
    }
    // open the borrow modal
    const openBook = (book) => {
      setSelectedBook(book)
      setBorrow(true)
    }
  
    // close the borrow modal
    const closeBorrow = () => {
      setSelectedBook(emptyBook)
      setBorrow(false)
    }
  
    // open the myBooks modal
    const openMyBooks = () => {
      setMyBooks(true)
    }
  
    // close the myBooks modal
    const closeMyBooks = () => {
      setMyBooks(false)
    }
  
    // parameter: client id
    const confirmBorrow = async () => {
      await borrowBook(selectedBook.Id,user)
      setBorrow(false)
      fetchAllBooks()
    }

    const onReturn = (bookId) => {
      borrowBook(bookId, -1)
      fetchAllBooks()
      
    }
  
    useEffect(() => {
      user !== null && fetchAllBooks()
    }, [user])

  
    return (
      <div className="App">
        <Login show={user===null} onSetUser={(id)=>setUser(id)} />
        <header className="App-header">
          <div style={{display:"flex"}}>
            <img src="https://cdn-icons.flaticon.com/png/512/3389/premium/3389081.png?token=exp=1650885181~hmac=f24fcd79094309e45141ea043ab7ae48" className="App-Icon"/>
            <h1 className='Library-title'>GoGoGoGo - Digital Library</h1>
          </div>
          <Button variant="outline-primary" onClick={() => openMyBooks()}>View My Books</Button>{' '}
        </header>
  
        <div className="Books-library">
          <h4>All Available Books</h4>
          <div style={{display:'grid',gridTemplateColumns:'600px 600px',columnGap:'12px',rowGap:'12px'}}>  
            {books && books.map((b) => !b.Borrowed && <Book key={b.Id} book={b} openBook={openBook} />)}
          </div>
        </div>
        
        <Borrow show={borrow} book={selectedBook} closeBorrow={closeBorrow} confirmBorrow={confirmBorrow} />
        <MyBooks show={myBooks} closeMyBooks={closeMyBooks} myBooks={myBookList} onReturn={onReturn} />
      </div>
    );
  }

export default Home;