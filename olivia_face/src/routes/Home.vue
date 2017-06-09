<template>
	<Page class="home">
		<div :class="{ 'column-2-3': this.isLeadSelected, 'column-1': !this.isLeadSelected }">
			<Card nopadding class="chat-card">
				<div class="clients">
					<div class="header">
						<i class="material-icons search-icon">search</i>
						<input type="search" placeholder="Search in clients">
					</div>
					<div class="list">
						<div @click="selectLead(lead)" v-for="lead in this.leads"
								 :class="{ client: true, active: (lead.lead.id == selectedLead.lead.id) }">
							<Avatar :user="lead.lead"></Avatar>
							<div class="info">
								<div class="namedate">
									<div class="name">{{ lead.lead.name || 'N/A' }}</div>
									<div class="date">{{ lead.lead.created_at }}</div>
								</div>
								<div class="message">{{ lead.last_message.message }}</div>
							</div>
						</div>
					</div>
				</div>
				<div class="dialog">
					<div class="header" v-if="this.isLeadSelected">
						<div class="name">{{ this.selectedLead.lead.name }}</div>
						<div class="lastactive" v-if="this.messages.length > 0">Last active {{ this.messages[this.messages.length - 1].created_at | moment("from") }}</div>
					</div>

					<div class="messages">
						<div v-for="message in this.messages">
							<div class="date-divider" v-if="message.dateDivider">{{ message.dateDivider }}</div>
							<div :class="{ message: true, forLead: message.for_lead, loading: message.loading }">
								<Avatar v-if="message.for_lead" :user="user"></Avatar>
								<Avatar v-if="!message.for_lead" :user="selectedLead.lead"></Avatar>
								<div class="body" v-html="message.message"></div>
								<div class="date">{{ message.created_at_formatted }}</div>
							</div>
						</div>
					</div>

					<div class="input-wrapper">
						<textarea placeholder="Your message..." v-model="message"></textarea>
						<button @click="sendSMS"><i class="material-icons">send</i></button>
					</div>
				</div>
			</Card>
		</div>

		<div :class="{ 'profile-column': true, 'column-1-3': this.isLeadSelected, 'hidden': !this.isLeadSelected }"
				 v-if="this.isLeadSelected">
			<Card nopadding :class=" { 'profile-card': true, 'autoheight': this.selectedLead.lead.id == 20 }">
				<div class="profile">
					<Avatar big :user="this.selectedLead.lead"></Avatar>

					<div class="info">
						<div class="name-button-wrapper">
							<div class="name">{{ this.selectedLead.lead.name || 'N/A' }}</div>

							<button @click="setOwnership(selectedLead.lead, true)" class="take-ownership" v-if="this.selectedLead.lead.id != 20 && this.selectedLead.lead.realtor != this.user.id">Take ownership</button>
							<button @click="setOwnership(selectedLead.lead, false)" class="take-ownership" v-if="this.selectedLead.lead.id != 20 && this.selectedLead.lead.realtor == this.user.id">Release ownership</button>
						</div>
						<div class="phone" v-if="this.selectedLead.lead.id != 20" v-show="this.selectedLead.lead.phone">
							<i class="material-icons">phone</i> {{ this.selectedLead.lead.phone }}
						</div>
						<div class="email" v-show="this.selectedLead.lead.email">
							<i class="material-icons">email</i> {{ this.selectedLead.lead.email }}
						</div>

						<div class="pair" v-if="this.selectedLead.lead.id != 20">
							<div class="key">Status:</div>
							<div class="value">{{ this.selectedLead.email.status }}</div>
						</div>

						<div class="pair" v-if="this.selectedLead.lead.id != 20">
							<div class="key">Enquiry type:</div>
							<div class="value">{{ this.selectedLead.email.enquiry_type }}</div>
						</div>
					</div>
				</div>
			</Card>

			<div :class=" {'lots-controls': true} " v-if="this.selectedLead.lead.interested_in.length > 1">
				<button @click="this.prevLot">Prev</button>
				<div class="position">{{ this.currentLot }}/{{ this.selectedLead.lead.interested_in.length }}</div>
				<button @click="this.nextLot">Next</button>
			</div>
			<div :class="{ 'lot-card-wrapper': true, multilot: this.selectedLead.lead.interested_in.length > 1, 'autoheight': this.selectedLead.lead.id == 20 }" v-if="this.isLeadSelected">
				<div class="lot-cards" :style="{ 'margin-left': '-' + ((this.currentLot - 1) * (445 + 20)) + 'px', width: ((this.selectedLead.lead.interested_in.length * (445 + 20)) - 20) + 'px' }">
					<Card nopadding v-for="property in this.selectedLead.lead.interested_in" class="lot-card" :key="property.id">
						<div class="section twocol" v-if="property.id != 0">
							<div class="info">
								<div class="title">Property Information</div>

								<div class="text" v-html="property.address"></div>
								<div class="text" v-html="property.description"></div>

								<div class="pair">
									<div class="key">Sale price:</div>
									<div class="value">{{ property.price }}</div>
								</div>
							</div>

							<div class="image">
								<a :href="property.link" target="_blank">
									<img :src="property.image_url"
											 :alt="property.description">
								</a>
							</div>
						</div>
						<div class="divider" v-if="property.id != 0"></div>
						<div class="section" v-if="answers.length > 0">
							<div class="title">Requested Information</div>

							<div class="pair" v-for="answer in answers">
								<div class="key">{{ answer.script_title }}</div>
								<div class="value">{{ (answer.text) ? answer.text : '&ndash;' }}</div>
							</div>
						</div>
						<div class="divider"></div>
						<div class="section">
							<div class="title">Additional Information</div>
							<div class="nodata">No Data Available</div>

							<!--<div class="row">-->
								<!--<div class="pair-inline">-->
									<!--<span class="title">Age:</span>-->
									<!--<span class="value">64</span>-->
								<!--</div>-->

								<!--<div class="socials">-->
									<!--<a href="#">F</a>-->
									<!--<a href="#">T</a>-->
									<!--<a href="#">L</a>-->
								<!--</div>-->
							<!--</div>-->

							<!--<div class="row">-->
								<!--<div class="pair-inline">-->
									<!--<div class="title">Career</div>-->
									<!--<div class="value">-->
										<!--Internet <br>-->
										<!--HTML <br>-->
										<!--WinAMP-->
									<!--</div>-->
								<!--</div>-->

								<!--<div class="pair-inline">-->
									<!--<div class="title">Places</div>-->
									<!--<div class="value">-->
										<!--Internet <br>-->
										<!--Extranet-->
									<!--</div>-->
								<!--</div>-->
							<!--</div>-->
						</div>
					</Card>
				</div>
			</div>
		</div>
	</Page>
