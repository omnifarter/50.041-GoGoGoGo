import React, { useEffect, useState } from 'react'

import './App.css';
import Book from './components/Book';
import Borrow from './components/Borrow';

import { Button } from 'react-bootstrap'
import MyBooks from './components/MyBooks';
import Router from './router';
import { BrowserRouter } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <Router />
    </BrowserRouter>
  );
}

export default App;
