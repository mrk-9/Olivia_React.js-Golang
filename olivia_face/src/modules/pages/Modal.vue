<template>
	<div class="page modal">
		<div class="content">
			<Logo/>

			<Card>
				<slot/>
			</Card>

			<div class="footer">
				&copy; 2017 Olivia
			</div>
		</div>
	</div>
</template>

<script>
  import Logo from './../Logo.vue'
  import Card from './../ui/Card.vue'

  import intercom from './../vendor/intercom'

  export default {
    components: {
      Logo,
      Card,
    },
		data () {
      intercom(null, null, null, true);
      ga('set', 'page', window.location.pathname);
      ga('send', 'pageview');

			return {}
    }
  }
</script>

<style lang="scss">
	@import './../../assets/scss/_settings.scss';

	.modal {
		height: auto !important;

		.content {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: space-between;
			width: 379px;
			min-height: 100vh;

			.logo {
				margin: 93px 0 50px;
			}

			> .footer {
				margin: 100px 0 40px;
				color: $placeholder;
				font-size: 13px;
			}
		}
	}
</style>