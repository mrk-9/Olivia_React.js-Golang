import Vue from 'vue'
import VueRouter from 'vue-router'
import VueResource from 'vue-resource'
import VueCookie from 'vue-cookie'
import VueMoment from './modules/vendor/vue-moment'
import AuthResource from './resources/Auth'

import './assets/scss/index.scss'

global.color = '#39c3dd'

Vue.use(VueRouter)
Vue.use(VueResource)
Vue.use(VueCookie)
Vue.use(VueMoment)

if (window.location.host == 'support.olivia.yakovlevi.ch') {
  global.host = 'api.olivia.yakovlevi.ch'
} else if (
  window.location.host == 'getolivia.co' ||
  window.location.host == 'www.getolivia.co' ||
  window.location.host == 'support.getolivia.co' ||
  window.location.host == 'www.support.getolivia.co'
) {
  global.host = 'api.getolivia.co'
} else {
  global.host = 'localhost:8081'
}

Vue.http.options.root = `http://${global.host}`

var view = (name) => {
  return function (resolve) {
    require(['./routes/' + name + '.vue'], resolve);
  }
}

const routes = [
  {path: '/', component: view('Home')},
  {path: '/auth/signin', component: view('Signin')},
]

var router = new VueRouter({
  mode: 'history',
  routes
})

router.beforeEach((to, from, next) => {
  if (to.path == '/auth/signin') {
    next()
    return
  }

  AuthResource.getProfile().then(response => {
    if (response.ok) {
      next()
    } else {
      window.location = '/auth/signin'
    }
  }, response => {
    console.log(response)
    window.location = '/auth/signin'
  })
})

new Vue({
  router,
}).$mount('#app')