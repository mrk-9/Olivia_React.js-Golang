import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
  getScriptsAndAnswers(user, lead) {
    const token = VueCookie.get('support_token')
    Vue.http.headers.common['Authorization'] = `SupportBearer ${token}`
    return Vue.http.get('support/scriptsandanswers', {params: {user, lead}})
  },

  saveAnswer(scriptId, lead, answer) {
    const token = VueCookie.get('support_token')
    Vue.http.headers.common['Authorization'] = `SupportBearer ${token}`
    return Vue.http.post('support/scripts/answer', {scriptId, lead, answer})
  }
}