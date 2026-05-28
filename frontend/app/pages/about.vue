<script setup lang="ts">
// No logic needed
</script>

<template>
  <div class="min-h-screen bg-white text-gray-900 font-sans p-6 md:p-12">
    <div class="max-w-3xl mx-auto">
      <header class="mb-12 border-b border-gray-100 pb-8">
        <NuxtLink to="/" class="text-xs uppercase tracking-widest text-gray-400 hover:text-gray-900 transition-colors mb-4 inline-block">
          ← Back to Chat
        </NuxtLink>
        <h1 class="text-4xl font-light tracking-tight text-gray-950 mb-2">About</h1>
        <p class="text-lg text-gray-500 font-light">The story behind the video chat.</p>
      </header>

      <main class="prose prose-gray max-w-none">
        <section class="mb-12">
          <h2 class="text-xl font-medium text-gray-900 mb-4">Creator</h2>
          <p class="text-gray-600 leading-relaxed mb-6">
            This platform was designed and developed by <span class="text-gray-950 font-medium">Taiyan Liu</span>. 
            It was built with a focus on privacy, peer-to-peer performance, and a clean, undistracted user experience.
          </p>
          <a 
            href="https://github.com/neoromantics" 
            target="_blank" 
            class="inline-flex items-center gap-2 px-6 py-3 bg-gray-950 text-white text-sm font-medium rounded hover:bg-gray-800 transition-colors"
          >
            Follow on GitHub
          </a>
        </section>

        <section class="mb-12">
          <h2 class="text-xl font-medium text-gray-900 mb-4">Technology</h2>
          <p class="text-gray-600 leading-relaxed">
            The application leverages <span class="font-medium">Go</span> and <span class="font-medium">Redis</span> for scalable signaling, 
            and <span class="font-medium">Nuxt 3</span> for a responsive, modern frontend. All video communication is 
            encrypted and handled directly between peers using <span class="font-medium">WebRTC</span>.
          </p>
        </section>
      </main>
    </div>
  </div>
</template>
