<template>
	<div :class="{ card: true, nopadding: this.nopadding }">
		<slot/>
	</div>
</template>

<script>
  export default {
    props: {
      nopadding: Boolean,
    },
    data () {
      return {}
    }
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../../assets/scss/_settings.scss';

	.card {
		position: relative;
		width: 100%;
		min-height: 100px;
		padding: 40px;
		background-color: #ffffff;
		box-sizing: border-box;
		@include shadow();

		&.nopadding {
			padding: 0;
		}

		.divider {
			width: 100%;
			height: 0;
			margin: 20px 0;
			border-bottom: 1px solid $border;
		}
	}
</style>