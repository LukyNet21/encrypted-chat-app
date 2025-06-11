import { useState } from 'react'
import * as openpgp from 'openpgp'

export function RegisterForm() {
  const [username, setUsername] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      const { privateKey, publicKey } = await openpgp.generateKey({
        type: 'ecc',
        curve: 'nistP521',
        userIDs: [{ name: username, email: `${username}@chat.local` }],
        format: 'armored'
      })

      const privateKeyBlob = new Blob([privateKey], { type: 'text/plain' })
      const privateKeyURL = URL.createObjectURL(privateKeyBlob)
      const downloadLink = document.createElement('a')
      downloadLink.href = privateKeyURL
      downloadLink.download = `${username}_private.key`
      document.body.appendChild(downloadLink)
      downloadLink.click()
      document.body.removeChild(downloadLink)
      URL.revokeObjectURL(privateKeyURL)

      // Register user with public key
      const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: username,
          public_key: publicKey,
        }),
      })

      if (!response.ok) {
        const errorData = await response.text()
        throw new Error(errorData || 'Registration failed')
      }

      alert('Registration successful! Your private key has been downloaded. Please keep it safe!')

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="form">
      <h2>Register</h2>
      {error && <div className="error">{error}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Username (min. 6 characters):</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            minLength={6}
            required
            disabled={isLoading}
          />
        </div>
        <button type="submit" disabled={isLoading || username.length < 6}>
          {isLoading ? 'Generating keys...' : 'Register'}
        </button>
      </form>
    </div>
  )
}