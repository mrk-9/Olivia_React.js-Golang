<template>
	<div :class="{ input: true, nomargin: this.nomargin}">
		<input :class="{ error: error.length > 0 }" :type="type" :placeholder="placeholder"
					 :value="value" @input="change" :maxlength="maxlength">
		<div :class="{ errors: true, visible: error.length > 0 }">{{ error }}</div>
	</div>
</template>

<script>
  export default {
    props: {
      type: {
        type: String,
        default: 'text',
      },
      placeholder: {
        type: String,
        default: 'Placeholder',
      },
      error: {
        type: String,
        default: '',
      },
      nomargin: {
        type: Boolean,
				default: false,
			},
      maxlength: Number,
      value: String,
    },
    methods: {
      change: function (e) {
        this.$emit('input', e.target.value)
      },
    },
    data () {
      return {}
    },
  }
</script>

<style lang="scss">
	@import './../../assets/scss/_settings.scss';

	.input {
		width: 100%;

		&:not(.nomargin) {
			margin-bottom: 18px;
		}

		input {
			width: 100%;
			height: 34px;
			padding: 0px 10px;
			font-size: 13px;
			line-height: 34px;
			border: 1px solid $border;
			box-sizing: border-box;
			transition: border-color .2s;

			&.error {
				border-color: $red;
			}
		}

		.errors {
			position: absolute;
			color: $red;
			font-size: 10px;
			opacity: 0;
			transition: opacity .2s;

			&.visible {
				opacity: 1;
			}
		}
	}
</style>