</template>

<script>
  import moment from "moment"
  import AuthResource from './../resources/Auth.js'
  import LeadsResource from './../resources/Leads.js'
  import SMSResource from './../resources/SMS.js'
  import ScriptsResource from './../resources/Scripts.js'

  import RegularPage from './../modules/pages/Regular.vue'

  import Button from './../modules/ui/Button.vue'
  import Card from './../modules/ui/Card.vue'
	import Avatar from './../modules/ui/Avatar.vue'

  export default {
    components: {
      Page: RegularPage,
      Btn: Button,
      Card,
			Avatar,
    },
    methods: {
      convertMessagesDates: function () {
        if (!this.messages) {
          return
				}

        let lastDateDivider = ''

        this.messages = this.messages.map((message) => {
          let dateDivider = moment.utc(message.created_at).local().calendar(null, {
            sameDay: '[Today]',
            lastDay: 'MMMM D, dddd',
            lastWeek: 'MMMM D, dddd',
            sameElse: 'MMMM D, dddd'
          })

					if (dateDivider != lastDateDivider) {
            lastDateDivider = dateDivider
					} else {
            dateDivider = undefined
					}

					message.created_at_formatted = moment.utc(message.created_at).local().format('LT')

					message.dateDivider = dateDivider

          return message
        })
			},
      selectLead: function (lead) {
        this.selectedLead = undefined
        this.selectedLead = lead

//				window.location.hash = '#chat:' + lead.lead.id

				this.$forceUpdate()

				this.messages = []

        SMSResource.getChat(this.selectedLead.lead.id).then(
          response => {
            this.messages = response.body
						this.convertMessagesDates()
            this.$nextTick(this.scrollToBottom)
          },
          response => {
            console.error(response)
          }
        )

				this.answers = []

        ScriptsResource.getScriptsAnswers(this.selectedLead.lead.id).then(
          response => {
            this.answers = response.body
					},
          response => {
            console.error(response)
          }
				)
      },
      sendSMS: function () {
        if (!this.message) {
          return
				}

        var leadId = this.selectedLead.lead.id

        var newMessage = {
          for_lead: true,
          message: this.message,
          loading: true,
        }

        this.message = ""

        this.messages.push(newMessage)
				this.convertMessagesDates()
        this.$nextTick(this.scrollToBottom)

        SMSResource.sendSMS(leadId, newMessage.message).then(
          response => {
            newMessage.loading = false
          },
          response => {
          }
        )
      },
      scrollToBottom: function () {
        var messagesWrapper = document.querySelector('.messages')
        messagesWrapper.scrollTop = messagesWrapper.scrollHeight
      },
			prevLot: function () {
				if (this.currentLot > 1) {
				  this.currentLot--;
				}
			},
			nextLot: function () {
				if (this.currentLot < this.selectedLead.lead.interested_in.length) {
				  this.currentLot++;
				}
			},
			setOwnership: function (lead, flag) {
				LeadsResource.setOwnership(lead.id, flag).then(
          response => {
            if (response.body.owning) {
							lead.realtor = this.user.id
            } else {
              lead.realtor = 0
						}
          }
        )
			},
    },
    computed: {
      isLeadSelected: function () {
        return (!!this.selectedLead.lead && !!this.selectedLead.lead.id != 0)
      },
    },
    beforeDestroy() {
      clearInterval(this.leadsInterval)
		},
    data () {
      if (!this.$cookie.get('token')) {
        this.$router.push('/auth/signin')
        return
      }

      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');
      fbq('track', 'PageView');

      AuthResource.getProfile().then(
        response => {
          this.user = response.body

          this.ws = new WebSocket(`ws://${global.host}/ws?id=${this.user.id}`)

          this.ws.onopen = () => {
            this.ws.send(JSON.stringify({
              type: "hello",
              data: {
                userId: '' + this.user.id,
              },
            }))

            this.ws.onmessage = (event) => {
              var message = JSON.parse(event.data)

              if (message.type == 'message') {
                message.created_at = (message.created_at) ? message.created_at : Date.now()

                this.messages.push(message.message)
                this.convertMessagesDates()
                this.$nextTick(this.scrollToBottom)
              }
            }
          }
        }
      )

      LeadsResource.getLeads().then(
        response => {
          this.leads = response.body

          if (this.leads.length > 0) {
            var leadIndex = 0

						if (this.$route.query) {
							let leadId = this.$route.query.chat

							for (let i = 0; i < this.leads.length; i++) {
							  if (this.leads[i].lead.id == leadId) {
							    leadIndex = i;
								}
							}
						}

            this.selectLead(this.leads[leadIndex])
          }
        }
      )

			const leadsInterval = setInterval(() => {
        LeadsResource.getLeads().then(
          response => {
            this.leads = response.body

            SMSResource.getChat(this.selectedLead.lead.id).then(
              response => {
                this.messages = response.body
                this.convertMessagesDates()
                this.$nextTick(this.scrollToBottom)
              },
              response => {
                console.error(responseƒ)
              }
            )
          }
				)
			}, 5000)

      return {
        user: {},
        leads: [],
        selectedLead: {
          id: 0,
					name: '',
				},
        messages: [],
        message: "",
        ws: {},
				answers: [],
				currentLot: 1,
        leadsInterval,
      }
    }
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../assets/scss/_settings.scss';

	.page.home {
		.wrapper {
			display: flex;
			justify-content: space-between;

			.column-1-3, .column-2-3 {
				height: 100%;
			}

			.column-2-3 {
				width: calc(100% - 445px - 20px) !important;
			}

			.card {
				margin-bottom: 24px;
			}

			.chat-card {
				display: flex;
				flex-direction: row;
				justify-content: space-between;
				align-items: stretch;
				height: 100%;

				.clients, .dialog {
					box-sizing: border-box;

					.header {
						width: 100%;
						height: 47px;
						border-bottom: 1px solid $border;
						line-height: 47px;
					}
				}

				.clients {
					width: 33.33%;
					min-width: 300px;
					border-right: 1px solid $border;

					.header {
						display: flex;
						align-items: center;

						.search-icon {
							display: block;
							margin: 0 16px;
							color: $placeholder;
							font-size: 20px;
						}

						input[type=search] {
							width: 100%;
							margin-right: 16px;
							border: none;
							font-size: 13px;
							outline: none;
							color: $text;
							-webkit-appearance: none;
							-moz-appearance: none;
							appearance: none;

							@include placeholder {
								color: $placeholder;
							}
						;
						}
					}

					.list {
						.client {
							$transTime: .05s;

							display: flex;
							width: 100%;
							height: 73px;
							padding: 16px 19px;
							box-sizing: border-box;
							border-bottom: 1px solid $border;
							cursor: pointer;
							transition: background-color $transTime;

							.avatar {
								margin-right: 12px;
							}

							.info {
								width: calc(100% - 42px - 12px);

								.namedate {
									display: flex;
									justify-content: space-between;

									.name, .date {
										transition: color $transTime;
									}

									.name {
										font-size: 14px;
										font-weight: 700;
									}

									.date {
										font-size: 11px;
										color: rgba($placeholder, .7);
									}
								}

								.message {
									font-size: 13px;
									color: $placeholder;
									transition: color $transTime;
									text-overflow: ellipsis;
									overflow: hidden;
									white-space: nowrap;
								}
							}

							&:hover, &.active {
								background-color: $blue;

								.info {
									.namedate {
										.name, .date {
											color: #ffffff;
										}
									}

									.message {
										color: rgba(#ffffff, .8);
									}
								}
							}
						}
					}
				}

				.dialog {
					display: flex;
					flex-direction: column;
					justify-content: center;
					width: 66.66%;
					height: 100%;
					position: relative;

					.header {
						display: flex;
						align-items: center;

						.name, .lastactive {
							font-size: 14px;
						}

						.name {
							font-weight: 700;
							padding-left: 16px;
							padding-right: 6px;
						}

						.lastactive {
							color: $placeholder;
							padding-right: 16px;
						}
					}

					.messages {
						width: 100%;
						height: calc(100% - 47px - 81px);
						padding: 18px;
						box-sizing: border-box;
						overflow: auto;

						.date-divider {
							margin: 35px 0 20px;
							font-size: 13px;
							color: rgba($placeholder, .7);
							text-align: center;
						}

						.message {
							display: flex;
							margin-top: 16px;

							&.loading {
								opacity: .5;
							}

							&.forLead {
								flex-direction: row-reverse;

								.body {
									margin-right: 12px;
									background-color: rgba(#31ACF8, .12);
								}

								.date {
									margin-right: 14px;
									text-align: right;
								}
							}

							&:not(.forLead) {
								.body {
									margin-left: 12px;
									background-color: rgba($placeholder, .1);
								}

								.date {
									margin-left: 14px;
									text-align: left;
								}
							}

							.avatar {
								min-width: 42px;
							}

							.body {
								padding: 12px 16px;
								font-size: 13px;
								border-radius: 3px;
							}

							.date {
								width: 60px;
								min-width: 60px;
								height: 42px;
								color: $placeholder;
								font-size: 13px;
								line-height: 42px;
							}
						}
					}

					.input-wrapper {
						display: flex;
						justify-content: space-between;
						align-items: center;
						width: 100%;
						height: 81px;
						padding: 14px 18px;
						box-sizing: border-box;
						border-top: 1px solid $border;

						textarea {
							width: 100%;
							height: 100%;
							@include placeholder {
								color: $placeholder;
							}
						;
							font-size: 13px;
							border: none;
							resize: none;
						}

						button {
							border: none;
							background: transparent;
							color: $blue;
							cursor: pointer;
						}
					}
				}
			}

			.profile-column {
				$colWidth: 445px;

				width: $colWidth;
				height: 100%;
				margin: -5px;
				padding: 5px;
				overflow: auto;

				&.hidden {
					display: none;
				}

				.profile-card {
					&.autoheight {
						.profile {
							height: auto;
						}
					}

					.profile {
						display: flex;
						justify-content: space-between;
						height: 149px;
						padding: 22px 22px;

						.info {
							width: calc(100% - 82px - 22px);
							padding-top: 4px;

							.name-button-wrapper {
								display: flex;
								justify-content: space-between;
								align-items: center;
								margin-bottom: 15px;

								.name {
									font-size: 14px;
									font-weight: 700;
								}

								button {
									width: 123px;
									height: 24px;
									background-color: $blue;
									color: #ffffff;
									border-radius: 12px;
									-webkit-appearance: none;
									-moz-appearance: none;
									appearance: none;
									border: none;
									cursor: pointer;
									font-size: 12px;
									font-weight: 600;
								}
							}

							.phone, .email {
								display: flex;
								align-items: center;
								padding-bottom: 5px;
								font-size: 13px;

								i {
									padding-right: 8px;
									font-size: 16px;
								}
							}

							.pair {
								display: flex;
								padding-top: 12px;

								.key, .value {
									font-size: 13px;
								}

								.key {
									width: 40%;
								}

								.value {
									width: 60%;
									color: $lightText;
								}
							}
						}
					}

					.button-wrapper {
						display: flex;
						justify-content: center;
						align-items: center;
						padding-bottom: 20px;
					}
				}

				.lots-controls {
					display: flex;
					justify-content: space-between;
					width: 100%;
					height: 24px;
					margin-bottom: 5px;

					button {

					}

					.position {
						font-size: 13px;
						color: $placeholder;
					}
				}

				.lot-card-wrapper {
					height: calc(100% - 193px - 24px);
					overflow-y: auto;
					@include shadow();

					&.autoheight {
						height: auto;
					}

					&.multilot {
						height: calc(100% - 193px - 24px - 29px);
					}

					.lot-cards {
						display: flex;

						.lot-card {
							width: $colWidth;
							margin: 0 20px 0 0;
							padding: 20px 0;

							&:nth-last-child(1) {
								margin-right: 0;
							}

							.section {
								padding: 0 20px;

								&.twocol {
									display: flex;

									> div {
										width: 50%;

										&.info {
											padding-right: 8px;
											box-sizing: border-box;
										}

										&.image {
											position: relative;

											img {
												width: 100%;
												border-radius: 2px;
											}
										}
									}
								}

								.title {
									padding-bottom: 4px;
									font-size: 14px;
									font-weight: 700;
								}

								.nodata {
									font-size: 14px;
								}

								.pair {
									display: flex;
									padding-top: 12px;

									.key, .value {
										width: 50%;
										font-size: 13px;
									}

									.key {
										font-weight: 600;
									}

									.value {
										color: $lightText;
									}
								}

								.pair-inline {
									.title, .value {
										font-size: 13px;
									}

									.value {
										color: $lightText;
									}
								}

								.text {
									padding-top: 12px;
									font-size: 13px;
								}

								.row {
									display: flex;
									padding-top: 12px;

									> div {
										width: 50%;
									}
								}
							}
						}
					}
				}
			}
		}
	}
</style>