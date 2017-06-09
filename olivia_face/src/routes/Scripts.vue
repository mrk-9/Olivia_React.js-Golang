<template>
	<Page>
		<div class="column">
			<div class="title">Lead Info Messaging</div>

			<div class="item" v-for="script in scripts.filter((s) => !s.removed && s.id != 0)" :key="script.id">
				<Card class="script" v-if="!script.editing">
					<div class="heading">
						<div class="title" v-if="!script.editing">{{ script.title }}</div>
						<input class="editingTitle" type="text" v-if="script.editing" v-model="script.title">
						<div class="actions" v-if="!script.editing">
							<a @click="script.remove()" href="javascript:void(0);"><i class="material-icons">close</i> <span>Remove</span></a>
							<a @click="script.edit()" href="javascript:void(0);"><i class="material-icons">mode_edit</i> <span>Edit</span></a>
						</div>
						<div class="actions visible" v-if="script.editing">
							<a @click="script.cancelEditing()" href="javascript:void(0);"><i class="material-icons">close</i> <span>Cancel</span></a>
							<a @click="script.save()" href="javascript:void(0);"><i class="material-icons">mode_edit</i> <span>Save</span></a>
						</div>
					</div>
					<div class="text" v-if="!script.editing" v-html="script.text"></div>
					<textarea class="editingText" v-if="script.editing" v-model="script.text"></textarea>
				</Card>
				<Card class="form" v-if="script.editing">
					<div class="heading">
						<div class="script-title">
							<TextInput v-model="script.title" placeholder="Label"
												 :maxlength="maxlengths.title" nomargin></TextInput>

							<div class="length">{{ maxlengths.title - script.title.length }}</div>
						</div>

						<div class="actions">
							<a @click="script.cancelEditing()" href="javascript:void(0);"><i class="material-icons">close</i> <span>Cancel</span></a>
							<a @click="script.save()" href="javascript:void(0);" class="blue"><i class="material-icons">check</i> <span>Save</span></a>
						</div>
					</div>

					<TA v-model="script.text" placeholder="Message body"
							:maxlength="maxlengths.text" nomargin></TA>

					<div class="footer">
						<div :class="{ variables: true, visible: script.variablesVisible }">
							<a @click="insertVariable($event, v)" href="javascript:void(0);" v-for="title, v in variables">{{ title }}</a>
						</div>

						<div class="length">{{ maxlengths.text - script.text.length }}</div>
					</div>
				</Card>
			</div>

			<div class="divider"></div>

			<Card class="form" v-for="script, index in scripts.filter((s) => !s.removed && s.id == 0)" :key="index">
				<div class="heading">
					<div class="script-title">
						<TextInput v-model="script.title" placeholder="Label"
											 :maxlength="maxlengths.title" nomargin></TextInput>

						<div class="length">{{ maxlengths.title - script.title.length }}</div>
					</div>

					<div class="actions">
						<a @click="script.remove()" href="javascript:void(0);"><i class="material-icons">close</i> <span>Cancel</span></a>
						<a @click="script.create()" href="javascript:void(0);" class="blue"><i class="material-icons">check</i> <span>Save</span></a>
					</div>
				</div>

				<TA v-model="script.text" placeholder="Message body"
						:maxlength="maxlengths.text" nomargin></TA>

				<div class="footer">
					<div :class="{ variables: true, visible: script.variablesVisible }">
						<a @click="insertVariable($event, v)" href="javascript:void(0);" v-for="title, v in variables">{{ title }}</a>
					</div>

					<div class="length">{{ maxlengths.text - script.text.length }}</div>
				</div>
			</Card>

			<div class="add-another-message" v-if="!formVisible">
				<Btn :onclick="addEmptyScript" outline width="185">Add another message</Btn>
			</div>
		</div>

		<div class="column"></div>
	</Page>
</template>

