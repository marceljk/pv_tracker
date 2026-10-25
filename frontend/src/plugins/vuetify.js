/**
 * plugins/vuetify.js
 */

import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { createVuetify } from 'vuetify'

export default createVuetify({
  theme: {
    defaultTheme: 'dark',
    themes: {
      dark: {
        dark: true,
        colors: {
          background: '#090a0f',
          surface: '#11131c',
          'surface-bright': '#1a1d2b',
          primary: '#10b981',
          secondary: '#94a3b8',
          emerald: '#10b981',
          blue: '#3b82f6',
          amber: '#f59e0b',
          cyan: '#06b6d4',
          rose: '#f43f5e',
          error: '#ef4444',
          info: '#3b82f6',
          success: '#10b981',
          warning: '#f59e0b',
        },
      },
      light: {
        dark: false,
        colors: {
          background: '#f4f6fb',
          surface: '#ffffff',
          'surface-bright': '#f8fafc',
          primary: '#059669',
          secondary: '#64748b',
          emerald: '#059669',
          blue: '#2563eb',
          amber: '#d97706',
          cyan: '#0891b2',
          rose: '#e11d48',
          error: '#dc2626',
          info: '#2563eb',
          success: '#059669',
          warning: '#d97706',
        },
      },
    },
  },
})
