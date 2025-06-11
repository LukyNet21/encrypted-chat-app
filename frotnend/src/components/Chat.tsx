import { useEffect } from "react"

interface ChatProps {
  ws: WebSocket
}

export function Chat({ ws }: ChatProps) {
  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      console.log('Received message:', event.data)
    }

    ws.addEventListener('message', handleMessage)

    return () => {
      ws.removeEventListener('message', handleMessage)
    }
  }, [ws])
  
  return (
    <div>Chat interface will go here</div>
  )
}