<script>
  import RegularPage from './../modules/pages/Regular.vue'

  import Card from './../modules/ui/Card.vue'
  import Input from './../modules/ui/Input.vue'
  import Textarea from './../modules/ui/Textarea.vue'
	import Button from './../modules/ui/Button.vue'

  import ScriptsResource from './../resources/Scripts.js'

	class Script {
    constructor (object) {
      this.id = object.id
      this.title = object.title
			this.text = object.text

			this.originalTitle = ''
			this.originalText = ''

			this.editing = false

			this.removed = false

			this.variablesVisible = false
		}

		edit() {
      this.originalTitle = this.title
			this.originalText = this.text

      this.editing = true
		}

		cancelEditing() {
      this.title = this.originalTitle
			this.text = this.originalText

      this.editing = false
		}

		remove() {
      if (this.id == 0) {
        this.removed = true
			} else {
        ScriptsResource.removeScript(this.id).then(
          response => {
            this.removed = true
          }, response => {
            alert('Unable to remove script')
          }
        )
			}
		}

		save() {
      ScriptsResource.updateScript(this.id, this.title, this.text).then(
        response => {
          this.editing = false
				}, response => {

				}
			)
		}

		create() {
      if (this.title == '' || this.text == '') {
        return
      }

      ScriptsResource.createScript(this.title, this.text).then(
        response => {
					this.id = response.body.id
        }, response => {
          alert('Unable to create script')
        }
      )
		}
	}

  export default {
    components: {
      Page: RegularPage,
      TextInput: Input,
      TA: Textarea,
      Btn: Button,
      Card,
    },
    methods: {
      insertVariable: function (e, variable) {
        variable = `{${variable}}`

        let ta = e.path[3].querySelector('textarea')

        //IE support
        if (document.selection) {
          ta.focus();
          sel = document.selection.createRange();
          sel.text = variable;
        }
        //MOZILLA and others
        else if (ta.selectionStart || ta.selectionStart == '0') {
          var startPos = ta.selectionStart;
          var endPos = ta.selectionEnd;
          ta.value = ta.value.substring(0, startPos)
            + variable
            + ta.value.substring(endPos, ta.value.length);
          ta.selectionStart = startPos + variable.length;
          ta.selectionEnd = startPos + variable.length;
        } else {
          ta.value += variable;
        }

        let event = new Event('input')
				ta.dispatchEvent(event)
      },
      getScripts: function () {
        ScriptsResource.getScripts().then(
          response => {
            this.scripts = response.body.map((script) => new Script(script))

            this.addEmptyScript()
          }, response => {
            console.error(response)
          }
        )
      },
			addEmptyScript: function () {
        let script = new Script({
					id: 0,
					title: '',
					text: '',
				})

        this.scripts.push(script)
			},
    },
    data: function () {
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');
      fbq('track', 'PageView');


      this.getScripts()

      return {
        scripts: [],
        maxlengths: {
          title: 40,
          text: 200,
        },
        script: {
          title: '',
          text: '',
        },
        editingScript: {
          id: 0,
					title: '',
					text: '',
				},
        variables: {
          'first_name': 'first name',
          'last_name': 'last name',
          'sale_price': 'sale price',
          'address': 'address',
          'description': 'description',
        },
				formVisible: false,
      }
    },
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../assets/scss/_settings.scss';

	.wrapper {
		display: flex;

		.column {
			width: 50%;
			padding-bottom: 20px;

			> .title {
				margin-bottom: 30px;
				padding-bottom: 25px;
				font-size: 18px;
				font-weight: 700;
				border-bottom: 1px solid $border;
			}

			.script {
				margin-bottom: 14px;
				padding: 25px;

				.heading {
					display: flex;
					justify-content: space-between;
					margin-bottom: 23px;

					.title, input.editingTitle {
						height: 19px;
						font-size: 14px;
						font-weight: 700;
					}

					input.editingTitle {
						width: 100%; height: 23px;
						margin: -2px 0 -2px -2px;
						border: 1px solid $border;
						box-sizing: border-box;
					}

					.actions {
						display: flex;

						a {
							display: flex;
							align-items: center;
							margin-left: 37px;
							color: $placeholder;
							text-decoration: none;
							transition: color .1s;

							&:hover {
								color: $text;
							}

							i {
								margin-right: 5px;
								font-size: 16px;
							}

							span {
								font-size: 12px;
								font-weight: 600;
							}
						}
					}
				}

				.text, textarea.editingText {
					font-size: 13px;
				}

				textarea.editingText {
					width: 100%; min-height: 100px;
					margin: -3px 0 -2px -3px;
					border: 1px solid $border;
				}
			}

			.divider {
				margin: 30px 0;
				border-bottom: 1px solid $border;
			}

			.form {
				margin-bottom: 14px;
				padding: 19px 21px;

				.heading {
					display: flex;
					justify-content: space-between;
					align-items: center;
					margin-bottom: 12px;

					.script-title {
						display: flex;
						align-items: center;
						width: calc(100% - 160px);

						.length {
							width: 16px;
							margin-left: 8px;
							color: $placeholder;
							font-size: 13px;
						}
					}

					.actions {
						display: flex;
						justify-content: flex-end;
						width: 150px;

						a {
							display: flex;
							align-items: center;
							margin-left: 4px;
							padding: 0 8px;
							height: 24px;
							color: $placeholder;
							text-decoration: none;
							transition: color .1s;

							&.blue {
								background-color: $blue;
								color: #ffffff;
								border-radius: 12px;

								&:hover {
									color: #ffffff;
									background-color: lighten($blue, 5%);
								}
							}

							&:hover {
								color: $text;
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

				textarea {
					height: 90px;
					resize: none;
				}

				.footer {
					display: flex;
					align-items: center;
					justify-content: space-between;
					height: 24px;

					.variables {
						a {
							display: inline-block;
							height: 24px;
							padding: 0 10px;
							line-height: 24px;
							color: $placeholder;
							font-size: 12px;
							border-radius: 12px;
							background-color: rgba($placeholder, .1);
							transition: background-color .1s;

							&:hover {
								background-color: rgba($placeholder, .2);
								text-decoration: none;
							}
						}
					}

					.length {
						width: 25px;
						margin-left: 8px;
						color: $placeholder;
						font-size: 13px;
					}
				}
			}

			.add-another-message {
				padding-top: 25px;
				padding-bottom: 20px;
				text-align: center;
				
				button {
					font-size: 13px;
				}
			}
		}
	}
</style>