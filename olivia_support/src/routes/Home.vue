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
								<Avatar v-if="message.for_lead" :user="selectedLead.realtor"></Avatar>
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

			<div class="tabs">
				<a href="#" :class="{ active: activeTab == 0 }" @click="setTab(0)">Script</a>
				<a href="#" :class="{ active: activeTab == 1 }" @click="setTab(1)">Notes</a>
				<a href="#" :class="{ active: activeTab == 2 }" @click="setTab(2)">Property Info</a>
			</div>

			<Card nopadding :class=" { 'profile-card': true, 'autoheight': this.selectedLead.lead.id == 20 }">
				<div class="profile">
					<Avatar big :user="this.selectedLead.lead"></Avatar>

					<div class="info">
						<div class="name-button-wrapper">
							<div class="name">{{ this.selectedLead.lead.name || 'N/A' }}</div>

							<button @click="setAssistance(selectedLead.lead, true)" class="take-ownership" v-if="!this.selectedLead.lead.needs_assistance">Request assistance</button>
							<button @click="setAssistance(selectedLead.lead, false)" class="take-ownership" v-if="this.selectedLead.lead.needs_assistance">Resolved</button>
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

			<div class="tab-content" v-show="activeTab == 0">
				<Card nopadding v-for="script in this.scripts" :key="script.data.id" class="script-card">
					<div class="heading">
						<div class="title">{{ script.data.title }}</div>
						<a href="javascript:void(0);" @click="message = script.data.script" class="copy"></a>
					</div>

					<div class="script" v-html="script.data.script"></div>

					<div class="answer" v-html="script.data.answer" v-if="!script.edited"></div>
					<textarea placeholder="Response" v-model="script.data.answer" v-if="script.edited"></textarea>

					<div class="controls">
						<a href="javascript:void(0)" @click="script.toggleEdit()" v-if="!script.edited">Edit</a>
						<a href="javascript:void(0)" @click="script.toggleEdit()" v-if="script.edited">Cancel</a>
						<a href="javascript:void(0)" @click="script.save(selectedLead.lead.id)" v-if="script.edited">Save</a>
					</div>
				</Card>
			</div>
			<div class="tab-content" v-show="activeTab == 1"></div>

			<div class="tab-content" v-show="activeTab == 2">
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

	class Script {
    constructor(script) {
      this.data = Object.assign({}, script)
			this.originalData = Object.assign({}, script)

			this.edited = false
		}

		toggleEdit() {
      if (this.edited) {
        this.data = Object.assign({}, this.originalData)
        this.edited = false
			} else {
        this.edited = true
			}
		}

		save(lead) {
      ScriptsResource.saveAnswer(this.data.id, lead, this.data.answer).then(response => {
        this.originalData = this.data
        this.edited = false
			})
		}
	}

  export default {
    components: {
      Page: RegularPage,
      Btn: Button,
      Card,
			Avatar,
    },
    methods: {
      setTab: function (index) {
        this.activeTab = index;
			},
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

				this.$forceUpdate()

				this.messages = []
				this.scripts = []

        SMSResource.getChat(this.selectedLead.lead.realtor, this.selectedLead.lead.id).then(
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

        ScriptsResource.getScriptsAndAnswers(this.selectedLead.lead.realtor, this.selectedLead.lead.id).then(
          response => {
            this.scripts = response.body.map(s => new Script(s))
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

        SMSResource.sendSMS(this.selectedLead.lead.realtor, leadId, newMessage.message).then(
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
			setAssistance: function (lead, flag) {
				LeadsResource.setAssistance(lead.id, flag).then(
          response => {
            lead.needs_assistance = response.body.assistance
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
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');

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
            this.selectLead(this.leads[0])
          }
        }
      )

			const leadsInterval = setInterval(() => {
        LeadsResource.getLeads().then(
          response => {
            this.leads = response.body

            SMSResource.getChat(this.selectedLead.lead.realtor, this.selectedLead.lead.id).then(
              response => {
                this.messages = response.body
                this.convertMessagesDates()
                this.$nextTick(this.scrollToBottom)
              },
              response => {
                console.error(response)
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
        scripts: [],
				currentLot: 1,
        leadsInterval,
        activeTab: 0,
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
					position: relative;
					width: 33.33%; height: 100%;
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
						height: calc(100% - 48px);
						overflow: auto;
						
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
									color: #ffffff;
									background-color: $blue;
								}

								.date {
									margin-right: 14px;
									text-align: right;
								}

								.avatar {
									background-color: #B465AE;
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

				.tabs {
					display: flex;
					margin-bottom: 19px;

					a {
						display: block;
						width: 33.333%;
						height: 36px;
						line-height: 36px;
						font-size: 14px;
						font-weight: 600;
						color: rgba($text, .45);
						text-align: center;
						transition: color .2s;

						&:after {
							display: block;
							content: '';
							width: 100%;
							height: 3px;
							background-color: transparent;
							box-shadow: none;

							transition: background-color .2s, box-shadow .2s;
						}

						&:hover {
							text-decoration: none;
							color: $text;
						}

						&.active {
							font-weight: 700;
							color: $text;

							&:after {
								display: block;
								content: '';
								width: 100%;
								height: 3px;
								background-color: $text;
								box-shadow: 0px 3px 6px rgba(#000000, .16);
							}
						}
					}
				}

				.tab-content {
					height: calc(100% - 193px - 24px - 36px - 19px);
					overflow: auto;
				}

				.script-card {
					padding: 17px 21px;

					&:last-child {
						margin-bottom: 0;
					}

					.heading {
						display: flex;
						justify-content: space-between;
						margin-bottom: 12px;

						.title {
							font-size: 14px;
							font-weight: 700;
						}

						.copy {
							display: block;
							width: 24px; height: 24px;
							border-radius: 12px;
							background-color: $blue;
							background-image: url('../assets/quote-sign.svg');
							background-repeat: no-repeat;
							background-position: center;
							background-size: 18px;
						}
					}

					.script {
						margin-bottom: 15px;
						font-size: 13px;
					}

					.answer {
						font-size: 13px;
					}

					textarea {
						width: 100%; height: 100px;
						padding: 7px 10px;
						border: 1px solid $border;
						box-sizing: border-box;
						font-size: 13px;
					}

					.controls {
						display: flex;
						margin-top: 10px;

						a {
							display: flex;
							align-items: center;
							margin-left: 4px;
							padding: 0 8px;
							height: 24px;
							font-size: 12px;
							background-color: $blue;
							color: #ffffff;
							border-radius: 12px;
							text-decoration: none;
							transition: color .1s;

							&:hover {
								color: #ffffff;
								background-color: lighten($blue, 5%);
							}

							i {
								margin-right: 5px;
								font-size: 16px;
							}

							span {
								font-size: 12px;
							}
						}
					}
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