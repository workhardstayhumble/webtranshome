<script setup>
defineProps({
  columns: { type: Array, required: true },
  rows: { type: Array, default: () => [] },
  rowKey: { type: String, default: 'id' },
  caption: { type: String, default: '' },
  bare: { type: Boolean, default: false },
  minWidth: { type: String, default: 'min-w-160' },
})
const RATA = { left: 'text-left', right: 'text-right', center: 'text-center' }
</script>
<template>
  <div class="relative overflow-x-auto" :class="!bare && 'overflow-hidden rounded-xl border border-border bg-surface shadow-[0_2px_15px_rgba(28,36,48,0.03)]'">
    <table class="w-full border-collapse text-sm" :class="minWidth">
      <caption v-if="caption" class="sr-only">{{ caption }}</caption>
      <thead class="bg-subtle/75 text-xs font-bold tracking-wide text-muted uppercase">
        <tr>
          <th v-for="kolom in columns" :key="kolom.key" scope="col" class="border-b border-border px-4 py-3.5" :class="[RATA[kolom.align ?? 'left'], kolom.class]">
            <span :class="kolom.srOnly && 'sr-only'">{{ kolom.label }}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="baris in rows" :key="baris[rowKey]" class="border-b border-border/70 last:border-0 transition-colors duration-150 hover:bg-primary-soft/35">
          <td v-for="kolom in columns" :key="kolom.key" class="px-4 py-3.5" :class="[RATA[kolom.align ?? 'left'], kolom.class]">
            <slot :name="`cell-${kolom.key}`" :row="baris" :value="baris[kolom.key]">{{ baris[kolom.key] }}</slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>