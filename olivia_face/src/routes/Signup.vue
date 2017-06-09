<template>
	<ModalPage id="signup">
		<form v-on:submit.prevent="submit" action="" method="post">
			<h3>Sign up</h3>

			<TextInput v-model="fname" placeholder="First Name" :error="errors.fname"/>

			<TextInput v-model="sname" placeholder="Last Name" :error="errors.sname"/>

			<TextInput v-model="email" placeholder="Email" :error="errors.email"/>

			<PhoneInput v-model="phone" placeholder="Phone number" :error="errors.phone"/>

			<TextInput type="password" v-model="password" placeholder="Password" :error="errors.password"/>

			<TextInput type="password" v-model="cpassword" placeholder="Confirm password" :error="errors.cpassword"/>

			<div class="footer">
				<router-link to="/auth/signin">Log in</router-link>

				<Btn>Sign up</Btn>
			</div>
		</form>
	</ModalPage>
</template>

<script>
  import ModalPage from './../modules/pages/Modal.vue'
  import Loader from './../modules/Loader.vue'
  import Input from './../modules/ui/Input.vue'
  import PhoneInput from './../modules/ui/PhoneInput.vue'
  import Button from './../modules/ui/Button.vue'

  import AuthResource from './../resources/Auth.js'

  class FormErrors {
    constructor() {
      this.reset()
    }

    reset() {
      this.fname = ''
      this.sname = ''
      this.email = ''
			this.code = ''
      this.phone = ''
      this.password = ''
      this.cpassword = ''
    }

    get allowed() {
      if (this.fname || this.sname || this.email || this.phone || this.password || this.cpassword) {
        return false
      }

      return true
    }
  }

  let formErrors = new FormErrors()

  export default {
    name: 'Auth',
    components: {
      ModalPage,
      Loader,
      TextInput: Input,
			PhoneInput: PhoneInput,
      Btn: Button,
    },
    methods: {
      submit: function (event) {
        this.errors.reset()

        if (!this.fname) {
          this.errors.fname = 'First name cannot be blank'
        }

        if (!this.sname) {
          this.errors.sname = 'Second name cannot be blank'
        }

        if (!this.email) {
          this.errors.email = 'Email cannot be blank'
        }

        if (!(/^\S+@\S+$/).test(this.email)) {
          this.errors.email = 'Email should be valid'
        }

        if (!this.phone.code) {
          this.errors.phone = 'Country code is required'
				}

        if (!this.phone.number) {
          this.errors.phone = 'Phone number cannot be blank'
        }

        if (!this.password) {
          this.errors.password = 'Password cannot be blank'
        }

        if (!this.cpassword) {
          this.errors.cpassword = 'Confirm password cannot be blank'
        }

        if (this.password != this.cpassword) {
          this.errors.cpassword = 'Passwords should be the same'
        }

        if (!this.errors.allowed) {
          return
        }

        this.loading = true

				Promise.all([
				  AuthResource.validatePhoneNumber(this.phone)
				]).then((results) => {
          if (!results[0].body.valid) {
            this.errors.phone = 'Phone number appears to be invalid'
						return
					}

          AuthResource.signUp(this.fname, this.sname, this.email, this.phone, this.password).then(
            response => {
              AuthResource.signIn(this.email, this.password).then(
                response => {
                  this.$cookie.set('token', response.body.token, 30)
                  this.loading = false;
                  fbq('trackCustom', 'SignUp');
                  this.$router.push('/home')
                }, response => {
                  this.errors.email = response.body.message
                  this.loading = false
                }
              )
            }, response => {
              this.errors.email = response.body.message
              this.loading = false
            }
          )
				})
      }
    },
    data () {
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');
      fbq('track', 'PageView');


      return {
        loading: false,
        loaderColor: global.color,
        fname: '',
        sname: '',
        email: '',
        phone: {code: '', number: ''},
        password: '',
        cpassword: '',
        errors: formErrors,
      }
    }
  }
</script>

<style lang="scss">
	#signup {
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