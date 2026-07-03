import {ref, computed} from 'vue'
import {defineStore} from 'pinia'
import * as App from '../../wailsjs/go/main/App'
import type {model} from '../../wailsjs/go/models'
import {useWorkspaceStore} from './workspace'

export const usePluginsStore = defineStore('plugins', () => {
  const available = ref<model.PluginManifest[]>([])
  const targets = ref<model.Target[]>([])
  const presets = ref<model.Target[]>([])
  const styles = ref<model.CustomStyle[]>([])

  const enabledOptionalNames = computed(() => available.value.filter(p => !p.core && p.enabled).map(p => p.name))

  function isPluginEnabled(name: string) {
    return available.value.find(p => p.name === name)?.enabled ?? false
  }

  async function refreshAll() {
    ;[available.value, targets.value, presets.value, styles.value] = await Promise.all([
      App.ListAvailablePlugins(),
      App.ListTargets(),
      App.ListTargetPresets(),
      App.ListStyles(),
    ])
  }

  async function setPluginEnabled(name: string, enabled: boolean) {
    await App.SetPluginEnabled(name, enabled)
    available.value = await App.ListAvailablePlugins()
  }

  async function saveStyles(newStyles: model.CustomStyle[]) {
    await App.SaveStyles(newStyles)
    styles.value = await App.ListStyles()
  }

  async function saveTarget(t: model.Target) {
    await App.SaveTarget(t)
    targets.value = await App.ListTargets()
  }

  async function deleteTarget(id: string) {
    await App.DeleteTarget(id)
    targets.value = await App.ListTargets()
  }

  async function setActiveTarget(id: string) {
    await App.SetActiveTarget(id)
    await useWorkspaceStore().refreshChapters() // também recarrega project.json (activeTarget)
  }

  return {
    available,
    targets,
    presets,
    styles,
    enabledOptionalNames,
    isPluginEnabled,
    refreshAll,
    setPluginEnabled,
    saveStyles,
    saveTarget,
    deleteTarget,
    setActiveTarget,
  }
})
