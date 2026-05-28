<template>
  <div class="min-h-screen bg-white text-gray-900 font-sans p-6 md:p-12">
    <div class="max-w-3xl mx-auto">
      <header class="mb-12 border-b border-gray-100 pb-8">
        <NuxtLink to="/" class="text-xs uppercase tracking-widest text-gray-400 hover:text-gray-900 transition-colors mb-4 inline-block">
          ← Back to Chat
        </NuxtLink>
        <h1 class="text-4xl font-light tracking-tight text-gray-950 mb-2">Privacy Policy</h1>
        <p class="text-lg text-gray-500 font-light">Your security is our priority.</p>
      </header>

      <main class="prose prose-gray max-w-none text-gray-600 space-y-8">
        <section>
          <h2 class="text-xl font-medium text-gray-900 mb-4">1. Peer-to-Peer Privacy</h2>
          <p>
            All video and audio communications are handled directly between users using WebRTC. 
            We do not route media through our servers, and we do not record or store any of your conversations.
          </p>
        </section>

        <section>
          <h2 class="text-xl font-medium text-gray-900 mb-4">2. Data Collection</h2>
          <p>
            We do not require accounts or personal information. Temporary session IDs are used only 
            for matchmaking and are discarded immediately after you disconnect.
          </p>
        </section>

        <section>
          <h2 class="text-xl font-medium text-gray-900 mb-4">3. Security</h2>
          <p>
            We use industry-standard encryption for signaling and enforce secure connections (HTTPS) 
            to protect your session data from interception.
          </p>
        </section>
      </main>
    </div>
  </div>
</template>
