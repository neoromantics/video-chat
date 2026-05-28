<script setup lang="ts">
const roomInput = ref('')

const handleJoin = () => {
  if (roomInput.value.trim()) {
    navigateTo('/' + roomInput.value.trim())
  }
}
</script>

<template>
  <div class="min-h-screen bg-white text-gray-900 font-sans p-6 md:p-12">
    <div class="max-w-6xl mx-auto">
      <header class="mb-24 flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-gray-100 pb-8">
        <div>
          <h1 class="text-3xl font-light tracking-tight text-gray-950 mb-1">
            Video Chat
          </h1>
          <p class="text-sm text-gray-400 font-normal">
            Minimal, secure, peer-to-peer.
          </p>
        </div>

        <div class="flex flex-col md:flex-row items-center gap-3">
          <input
            v-model="roomInput"
            type="text"
            placeholder="Enter Room ID"
            class="px-4 py-2 bg-gray-50 border border-gray-200 rounded text-sm focus:outline-none focus:border-gray-400 transition-colors w-full md:w-64"
            @keyup.enter="handleJoin"
          />
          <button
            @click="handleJoin"
            class="w-full md:w-auto px-8 py-2 bg-gray-950 text-white text-sm font-medium rounded hover:bg-gray-800 transition-colors"
          >
            Join Room
          </button>
        </div>
      </header>

      <main class="text-center py-20">
        <h2 class="text-5xl font-extralight tracking-tight text-gray-900 mb-6">Connect instantly.</h2>
        <p class="text-gray-400 max-w-lg mx-auto leading-relaxed">
          A high-performance video communication platform designed for privacy and speed. 
          No tracking, no accounts, just peer-to-peer video.
        </p>
      </main>

      <footer class="mt-48 pt-8 border-t border-gray-50 flex flex-col md:flex-row justify-between gap-6">
        <div class="flex flex-wrap gap-x-8 gap-y-4">
          <NuxtLink to="/about" class="text-[10px] uppercase tracking-[0.2em] text-gray-400 hover:text-gray-900 transition-colors">About</NuxtLink>
          <NuxtLink to="/privacy" class="text-[10px] uppercase tracking-[0.2em] text-gray-400 hover:text-gray-900 transition-colors">Privacy</NuxtLink>
          <NuxtLink to="/terms" class="text-[10px] uppercase tracking-[0.2em] text-gray-400 hover:text-gray-900 transition-colors">Terms</NuxtLink>
          <NuxtLink to="/cookies" class="text-[10px] uppercase tracking-[0.2em] text-gray-400 hover:text-gray-900 transition-colors">Cookies</NuxtLink>
        </div>
        <p class="text-[10px] uppercase tracking-[0.2em] text-gray-300">
          &copy; 2026 Taiyan Liu
        </p>
      </footer>
    </div>
  </div>
</template>
