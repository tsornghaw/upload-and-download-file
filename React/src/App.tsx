import React, {useEffect, useState} from 'react';
import './App.css';
import Login from "./pages/Login";
import Nav from "./components/Nav";
import {BrowserRouter, Route, Redirect} from "react-router-dom";
import Home from "./pages/Home";
import Admin from "./pages/Admin";
import Register from "./pages/Register";

function App() {
    const [name, setName] = useState('');

    useEffect(() => {
        (
            async () => {
                const response = await fetch('http://localhost:8000/api/auth/user', {
                    headers: {'Content-Type': 'application/json'},
                    credentials: 'include',
                });
                console.log(name)

                const content = await response.json();

                setName(content.name);
            }
        )();
    }, []);


    return (
        <div className="App">
            <BrowserRouter>
                <Nav name={name} setName={setName}/>

                {/* <main className="form-signin">
                    <Route path="/" exact component={() => <Home name={name}/>}/>
                    <Route path="/login" component={() => <Login setName={setName}/>}/>
                    <Route path="/register" component={Register}/>
                </main> */}
                <main className="form-signin">
                    <Route path="/" exact component={() => <Home name={name} />} />
                    <Route path="/login" component={() => <Login setName={setName} />} />
                    <Route path="/register" component={Register} />

                    {/* Conditional Redirect */}
                    {name === "admin" ? (
                        <Redirect from="/" to="/admin" />
                    ) : (
                        <Redirect from="/" to="/home" />
                    )}

                    <Route path="/admin" component={() => <Admin name={name} />} />
                    {/* <Route path="/home" component={Home} /> */}
                </main>
            </BrowserRouter>
        </div>
    );
}

export default App;
