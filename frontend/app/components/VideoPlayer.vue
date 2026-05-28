<script setup lang="ts">
const props = defineProps<{
  stream: MediaStream | null
  muted?: boolean
  label?: string
  mirrored?: boolean
}>()

const videoRef = ref<HTMLVideoElement | null>(null)

watch(() => props.stream, (newStream) => {
  if (videoRef.value && newStream) {
    videoRef.value.srcObject = newStream
  }
}, { immediate: true })

onMounted(() => {
  if (videoRef.value && props.stream) {
    videoRef.value.srcObject = props.stream
  }
})
</script>

<template>
  <div class="relative w-full aspect-video bg-gray-100 rounded border border-gray-200 overflow-hidden">
    <video
      ref="videoRef"
      autoplay
      playsinline
      :muted="muted"
      class="w-full h-full object-cover"
      :class="{ 'scale-x-[-1]': mirrored }"
    ></video>
    <div class="absolute bottom-2 left-2 px-2 py-0.5 bg-white/80 border border-gray-200 text-[10px] uppercase tracking-wider text-gray-600 rounded">
      {{ label || 'User' }}
    </div>
  </div>
</template>
