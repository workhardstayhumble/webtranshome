<script setup>
import SectionHeading from '@/components/ui/SectionHeading.vue'
defineProps({ id: { type: String, required: true }, judul: { type: String, required: true }, tautan: { type: String, default: '' }, labelTautan: { type: String, default: '' } })
</script>
<template>
  <SectionHeading :id="id" :title="judul" :to="tautan" :link-label="labelTautan || 'Lihat semua'" />
</template>