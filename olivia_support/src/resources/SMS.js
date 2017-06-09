import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
    sendSMS(user, lead, message) {
        const token = VueCookie.get('support_token')
        Vue.http.headers.common['Authorization'] = `SupportBearer ${token}`
        return Vue.http.post('support/sms/send', {user, lead, message})
    },
    getChat(user, lead) {
        const token = VueCookie.get('support_token')
        Vue.http.headers.common['Authorization'] = `SupportBearer ${token}`
        return Vue.http.get('support/sms/chat', {params: {user, lead}})
    },
}