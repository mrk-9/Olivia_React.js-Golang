<template>
	<div :class="{ textarea: true, nomargin: this.nomargin}">
		<textarea :class="{ error: error.length > 0 }" :placeholder="placeholder"
							:value="value" @input="change" :maxlength="maxlength"
							@blur="blur" @focus="focus"></textarea>
		<div :class="{ errors: true, visible: error.length > 0 }">{{ error }}</div>
	</div>
</template>

<script>
  export default {
    props: {
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
      blur: function (e) {
        this.$emit('blur')
      },
      focus: function (e) {
        this.$emit('focus')
      },
    },
    data () {
      return {}
    },
  }
</script>

<style lang="scss">
	@import './../../assets/scss/_settings.scss';

	.textarea {
		width: 100%;

		&:not(.nomargin) {
			margin-bottom: 18px;
		}

		textarea {
			width: 100%;
			height: 68px;
			padding: 8px 10px;
			font-size: 13px;
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