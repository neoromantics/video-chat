import { useChatStore } from '~/stores/chat'

export const useWebRTC = () => {
  const store = useChatStore()
  const config = useRuntimeConfig()
  let ws: WebSocket | null = null
  const peerConnections: Record<string, RTCPeerConnection> = {}

  const iceServers = [
    { urls: 'stun:stun.l.google.com:19302' },
    {
      urls: config.public.turnUrl,
      username: config.public.turnUsername,
      credential: config.public.turnPassword
    }
  ]

  const init = async () => {
    try {
      store.localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
      connectWS()
    } catch (err) {
      console.error('Media error', err)
    }
  }

  const connectWS = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = process.dev 
      ? `ws://${window.location.hostname}:8080/ws`
      : `${protocol}//${window.location.host}/ws`
    
    ws = new WebSocket(wsUrl)

    ws.onmessage = async (event) => {
      const msg = JSON.parse(event.data)
      switch (msg.type) {
        case 'welcome':
          store.setUserId(msg.payload.user_id)
          break
        case 'joined_room':
          store.setJoined(msg.payload.room_id)
          if (msg.payload.peers) {
            for (const peerId of msg.payload.peers) {
              if (peerId !== store.userId) {
                await initiateCall(peerId)
              }
            }
          }
          break
        case 'peer_joined':
          if (msg.payload.roomId === store.roomId && msg.payload.peerId !== store.userId) {
            store.addPeer(msg.payload.peerId)
            // Note: The joiner is the one who initiates the call. 
            // Existing peers just wait for the offer.
          }
          break
        case 'signal':
          handleSignal(msg)
          break
      }
    }

  }

  const initiateCall = async (peerId: string) => {
    store.addPeer(peerId)
    const pc = createPeerConnection(peerId)
    const offer = await pc.createOffer()
    await pc.setLocalDescription(offer)
    sendSignal(peerId, { type: 'offer', sdp: offer.sdp })
  }

  const createPeerConnection = (peerId: string) => {
    if (peerConnections[peerId]) peerConnections[peerId].close()
    
    const pc = new RTCPeerConnection({ iceServers })
    peerConnections[peerId] = pc

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        sendSignal(peerId, { type: 'candidate', candidate: event.candidate })
      }
    }

    pc.ontrack = (event) => {
      store.updatePeerStream(peerId, event.streams[0])
    }

    if (store.localStream) {
      store.localStream.getTracks().forEach(track => {
        pc.addTrack(track, store.localStream!)
      })
    }

    return pc
  }

  const handleSignal = async (msg: any) => {
    const { from, payload } = msg
    let pc = peerConnections[from]

    if (payload.type === 'offer') {
      pc = createPeerConnection(from)
      await pc.setRemoteDescription(new RTCSessionDescription(payload))
      const answer = await pc.createAnswer()
      await pc.setLocalDescription(answer)
      sendSignal(from, { type: 'answer', sdp: answer.sdp })
    } else if (payload.type === 'answer') {
      await pc.setRemoteDescription(new RTCSessionDescription(payload))
    } else if (payload.type === 'candidate') {
      if (pc) await pc.addIceCandidate(new RTCIceCandidate(payload.candidate))
    }
  }

  const sendSignal = (to: string, payload: any) => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'signal', to, payload }))
    }
  }

  const join = (roomId: string) => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'join', payload: { roomId } }))
    }
  }

  const cleanup = () => {
    if (ws) ws.close()
    Object.values(peerConnections).forEach(pc => pc.close())
    store.localStream?.getTracks().forEach(t => t.stop())
    store.resetRoom()
  }

  onBeforeUnmount(cleanup)

  return { init, join, cleanup }
}
