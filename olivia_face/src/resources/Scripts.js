import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
  getScripts() {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.get('scripts/get')
  },

  createScript(title, text) {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.post('scripts/create', {title, text})
  },

  removeScript(id) {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.post('scripts/remove', {id})
  },

  updateScript(id, title, text) {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.post('scripts/update', {id, title, text})
  },

  getScriptsAnswers(lead) {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.get('scripts/answers', {params: {lead}})
  },
}