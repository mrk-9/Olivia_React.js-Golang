import Vue from 'vue'
import VueCookie from 'vue-cookie'
import moment from 'moment'

export default {
  getLeads() {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.get('leads/get').then(
      response => {
        let leads = response.body.map(l => {
          let date = (l.last_message.id != 0) ? l.last_message.created_at : l.lead.created_at

          if (l.lead.id == 20) {
            l.lead.image = require('./../assets/olivia-chat-logo.png')
          }

          if (!l.lead.interested_in) {
            l.lead.interested_in = [{id: 0}]
          }

          l.lead.created_at = moment.utc(date).local().calendar(null, {
            sameDay: 'LT',
            lastDay: 'ddd, LT',
            lastWeek: 'ddd',
            sameElse: 'MMM D'
          })

          return l
        })

        response.body = leads
        return response
      },
      response => {
        return response
      }
    )
  },
  setOwnership(lead, flag) {
    const token = VueCookie.get('token')
    Vue.http.headers.common['Authorization'] = `Bearer ${token}`
    return Vue.http.post('leads/setownership', {lead, flag})
  },
}