import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
    sendSMS(lead, message) {
        const token = VueCookie.get('token')
        Vue.http.headers.common['Authorization'] = `Bearer ${token}`
        return Vue.http.post('sms/send', {lead, message})
    },
    getChat(lead) {
        const token = VueCookie.get('token')
        Vue.http.headers.common['Authorization'] = `Bearer ${token}`
        return Vue.http.get('sms/chat', {params: {lead}})
    },
}