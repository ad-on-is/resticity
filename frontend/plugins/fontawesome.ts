import { defineNuxtPlugin } from 'nuxt/app'

import { library, config } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { fab } from '@fortawesome/free-brands-svg-icons'

library.add(fas, far, fab)
config.autoAddCss = false

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component('FaIcon', FontAwesomeIcon)
})
