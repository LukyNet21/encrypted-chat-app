import { useState } from 'react'
import './App.css'
import { RegisterForm } from './components/Register'
import { Login } from './components/Login'
import { Chat } from './components/Chat'

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [ws, setWs] = useState<WebSocket | null>(null)

  const handleLoginSuccess = (websocket: WebSocket) => {
    setWs(websocket)
    setIsLoggedIn(true)
  }

  return (
    <div>
      <h1>Chat</h1>
      {!isLoggedIn ? (
        <>
          <Login onLoginSuccess={handleLoginSuccess} />
          <RegisterForm />
        </>
      ) : (
        ws ?
          <Chat ws={ws} />
          :
          <p>WS not loaded</p>
      )}
    </div>
  )
}

export default App
