import { Tabs, Tab, Button, Row, Col,Container } from 'react-bootstrap';
import { useEffect, useState } from "react";
import {getBook,getAllBooks,borrowBook,addNode,removeNode, getNodes, addBook} from '../../helpers/APIs'
import ReactJson from 'react-json-view'
import NodeDisplay from "../../components/NodeDisplay";
import GetBook from '../../components/Testing/GetBook';
import BorrowBook from '../../components/Testing/BorrowBook';

function Test() {
    const [newKeyStructure, setNewKeyStructure] = useState({});
    const [oldKeyStructure, setOldKeyStructure] = useState({});
    const [response, setResponse] = useState()
    const [tab, setTab] = useState('get')

    const getAllBooksTest = async () => {
        const data = await getAllBooks()
        setResponse(data)
    }

    const addNodeTest = async () => {
        const data = await addNode()
        setOldKeyStructure(newKeyStructure)
        setNewKeyStructure(data.data)
    };

    const removeNodeTest = async () => {
        const data = await removeNode()
        setOldKeyStructure(newKeyStructure)
        setNewKeyStructure(data.data)
    };

    const getNodesTest = async () => {
        const data = await getNodes()
        setOldKeyStructure(newKeyStructure)
        setNewKeyStructure(data.data)
    };

    

    useEffect(() => {
        getNodesTest()
    }, [])

    return (
        <div>
            <header className="App-header">
                <h1 className="Library-title">GoGoGoGo - Test Page</h1>
            </header>
            <br />
            <div  style={{display:'flex',width:'100%', justifyContent:'space-evenly'}}>
                <div>
                    <h3>Old key structure</h3>
                    <NodeDisplay keyStructure={oldKeyStructure} />
                </div>
                <div>
                    <h3>New key structure</h3>
                    <NodeDisplay keyStructure={newKeyStructure} />
                </div>
            </div>
            <div style={{textAlign:'center',margin:'24px 0'}}>

                <Button variant="success" style={{marginRight:'12px'}} onClick={() => addNodeTest()}>
                    Add Node
                </Button>{" "}
                <Button variant="danger" onClick={() => removeNodeTest()}>
                    Remove Node
                </Button>{" "}
            </div>
            <Container>
                <Row>
                    <Col style={{width:'50%'}}>
                        <Tabs activeKey={tab} onSelect={(k) => setTab(k)}>
                            <Tab eventKey="get" title="Get Book">
                                <GetBook setResponse={setResponse} />
                            </Tab>
                            <Tab eventKey='borrow' title="Borrow Book">
                            <BorrowBook setResponse={setResponse} />
                            </Tab>
                            <Tab eventKey='getAll' title="Get All Books">
                                <h3 className="Library-title">Get All Books</h3>
                                <br />
                                <Button onClick={getAllBooksTest} style={{width:'100%'}}>Get all books</Button>
                            </Tab>
                        </Tabs>
                    </Col>
                    <Col style={{width:'50%'}}>
                        <ReactJson src={response} />
                    </Col>
                </Row>
            </Container>

        </div>
    );
}

export default Test;
