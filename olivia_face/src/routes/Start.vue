<template>
	<div id="landing">
		<div class="heading-wrapper">
			<div class="heading">
				<div class="topbar">
					<div class="logo"></div>

					<div class="menu">
						<a href="#features">Features</a>
						<router-link to="/plans">Pricing</router-link>
						<router-link to="/auth/signin">Log in</router-link>
						<router-link to="/auth/signup" class="outline">Sign up</router-link>
					</div>
				</div>

				<div class="slogan">
					Respond to inbound leads faster
				</div>

				<div class="small-slogan">
					Our real estate concierges respond to your prospects within minutes, 24/7
				</div>

				<router-link to="/auth/signup" class="trial">Start your free trial</router-link>
			</div>
		</div>

		<div class="content-wrapper">
			<a name="features"></a>
			<div class="content">
				<div class="row">
					<div class="text">
						<div class="title">Instantly Communicate With Prospects</div>

						<p>Contact your prospects when they’re most interested, on the device they check most often. Olivia’s agents communicate with your leads within minutes of their request.</p>
						<p>You can view all of their conversations and take over whenever you wish.</p>
					</div>

					<div class="image f1"></div>
				</div>

				<div class="row">
					<div class="image f2"></div>

					<div class="text">
						<div class="title">Get More Information</div>

						<p>Whenever possible, we’ll give you more detailed information about your prospects including, job titles, companies, education history, social media profiles, and more.</p>
					</div>
				</div>

				<div class="row">
					<div class="text">
						<div class="title">Create Your Own Script</div>

						<p>Olivia’s agents use the script that you provide to ensure that your leads are taken care of and that you have all the information that you need to make a sale.</p>
					</div>

					<div class="image f3"></div>
				</div>

				<div class="row">
					<div class="image f4"></div>

					<div class="text">
						<div class="title">Get An Update Every Morning</div>

						<p>Receive an up to the minute report in your inbox daily alerting you to new prospects received overnight and updates on existing leads.</p>
					</div>
				</div>
			</div>
		</div>

		<div class="footer1-wrapper">
			<div class="footer1">
				<div class="title">Start your free 10 day trial now!</div>
				<router-link to="/auth/signup" class="button">Start your free trial</router-link>
			</div>
		</div>

		<div class="footer2-wrapper">
			<div class="footer2">
				<div class="logo"></div>
				<div class="copy">&copy; 2017 Olivia</div>
				<div class="menu">
					<router-link to="/plans">Pricing</router-link>
					<a href="mailto:support@getolivia.co">Contact us</a>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
  import Logo from './../modules/Logo.vue'
  import Button from './../modules/ui/Button.vue'

  import intercom from './../modules/vendor/intercom'

  export default {
    components: {
      Logo: Logo,
      Btn: Button,
    },
    name: 'Home',
    data () {
      intercom(null, null, null, true);
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');
      fbq('track', 'PageView');


      var vm = this

      return {
        goSignUp: function () {
          vm.$router.push('/auth/signup')
          console.log('kek')
        },
      }
    }
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../assets/scss/_settings.scss';

	$dark: #1f333d;

	#landing {
		width: 100%; height: 100%;

		$innerWidth: 1173px;

		.heading-wrapper {
			width: 100%; height: 860px;
			background-color: $blue;
			background-image: url('./../assets/landing-bg.png');
			background-position: center;
			background-repeat: no-repeat;
			background-size: cover;

			.heading {
				max-width: $innerWidth;
				width: 100%;
				margin: 0 auto 0;
				padding-top: 38px;

				.topbar {
					display: flex;
					justify-content: space-between;
					align-items: center;

					.logo {
						width: 172px; height: 33px;
						background-image: url('./../assets/logo-white.png');
					}

					.menu {
						display: flex;
						align-items: center;

						a {
							display: block;
							margin-left: 56px;
							text-transform: uppercase;
							color: #ffffff;
							font-size: 13px;
							font-weight: 700;

							&.outline {
								width: 100px; height: 42px;
								border: 2px solid #ffffff;
								line-height: 42px;
								text-align: center;
								border-radius: 21px;
								box-sizing: border-box;
							}
						}
					}
				}

				.slogan {
					margin: 205px auto 0;
					font-size: 36px;
					font-weight: 700;
					color: #ffffff;
					text-align: center;
					text-shadow: rgba($dark, .3) 0 6px 12px;
				}

				.small-slogan {
					margin-top: 45px;
					font-size: 17px;
					font-weight: 600;
					line-height: 26px;
					color: #ffffff;
					text-align: center;
					text-shadow: rgba($dark, .3) 0 6px 12px;
				}

				.trial {
					display: block;
					width: 230px; height: 48px;
					margin: 56px auto 0;
					background-color: $blue;
					border-radius: 24px;
					font-size: 13px;
					font-weight: 700;
					color: #ffffff;
					line-height: 48px;
					text-align: center;
					text-transform: uppercase;
					letter-spacing: 1px;
					box-shadow: rgba($dark, .3) 0 6px 12px;
					transition: background-color .2s;

					&:hover {
						text-decoration: none;
						background-color: lighten($blue, 3%);
					}
				}
			}
		}

		.content-wrapper {
			.content {
				max-width: $innerWidth;
				width: 100%;
				margin: 0 auto 0;
				padding: 128px 0 0;
				box-sizing: border-box;

				.row {
					display: flex;
					align-items: center;
					margin-bottom: 120px;

					.text, .image {
						width: calc(50% - 118px);
						margin: 0 59px;
						box-sizing: border-box;
					}

					.text {
						padding-left: 20px;

						.title {
							padding-bottom: 20px;
							color: $blue;
							font-size: 26px;
							font-weight: 700;
						}

						p {
							font-size: 17px;
							color: #757d85;
							line-height: 26px;
						}
					}

					.image {
						width: 491px; height: 349px;
						background-position: center;
						background-repeat: no-repeat;
						background-size: contain;

						&.f1 {
							background-image: url('./../assets/figures/1.png');
						}

						&.f2 {
							background-image: url('./../assets/figures/2.png');
						}

						&.f3 {
							background-image: url('./../assets/figures/3.png');
						}

						&.f4 {
							background-image: url('./../assets/figures/4.png');
						}
					}
				}
			}
		}

		.footer1-wrapper {
			width: 100%;
			height: 390px;
			background-color: $blue;

			.footer1 {
				max-width: $innerWidth;
				width: 100%; height: 100%;
				margin: 0 auto 0;

				background-image: url('./../assets/footer1.png');
				background-repeat: no-repeat;
				background-position: center bottom;
				background-size: contain;

				.title {
					padding-top: 125px;
					font-size: 36px;
					font-weight: 700;
					color: #ffffff;
					text-align: center;
				}

				.button {
					display: block;
					width: 230px; height: 48px;
					margin: 57px auto 0;
					font-size: 13px;
					font-weight: 700;
					color: #ffffff;
					line-height: 44px;
					text-align: center;
					box-sizing: border-box;
					border: 2px solid #ffffff;
					border-radius: 24px;
					text-decoration: none;
				}
			}
		}

		.footer2-wrapper {
			width: 100%; height: 95px;

			.footer2 {
				display: flex;
				justify-content: space-between;
				align-items: center;
				max-width: $innerWidth;
				width: 100%;
				height: 100%;
				margin: 0 auto 0;

				.logo {
					width: 172px; height: 33px;
					background-image: url('./../assets/logo2.png');
				}

				.copy {
					color: #c3d3e2;
					font-size: 14px;
					text-align: center;
				}

				.menu {
					a {
						margin-left: 50px;
						color: #757d85;
						font-size: 14px;
					}
				}
			}
		}
	}
</style>