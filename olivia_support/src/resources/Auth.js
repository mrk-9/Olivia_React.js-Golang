import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
  signIn(email, password) {
    return Vue.http.post('support/auth/signin', {email, password})
  },

  getProfile() {
    const token = VueCookie.get('support_token')

    Vue.http.headers.common['Authorization'] = `SupportBearer ${token}`;

    if (window.localStorage.getItem('support_user') && token) {
      return new Promise((resolve, reject) => {
        var response = {}
        response.body = JSON.parse(window.localStorage.getItem('support_user'))
        response.ok = true
        return resolve(response)
      })
    }

    return Vue.http.get('support/auth/getprofile').then(response => {
      window.localStorage.setItem('support_user', JSON.stringify(response.body))
      return response
    }, response => response)
  },

  logout($router) {
    VueCookie.delete('support_token');
    window.localStorage.removeItem('support_user');
    $router.replace('/auth/signin');
  },
}