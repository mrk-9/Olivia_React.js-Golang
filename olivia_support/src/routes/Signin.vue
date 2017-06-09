<template>
	<ModalPage id="signin">
		<form v-on:submit.prevent="submit" action="" method="post">
			<h3>Sign in</h3>

			<TextInput v-model="email" placeholder="Email" :error="errors.email"/>

			<TextInput type="password" v-model="password" placeholder="Password" :error="errors.password"/>

			<div class="footer">
				<!--<Checkbox v-model="remember" title="Remember me" />-->

				<Btn>Log in</Btn>
			</div>
		</form>
	</ModalPage>
</template>

<script>
  import ModalPage from './../modules/pages/Modal.vue'
  import Loader from './../modules/Loader.vue'
  import Input from './../modules/ui/Input.vue'
  import Button from './../modules/ui/Button.vue'
  import Checkbox from './../modules/ui/Checkbox.vue'

  import AuthResource from './../resources/Auth.js'

  class FormErrors {
    constructor() {
      this.reset()
    }

    reset() {
      this.email = ''
      this.password = ''
    }

    get allowed() {
      if (this.email || this.password) {
        return false
      }

      return true
    }
  }

  let formErrors = new FormErrors()

  export default {
    name: 'Auth',
    components: {
      Loader,
      ModalPage,
      TextInput: Input,
      Btn: Button,
      Checkbox,
    },
    methods: {
      submit: function (event) {
        this.errors.reset()
        this.loading = true

        if (!this.email) {
          this.errors.email = 'Email cannot be blank'
        }

        if (!this.password) {
          this.errors.password = 'Password cannot be blank'
        }

        if (!this.errors.allowed) {
          this.loading = false
          return
        }

        AuthResource.signIn(this.email, this.password).then(
          response => {
            this.$cookie.set('support_token', response.body.token, 30)
            this.loading = false
            this.$router.push('/')
          }, response => {
            this.errors.email = response.body.message
            this.loading = false
          }
        )
      }
    },
    data () {
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');

      return {
        loading: false,
        loaderColor: global.color,
        email: '',
        password: '',
        remember: false,
        errors: formErrors,
      }
    }
  }
</script>

<style lang="scss">
	#signin {
		h3 {
			margin-bottom: 30px;
		}

		form > .footer {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-top: 25px;

			a, button {
				font-size: 13px;
			}

			a {
				font-weight: bold;
			}
		}
	}
</style>