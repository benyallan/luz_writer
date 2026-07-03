<script setup lang="ts">
import {computed} from 'vue'
import {SelectRoot, SelectTrigger, SelectValue, SelectIcon, SelectPortal, SelectContent, SelectViewport, SelectItem, SelectItemText} from 'reka-ui'
import type {model} from '../../wailsjs/go/models'

// Renderiza qualquer FormSchema (seção 8.2) genericamente. Em modo normal,
// `modelValue` é a config completa (já com defaults aplicados pelo backend
// via Resolve). Em modo override (seção 8.4), cada campo tem um toggle
// "sobrescrever": desligado mostra o valor herdado esmaecido; ligado grava em
// `modelValue` (que aqui representa só as sobrescritas, um objeto parcial).
const props = withDefaults(
  defineProps<{
    schema: model.FormSchema
    modelValue: Record<string, any>
    overrideMode?: boolean
    inheritedValues?: Record<string, any>
    inheritedLabel?: string
  }>(),
  {overrideMode: false, inheritedValues: () => ({}), inheritedLabel: 'do target'},
)

const emit = defineEmits<{'update:modelValue': [value: Record<string, any>]}>()

const DIMENSION_UNITS = ['in', 'cm', 'mm', 'pt']

function fieldValue(key: string, def: any) {
  if (props.modelValue && key in props.modelValue) return props.modelValue[key]
  return def
}

function isOverridden(key: string) {
  return props.overrideMode && props.modelValue && key in props.modelValue
}

function setValue(key: string, value: any) {
  emit('update:modelValue', {...props.modelValue, [key]: value})
}

function toggleOverride(key: string, on: boolean, def: any) {
  const next = {...props.modelValue}
  if (on) {
    next[key] = key in next ? next[key] : (props.inheritedValues[key] ?? def)
  } else {
    delete next[key]
  }
  emit('update:modelValue', next)
}

function parseDimension(v: string): {amount: string; unit: string} {
  const m = /^([\d.]*)\s*(in|cm|mm|pt)?$/.exec((v ?? '').trim())
  return {amount: m?.[1] || '', unit: m?.[2] || 'cm'}
}

function dimensionParts(field: model.FormField) {
  return parseDimension(fieldValue(field.key, field.default))
}

function setDimension(field: model.FormField, amount: string, unit: string) {
  setValue(field.key, `${amount}${unit}`)
}
</script>

<template>
  <div class="schema-form">
    <div v-for="field in schema.fields" :key="field.key" class="schema-form__field">
      <label class="schema-form__label">
        <input
          v-if="overrideMode"
          type="checkbox"
          class="schema-form__override-toggle"
          :checked="isOverridden(field.key)"
          @change="toggleOverride(field.key, ($event.target as HTMLInputElement).checked, field.default)"
        />
        {{ field.label }}
        <span v-if="overrideMode && !isOverridden(field.key)" class="schema-form__inherited-tag">
          herdado {{ inheritedLabel }}
        </span>
      </label>

      <div class="schema-form__control" :class="{'schema-form__control--disabled': overrideMode && !isOverridden(field.key)}">
        <template v-if="field.type === 'switch'">
          <input
            type="checkbox"
            :disabled="overrideMode && !isOverridden(field.key)"
            :checked="overrideMode ? (fieldValue(field.key, inheritedValues[field.key] ?? field.default)) : fieldValue(field.key, field.default)"
            @change="setValue(field.key, ($event.target as HTMLInputElement).checked)"
          />
        </template>

        <template v-else-if="field.type === 'select'">
          <SelectRoot
            :model-value="overrideMode ? fieldValue(field.key, inheritedValues[field.key] ?? field.default) : fieldValue(field.key, field.default)"
            :disabled="overrideMode && !isOverridden(field.key)"
            @update:model-value="v => setValue(field.key, v)"
          >
            <SelectTrigger class="select-trigger">
              <SelectValue />
              <SelectIcon>▾</SelectIcon>
            </SelectTrigger>
            <SelectPortal>
              <SelectContent class="select-content">
                <SelectViewport>
                  <SelectItem v-for="opt in field.options" :key="opt.value" :value="opt.value" class="select-item">
                    <SelectItemText>{{ opt.label }}</SelectItemText>
                  </SelectItem>
                </SelectViewport>
              </SelectContent>
            </SelectPortal>
          </SelectRoot>
        </template>

        <template v-else-if="field.type === 'dimension'">
          <input
            type="number"
            step="0.01"
            class="schema-form__dimension-amount"
            :disabled="overrideMode && !isOverridden(field.key)"
            :value="dimensionParts(field).amount"
            @change="setDimension(field, ($event.target as HTMLInputElement).value, dimensionParts(field).unit)"
          />
          <select
            class="schema-form__dimension-unit"
            :disabled="overrideMode && !isOverridden(field.key)"
            :value="dimensionParts(field).unit"
            @change="setDimension(field, dimensionParts(field).amount, ($event.target as HTMLSelectElement).value)"
          >
            <option v-for="u in DIMENSION_UNITS" :key="u" :value="u">{{ u }}</option>
          </select>
        </template>

        <template v-else-if="field.type === 'number'">
          <input
            type="number"
            :disabled="overrideMode && !isOverridden(field.key)"
            :value="overrideMode ? fieldValue(field.key, inheritedValues[field.key] ?? field.default) : fieldValue(field.key, field.default)"
            @change="setValue(field.key, Number(($event.target as HTMLInputElement).value))"
          />
        </template>

        <template v-else>
          <input
            type="text"
            :disabled="overrideMode && !isOverridden(field.key)"
            :value="overrideMode ? fieldValue(field.key, inheritedValues[field.key] ?? field.default) : fieldValue(field.key, field.default)"
            @change="setValue(field.key, ($event.target as HTMLInputElement).value)"
          />
        </template>
      </div>
    </div>

    <p v-if="!schema.fields.length" class="schema-form__empty">Este módulo não tem configurações.</p>
  </div>
</template>

<style scoped>
.schema-form__field {
  margin-bottom: 12px;
}

.schema-form__label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.82rem;
  margin-bottom: 4px;
}

.schema-form__inherited-tag {
  font-size: 0.68rem;
  color: var(--luz-fg-muted);
  font-style: italic;
}

.schema-form__control {
  display: flex;
  align-items: center;
  gap: 6px;
}

.schema-form__control--disabled {
  opacity: 0.55;
}

.schema-form__control input[type='text'],
.schema-form__control input[type='number'] {
  padding: 5px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  flex: 1;
  min-width: 0;
}

.schema-form__dimension-amount {
  width: 90px;
  padding: 5px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
}

.schema-form__dimension-unit {
  padding: 5px 6px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: var(--luz-bg-editor);
  color: var(--luz-fg);
}

.schema-form__empty {
  color: var(--luz-fg-muted);
  font-size: 0.82rem;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
  min-width: 160px;
}

/* .select-content/.select-item vivem em style.css (global) — ver comentário
   lá: conteúdo teleportado pelo Reka UI não é alcançado por <style scoped>. */
</style>
