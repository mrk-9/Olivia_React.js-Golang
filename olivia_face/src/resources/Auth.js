import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
  signIn(email, password) {
    VueCookie.delete('token');
    window.localStorage.removeItem('user');

    return Vue.http.post('auth/signin', {email, password})
  },

  signUp(fname, sname, email, phone, password) {
    VueCookie.delete('token');
    window.localStorage.removeItem('user');

    return Vue.http.post('auth/signup', {fname, sname, email, code: phone.code, number: phone.number, password})
  },

  validatePhoneNumber(phone_number) {
    return Vue.http.post('auth/validatephone', {code: phone_number.code, number: phone_number.number})
  },

  getProfile() {
    const token = VueCookie.get('token')

    Vue.http.headers.common['Authorization'] = `Bearer ${token}`;

    if (window.localStorage.getItem('user') && token) {
      return new Promise((resolve, reject) => {
        var response = {}
        response.body = JSON.parse(window.localStorage.getItem('user'))
        response.ok = true
        return resolve(response)
      })
    }

    return Vue.http.get('auth/getprofile').then(response => {
      window.localStorage.setItem('user', JSON.stringify(response.body))
      return response
    }, response => response)
  },

  logout($router) {
    VueCookie.delete('token');
    window.localStorage.removeItem('user');
    $router.replace('/');
  },
}