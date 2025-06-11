import { useState, useRef } from 'react'
import { createCleartextMessage, readPrivateKey, sign } from 'openpgp'

interface LoginProps {
  onLoginSuccess: (ws: WebSocket) => void;
}

export function Login({ onLoginSuccess }: LoginProps) {
  const [username, setUsername] = useState('')
  const [privateKey, setPrivateKey] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileRead = (e: ProgressEvent<FileReader>) => {
    const content = e.target?.result as string
    setPrivateKey(content)
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = handleFileRead
      reader.readAsText(file)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      if (!privateKey) {
        throw new Error('Please select your private key file')
      }

      const ws = new WebSocket('ws://localhost:8080/ws')

      ws.onopen = () => {
        ws.send(username)
      }

      ws.onmessage = async (event) => {
        const message = event.data
        
        if (message === 'user not found') {
          ws.close()
          setError('User not found')
          setIsLoading(false)
          return
        }

        if (message === 'invalid signature') {
          ws.close()
          setError('Invalid signature - wrong private key')
          setIsLoading(false)
          return
        }
        
        if (message === 'successfully signed in') {
          onLoginSuccess(ws)
          return
        }

        try {
          const privateKeyObj = await readPrivateKey({ armoredKey: privateKey })
          const signed = await createCleartextMessage({ text: message })
          const signature = await sign({
            message: signed,
            signingKeys: privateKeyObj
          })
          ws.send(signature)
        } catch (err) {
          ws.close()
          setError('Failed to sign challenge - invalid private key format')
          setIsLoading(false)
          return
        }
      }

      ws.onerror = () => {
        setError('WebSocket connection failed')
        setIsLoading(false)
      }

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
      setIsLoading(false)
    }
  }

  return (
    <div className="form">
      <h2>Login</h2>
      {error && <div className="error">{error}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Username:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            disabled={isLoading}
          />
        </div>
        <div className="form-group">
          <label htmlFor="privateKey">Private Key File:</label>
          <input
            type="file"
            id="privateKey"
            ref={fileInputRef}
            onChange={handleFileChange}
            accept=".key"
            required
            disabled={isLoading}
          />
        </div>
        <button type="submit" disabled={isLoading || !username || !privateKey}>
          {isLoading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  )
}