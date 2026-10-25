import { ref, watch } from 'vue'

const STORAGE_KEY = 'solar_os_settings_v1'

const defaultSettings = {
  theme: 'dark', // 'dark' | 'light'
  electricityPriceCent: 32.0, // ct/kWh Netzbezugspreis
  feedInTariffCent: 5.9, // ct/kWh Einspeisevergütung
  systemCapacityKwp: 5.0, // kWp Anlagenleistung
  batteryCapacityKwh: 6.0, // kWh Speicherkapazität
  co2GramPerKwh: 400, // g CO2 pro kWh Strommix
}

function loadSettings() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return { ...defaultSettings }
    return { ...defaultSettings, ...JSON.parse(raw) }
  } catch {
    return { ...defaultSettings }
  }
}

const settings = ref(loadSettings())
const activeSection = ref('overview') // 'overview' | 'flow' | 'forecast' | 'analytics'
const isOnline = ref(typeof navigator !== 'undefined' ? navigator.onLine : true)

if (typeof window !== 'undefined') {
  window.addEventListener('online', () => { isOnline.value = true })
  window.addEventListener('offline', () => { isOnline.value = false })
}

function applyThemeToDom(theme) {
  if (typeof document === 'undefined') return
  const html = document.documentElement
  if (theme === 'light') {
    html.classList.remove('dark')
    html.classList.add('light')
    const meta = document.querySelector('meta[name="theme-color"]')
    if (meta) meta.setAttribute('content', '#f8fafc')
  } else {
    html.classList.remove('light')
    html.classList.add('dark')
    const meta = document.querySelector('meta[name="theme-color"]')
    if (meta) meta.setAttribute('content', '#090a0f')
  }
}

applyThemeToDom(settings.value.theme)

watch(
  settings,
  (val) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(val))
    } catch {
      // ignore storage errors
    }
    applyThemeToDom(val.theme)
  },
  { deep: true }
)

export function useSettings() {
  const toggleTheme = () => {
    settings.value.theme = settings.value.theme === 'dark' ? 'light' : 'dark'
  }

  const resetSettings = () => {
    settings.value = { ...defaultSettings }
  }

  return {
    settings,
    activeSection,
    isOnline,
    toggleTheme,
    resetSettings,
  }
}

/**
 * Formats instantaneous power in W or kW
 */
export function formatPower(watts, { showUnit = true, compact = false } = {}) {
  if (watts === undefined || watts === null || isNaN(watts)) {
    return showUnit ? '0 W' : '0'
  }
  const absW = Math.abs(Number(watts))
  if (compact && absW >= 1000) {
    const kw = (absW / 1000).toLocaleString('de-DE', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })
    return showUnit ? `${kw} kW` : kw
  }
  const formatted = Math.round(absW).toLocaleString('de-DE')
  return showUnit ? `${formatted} W` : formatted
}

/**
 * Formats integrated energy in Wh or kWh
 */
export function formatEnergy(wh, { forceKwh = false } = {}) {
  if (wh === undefined || wh === null || isNaN(wh)) return '0 Wh'
  const absWh = Math.abs(Number(wh))
  if (forceKwh || absWh >= 1000) {
    const kwh = absWh / 1000
    const digits = kwh >= 100 ? 0 : 1
    return `${kwh.toLocaleString('de-DE', {
      minimumFractionDigits: digits,
      maximumFractionDigits: digits,
    })} kWh`
  }
  return `${Math.round(absWh).toLocaleString('de-DE')} Wh`
}

/**
 * Computes energy sums (in Wh) and KPI metrics from a list of Firebase `dailySum` entries.
 */
export function computePeriodMetrics(dailyList, currentSettings = settings.value) {
  const sum = {
    pvPower: 0, // Wh
    powerConsumption: 0, // Wh
    gridPowerIn: 0, // Wh (Netzeinspeisung / exported to grid)
    gridPowerOut: 0, // Wh (Netzbezug / imported from grid, stored negative or positive)
  }

  if (Array.isArray(dailyList)) {
    dailyList.forEach((day) => {
      if (!day || !day.count) return
      const measureTime = day.count / 120 // hours (30s interval = 120 samples/hr)
      sum.pvPower += ((day.pvPower || 0) / day.count) * measureTime
      sum.powerConsumption += ((day.powerConsumption || 0) / day.count) * measureTime
      sum.gridPowerIn += (Math.abs(day.gridPowerIn || 0) / day.count) * measureTime
      sum.gridPowerOut += (Math.abs(day.gridPowerOut || 0) / day.count) * measureTime
    })
  }

  const pvWh = Math.max(0, sum.pvPower)
  const consumptionWh = Math.max(0, sum.powerConsumption)
  const gridExportWh = Math.max(0, sum.gridPowerIn)
  const gridImportWh = Math.max(0, sum.gridPowerOut)

  // Direct/battery self-used solar energy
  const selfUsedSolarWh = Math.max(0, Math.min(pvWh - gridExportWh, consumptionWh))

  // Autarkiegrad: % of home consumption covered without grid import
  const autarkiePercent =
    consumptionWh > 0
      ? Math.min(100, Math.max(0, Math.round(((consumptionWh - gridImportWh) / consumptionWh) * 100)))
      : pvWh > 0
        ? 100
        : 0

  // Eigenverbrauchsquote: % of produced solar energy consumed locally (home + battery)
  const eigenverbrauchPercent =
    pvWh > 0
      ? Math.min(100, Math.max(0, Math.round(((pvWh - gridExportWh) / pvWh) * 100)))
      : 0

  // Financial calculation
  const pricePerKwhEuro = (currentSettings?.electricityPriceCent ?? 32) / 100
  const feedInPerKwhEuro = (currentSettings?.feedInTariffCent ?? 8.2) / 100
  const savedImportKwh = Math.max(0, consumptionWh - gridImportWh) / 1000
  const exportedKwh = gridExportWh / 1000

  const savingsEuro = savedImportKwh * pricePerKwhEuro + exportedKwh * feedInPerKwhEuro
  const co2SavedKg = ((pvWh / 1000) * (currentSettings?.co2GramPerKwh ?? 400)) / 1000

  return {
    pvWh,
    consumptionWh,
    gridExportWh,
    gridImportWh,
    selfUsedSolarWh,
    autarkiePercent,
    eigenverbrauchPercent,
    savingsEuro,
    co2SavedKg,
  }
}
