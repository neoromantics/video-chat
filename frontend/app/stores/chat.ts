import { defineStore } from 'pinia'

interface Peer {
  id: string
  stream: MediaStream | null
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    userId: null as string | null,
    roomId: null as string | null,
    isJoined: false,
    localStream: null as MediaStream | null,
    peers: [] as Peer[],
  }),
  actions: {
    setUserId(id: string) {
      this.userId = id
    },
    setJoined(roomId: string) {
      this.roomId = roomId
      this.isJoined = true
    },
    resetRoom() {
      this.roomId = null
      this.isJoined = false
      this.peers = []
    },
    addPeer(peerId: string) {
      if (!this.peers.find(p => p.id === peerId)) {
        this.peers.push({ id: peerId, stream: null })
      }
    },
    removePeer(peerId: string) {
      this.peers = this.peers.filter(p => p.id !== peerId)
    },
    updatePeerStream(peerId: string, stream: MediaStream) {
      const peer = this.peers.find(p => p.id === peerId)
      if (peer) {
        peer.stream = stream
      }
    }
  }
})
