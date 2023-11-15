/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'
import { VueFire } from 'vuefire'
import { firebaseApp } from './firebase'

// Plugins
import { registerPlugins } from '@/plugins'

const app = createApp(App)


app
  .use(VueFire, {
    // imported above but could also just be created here
    firebaseApp,
  })

registerPlugins(app)

app.mount('#app')
