<script setup lang="ts">
import {computed, ref, watch} from 'vue'
import {nodeViewProps, NodeViewWrapper} from '@tiptap/vue-3'
import {PopoverRoot, PopoverTrigger, PopoverPortal, PopoverContent} from 'reka-ui'

const props = defineProps(nodeViewProps)

const caption = ref<string>(props.node.attrs.caption ?? '')
const width = ref<number>(props.node.attrs.width ?? 80)

watch(
  () => props.node.attrs,
  attrs => {
    caption.value = attrs.caption ?? ''
    width.value = attrs.width ?? 80
  },
)

const src = computed(() => props.node.attrs.src as string | null)
const resolvedSrc = computed(() => (src.value ? `/workspace-files/${src.value}` : ''))

function commitCaption() {
  props.updateAttributes({caption: caption.value})
}

function commitWidth() {
  const n = Math.min(100, Math.max(10, Number(width.value) || 80))
  width.value = n
  props.updateAttributes({width: n})
}
</script>

<template>
  <NodeViewWrapper class="luz-image" :style="{width: `${node.attrs.width ?? 80}%`}">
    <PopoverRoot>
      <PopoverTrigger as-child>
        <img v-if="resolvedSrc" :src="resolvedSrc" class="luz-image__img" />
        <div v-else class="luz-image__missing">Imagem não encontrada</div>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent class="luz-image__popover" :side-offset="4">
          <label class="luz-image__field">
            Legenda
            <input type="text" v-model="caption" placeholder="(sem legenda)" @blur="commitCaption" @keyup.enter="commitCaption" />
          </label>
          <label class="luz-image__field">
            Largura (% do texto)
            <input type="number" min="10" max="100" v-model.number="width" @blur="commitWidth" @keyup.enter="commitWidth" />
          </label>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
    <p v-if="node.attrs.caption" class="luz-image__caption">{{ node.attrs.caption }}</p>
  </NodeViewWrapper>
</template>

<style scoped>
.luz-image {
  margin: 1em auto;
  display: block;
}

.luz-image__img {
  width: 100%;
  display: block;
  border: 1px solid var(--luz-border);
  border-radius: 4px;
  cursor: pointer;
}

.luz-image__missing {
  border: 1px dashed var(--luz-border);
  border-radius: 4px;
  padding: 24px;
  text-align: center;
  color: var(--luz-fg-muted);
  font-size: 0.8rem;
}

.luz-image__caption {
  margin: 6px 0 0;
  font-size: 0.8rem;
  color: var(--luz-fg-muted);
  text-align: center;
}

.luz-image__popover {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 6px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.8rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 220px;
}

.luz-image__field {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.luz-image__field input {
  padding: 4px 6px;
  border-radius: 4px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
}
</style>
