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

if (window.location.host == 'olivia.yakovlevi.ch') {
  global.host = 'api.olivia.yakovlevi.ch'
} else if (window.location.host == 'getolivia.co' || window.location.host == 'www.getolivia.co') {
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
  {path: '/', component: view('Start')},
  {path: '/auth/signin', component: view('Signin')},
  {path: '/auth/signup', component: view('Signup')},
  {path: '/home', component: view('Home')},
  {path: '/scripts', component: view('Scripts')},
  {path: '/account', component: view('Account')},
  {path: '/plans', component: view('Plans')}
]

var router = new VueRouter({
  mode: 'history',
  routes
})

router.beforeEach((to, from, next) => {
  if (['/', '/auth/signin', '/auth/signup', '/plans'].indexOf(to.path) != -1) {
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
