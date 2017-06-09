import Vue from 'vue'
import VueCookie from 'vue-cookie'

export default {
    getEmails() {
        const token = VueCookie.get('token')
        Vue.http.headers.common['Authorization'] = `Bearer ${token}`
        return Vue.http.get('emails/get')
    },
}