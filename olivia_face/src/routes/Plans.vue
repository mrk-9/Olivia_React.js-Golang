<template>
	<Page class="plans">
		<div class="title">Plans & Pricing</div>
		<div class="subtitle">
			Increase your sales and decrease the time you spend following up on leads. <br>
			Let Olivia handle your leads for you, 24/7.
		</div>

		<div class="plans-wrapper">
			<div class="plan">
				<div class="title">Basic</div>
				<div class="price">
					<div class="currency" v-html="this.currency"></div>
					<div class="value">{{ this.prices[0] }}</div>
				</div>
				<div class="per-month">per month</div>

				<div class="features">
					<div class="item">Up to 50 leads per month</div>
					<div class="item">24/7 premium live support</div>
					<div class="item">Detailed lead insights</div>
					<div class="item">Custom scripts</div>
					<div class="item">Realtime messaging</div>
				</div>

				<Btn width="165" outline v-if="this.user.id != 0">Upgrade now</Btn>
				<Btn width="165" outline v-if="this.user.id == 0" :onclick="this.signup">Sign Up</Btn>
			</div>
			<div class="plan center">
				<div class="heading">Most popular</div>

				<div class="title">Standard</div>
				<div class="price">
					<div class="currency" v-html="this.currency"></div>
					<div class="value">{{ this.prices[1] }}</div>
				</div>
				<div class="per-month">per month</div>

				<div class="features">
					<div class="item">Up to 100 leads per month</div>
					<div class="item">24/7 premium live support</div>
					<div class="item">Detailed lead insights</div>
					<div class="item">Custom scripts</div>
					<div class="item">Realtime messaging</div>
				</div>

				<Btn width="165" outline v-if="this.user.id != 0">Upgrade now</Btn>
				<Btn width="165" outline v-if="this.user.id == 0" :onclick="this.signup">Sign Up</Btn>
			</div>
			<div class="plan premium">
				<div class="title">Premium</div>
				<div class="price">
					<div class="currency" v-html="this.currency"></div>
					<div class="value">{{ this.prices[2] }}</div>
				</div>
				<div class="per-month">per month</div>

				<div class="features">
					<div class="item">Up to 200 leads per month</div>
					<div class="item">24/7 premium live support</div>
					<div class="item">Detailed lead insights</div>
					<div class="item">Custom scripts</div>
					<div class="item">Realtime messaging</div>
				</div>

				<Btn width="165" outline v-if="this.user.id != 0">Upgrade now</Btn>
				<Btn width="165" outline v-if="this.user.id == 0" :onclick="this.signup">Sign Up</Btn>
			</div>
		</div>

		<div class="footer">
			If you receive more than 200 leads per month or would like to discuss a team plan, contact <a
						href="mailto:sales@getolivia.co">sales@getolivia.co</a>
		</div>
	</Page>
</template>

<script>
	import AuthResource from './../resources/Auth'

  import RegularPage from './../modules/pages/Regular.vue'

  import Button from './../modules/ui/Button.vue'

  export default {
    components: {
      Page: RegularPage,
      Btn: Button,
    },
    methods: {
      signup() {
        console.log('kek')
        this.$router.push('/auth/signup')
      }
    },
    data() {
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');
      fbq('track', 'PageView');

      AuthResource.getProfile().then(
        response => {
          this.user = response.body
        }
      )

			let locales = [navigator.language, ...navigator.languages].map((l) => l.toLowerCase())

			let prices = [199, 299, 499]
			let currency = '&#36;'

			if (locales.indexOf('en-gb') != -1) {
        currency = '&pound;';
        prices = [159, 249, 399]
			}

      return {
        currency,
				prices,
        user: {
          id: 0,
        }
      }
    }
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../assets/scss/_settings.scss';

	.page.plans {
		.wrapper {
			display: block;

			> .title {
				margin-top: 5px;
				font-size: 18px;
				font-weight: bold;
				text-align: center;
			}

			.subtitle {
				margin-top: 24px;
				font-size: 16px;
				text-align: center;
			}

			.plans-wrapper {
				display: flex;
				width: 900px;
				/*height: 600px;*/
				margin: 31px auto 0;
				padding-bottom: 25px;

				.plan {
					width: 33.333%;
					background-color: #ffffff;
					box-sizing: border-box;

					.title {
						margin-top: 34px;
						font-size: 18px;
						font-weight: 700;
						text-align: center;
						text-transform: uppercase;
					}

					.price {
						display: flex;
						justify-content: center;
						margin-top: 11px;
						text-align: center;

						.currency {
							margin: 6px 5px 0 0;
							font-size: 24px;
						}

						.value {
							font-size: 42px;
						}
					}

					.per-month {
						font-size: 14px;
						text-align: center;
						color: $placeholder;
					}

					.features {
						margin-top: 17px;
						padding: 0 40px;

						.item {
							width: 100%;
							height: 50px;
							line-height: 50px;
							font-size: 14px;
							text-align: center;
							border-bottom: 1px solid $border;
						}
					}

					button {
						display: block;
						margin: 26px auto 0;
					}

					&:not(.center) {
						/*height: 560px;*/
						margin-top: 31px;
						margin-bottom: 15px;
						padding-bottom: 30px;
						border: 1px solid $border;
						box-shadow: 0px 3px 6px rgba(#000000, 0.05);
						z-index: 1;

						&:first-child {
							border-right-width: 0;
						}

						&:last-child {
							border-left-width: 0;
						}

						.title {
							color: $placeholder;
						}

						.price {
							color: #74828D;
						}

						.per-month {

						}
					}

					&.center {
						/*height: 600px;*/
						padding-bottom: 45px;
						z-index: 2;
						border: 1px solid $border;
						box-shadow: 0px 6px 12px rgba(#000000, 0.08);

						.heading {
							width: 100%;
							height: 31px;
							line-height: 31px;
							color: #ffffff;
							font-size: 14px;
							font-weight: 700;
							background-color: $blue;
							text-transform: uppercase;
							text-align: center;
						}

						.title, .price {
							color: $blue;
						}
					}

					&.premium {
						.title, .price {
							color: #61D1D5;
						}
					}
				}
			}

			.footer {
				padding-bottom: 25px;
				font-size: 13px;
				text-align: center;
			}
		}
	}
</style>