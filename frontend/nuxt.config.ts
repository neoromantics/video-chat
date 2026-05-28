// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  future: {
    compatibilityVersion: 4,
  },
  srcDir: 'app',
  devtools: { enabled: true },
  app: {
    head: {
      title: 'Video Chat',
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }
      ]
    }
  },
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      turnUrl: 'turn:vcm-51101.vm.duke.edu:3478',
      turnUsername: '', // Set via NUXT_PUBLIC_TURN_USERNAME
      turnPassword: '', // Set via NUXT_PUBLIC_TURN_PASSWORD
    }
  }
})
