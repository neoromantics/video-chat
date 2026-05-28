<script setup lang="ts">
import { useChatStore } from '~/stores/chat'
import { useWebRTC } from '~/composables/useWebRTC'

const route = useRoute()
const store = useChatStore()
const { init, join, cleanup } = useWebRTC()

onMounted(async () => {
  await init()
  join(route.params.id as string)
})

onBeforeUnmount(() => {
  cleanup()
})

watch(() => store.error, (newError) => {
  if (newError) {
    alert(newError)
    navigateTo('/')
  }
})

const handleLeave = () => {
  navigateTo('/')
}
</script>

<template>
  <div class="min-h-screen bg-white text-gray-900 font-sans p-6 md:p-12">
    <div class="max-w-6xl mx-auto">
      <header class="mb-12 flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-gray-100 pb-8">
        <div>
          <NuxtLink to="/" class="text-[10px] uppercase tracking-widest text-gray-400 hover:text-gray-900 transition-colors mb-4 inline-block">
            ← Home
          </NuxtLink>
          <h1 class="text-3xl font-light tracking-tight text-gray-950 mb-1">
            Room {{ route.params.id }}
          </h1>
        </div>

        <div class="flex items-center gap-4">
          <div class="flex flex-col items-end">
            <span class="text-xs uppercase tracking-widest text-gray-400 font-semibold text-[9px]">Status</span>
            <span class="text-sm font-medium text-green-500" v-if="store.isJoined">Connected</span>
            <span class="text-sm font-medium text-gray-400" v-else>Connecting...</span>
          </div>
          <button
            @click="handleLeave"
            class="px-6 py-2 border border-red-100 text-red-500 text-sm font-medium rounded hover:bg-red-50 transition-colors"
          >
            Leave
          </button>
        </div>
      </header>

      <main>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Local Video -->
          <div class="col-span-1">
            <VideoPlayer :stream="store.localStream" label="You" mirrored muted />
          </div>

          <!-- Peer Videos -->
          <template v-for="peer in store.peers" :key="peer.id">
            <div class="col-span-1">
              <VideoPlayer :stream="peer.stream" :label="`Peer ${peer.id.slice(0,4)}`" />
            </div>
          </template>
        </div>

        <div v-if="store.peers.length === 0 && store.isJoined" class="mt-24 text-center">
          <p class="text-xs uppercase tracking-[0.3em] text-gray-300 font-medium">Waiting for others to join...</p>
        </div>
      </main>
    </div>
  </div>
</template>
