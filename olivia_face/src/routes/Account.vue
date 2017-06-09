<template>
	<Page class="account">
		<div class="content">
			<div class="title">Account Information</div>

			<Card>
				<div class="pairs">
					<div class="pair">
						<div class="key">Name</div>
						<div class="value">{{ this.user.fname + ' ' + this.user.sname }}</div>
					</div>

					<div class="pair">
						<div class="key">Email</div>
						<div class="value">{{ this.user.email }}</div>
					</div>

					<div class="pair">
						<div class="key">Phone</div>
						<div class="value">{{ this.user.phone }}</div>
					</div>

					<div class="pair">
						<div class="key">Olivia email</div>
						<div class="value">{{ this.user.o_email + '@' + this.host }}</div>
					</div>

					<div class="pair">
						<div class="key">Olivia SMS number</div>
						<div class="value">{{ this.user.assigned_phone_number }}</div>
					</div>

					<div class="pair">
						<div class="key">Timezone</div>
						<div class="value">{{ this.timezone }}</div>
					</div>

					<div class="pair">
						<div class="key">Plan</div>
						<div class="value">Free trial <span class="red">(10 days remaining)</span></div>
					</div>
				</div>

				<div class="button-wrapper">
					<Btn width="155" outline>Upgrade now</Btn>
				</div>
			</Card>
		</div>
	</Page>
</template>

<script>
  import AuthResource from './../resources/Auth.js'

  import RegularPage from './../modules/pages/Regular.vue'

  import Card from './../modules/ui/Card.vue'
	import Button from './../modules/ui/Button.vue'

  export default {
    components: {
      Page: RegularPage,
			Card,
			Btn: Button,
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

      return {
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
				host: 'oliviaemail.com',
				user: {},
			}
		}
	}
</script>

<style lang="scss" rel="stylesheet/scss">
	@import "./../assets/scss/_settings.scss";

	.page.account {
		.content {
			max-width: 660px;
			margin: 0 auto 0;

			.red {
				color: $moreRed;
			}

			.title {
				margin: 20px 0 44px;
				font-size: 18px;
				font-weight: 700;
				text-align: center;
			}

			.card {
				.pair {
					display: flex;
					padding: 17px 0;
					font-size: 14px;
					border-bottom: 1px solid $border;

					&:nth-last-child(1) {
						border-bottom: none;
					}

					.key {
						width: 150px;
						font-weight: 700;
					}

					.value {
						width: calc(100% - 150px);
					}
				}

				.button-wrapper {
					padding-left: 150px;
				}
			}
		}
	}
</style>