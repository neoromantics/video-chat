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
      console.log('[WebRTC] Initializing media...')
      store.localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
      await connectWS()
    } catch (err) {
      console.error('[WebRTC] Media error:', err)
    }
  }

  const connectWS = () => {
    return new Promise<void>((resolve, reject) => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsUrl = process.dev 
        ? `ws://${window.location.hostname}:8080/ws`
        : `${protocol}//${window.location.host}/ws`
      
      console.log('[WebRTC] Connecting to WebSocket:', wsUrl)
      ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        console.log('[WebRTC] WebSocket connected')
        resolve()
      }
      
      ws.onerror = (err) => {
        console.error('[WebRTC] WebSocket error:', err)
        reject(err)
      }

      ws.onmessage = async (event) => {
        const msg = JSON.parse(event.data)
        console.log('[WebRTC] Received message:', msg.type, msg.from || '')
        
        switch (msg.type) {
          case 'welcome':
            console.log('[WebRTC] Your UserID:', msg.payload.user_id)
            store.setUserId(msg.payload.user_id)
            break
          case 'joined_room':
            console.log('[WebRTC] Joined room:', msg.payload.room_id, 'with peers:', msg.payload.peers)
            store.setJoined(msg.payload.room_id)
            if (msg.payload.peers && msg.payload.peers.length > 0) {
              // Wait 500ms to ensure others are ready for signaling
              setTimeout(async () => {
                for (const peerId of msg.payload.peers) {
                  if (peerId !== store.userId) {
                    await initiateCall(peerId)
                  }
                }
              }, 500)
            }
            break
          case 'peer_joined':
            if (msg.payload.roomId === store.roomId && msg.payload.peerId !== store.userId) {
              console.log('[WebRTC] New peer joined:', msg.payload.peerId)
              store.addPeer(msg.payload.peerId)
            }
            break
          case 'signal':
            handleSignal(msg)
            break
        }
      }
    })
  }

  const initiateCall = async (peerId: string) => {
    console.log('[WebRTC] Initiating call to:', peerId)
    console.log('[WebRTC] Using ICE Servers:', iceServers)
    store.addPeer(peerId)
    const pc = createPeerConnection(peerId)
    const offer = await pc.createOffer()
    await pc.setLocalDescription(offer)
    sendSignal(peerId, offer)
  }

  const createPeerConnection = (peerId: string) => {
    console.log('[WebRTC] Creating PeerConnection for:', peerId)
    if (peerConnections[peerId]) peerConnections[peerId].close()
    
    const pc = new RTCPeerConnection({ iceServers })
    peerConnections[peerId] = pc

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        sendSignal(peerId, { type: 'candidate', candidate: event.candidate })
      }
    }

    pc.oniceconnectionstatechange = () => {
      console.log(`[WebRTC] ICE state (${peerId}):`, pc.iceConnectionState)
    }

    pc.ontrack = (event) => {
      console.log('[WebRTC] Received remote track from:', peerId)
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
      console.log('[WebRTC] Handling offer from:', from)
      pc = createPeerConnection(from)
      await pc.setRemoteDescription(new RTCSessionDescription(payload))
      const answer = await pc.createAnswer()
      await pc.setLocalDescription(answer)
      sendSignal(from, answer)
    } else if (payload.type === 'answer') {
      console.log('[WebRTC] Handling answer from:', from)
      if (pc) await pc.setRemoteDescription(new RTCSessionDescription(payload))
    } else if (payload.type === 'candidate') {
      console.log('[WebRTC] Handling candidate from:', from)
      if (pc && payload.candidate) {
        try {
          await pc.addIceCandidate(new RTCIceCandidate(payload.candidate))
        } catch (e) {
          console.error('[WebRTC] Error adding candidate:', e)
        }
      }
    }
  }

  const sendSignal = (to: string, payload: any) => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'signal', to, payload }))
    }
  }

  const join = (roomId: string) => {
    if (ws?.readyState === WebSocket.OPEN) {
      const id = roomId || 'lobby'
      console.log('[WebRTC] Joining room:', id)
      ws.send(JSON.stringify({ type: 'join', payload: { roomId: id } }))
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
