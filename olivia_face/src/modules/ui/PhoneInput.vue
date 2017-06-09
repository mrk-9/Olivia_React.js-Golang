<template>
	<div :class="{ input: true, phone: true, nomargin: this.nomargin}">
		<div class="phone-wrapper">
			<input :class="{ code: true, error: error.length > 0 }" type="text"
						 :value="phone.code" @input="changeCode" maxlength="3" placeholder="+1">
			<input :class="{ number: true, error: error.length > 0 }" :type="type"
						 :value="phone.number" @input="changeNumber" :maxlength="maxlength" :placeholder="placeholder">
		</div>
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
      value: Object,
    },
    methods: {
      changeCode: function (e) {
        this.phone.code = e.target.value;
        console.log(this.phone)
        this.$emit('input', this.phone)
      },
			changeNumber: function (e) {
        this.phone.number = e.target.value;
        console.log(this.phone)
        this.$emit('input', this.phone)
      },
    },
    data () {
      return {
        phone: {
          code: '',
          phone: '',
        },
			}
    },
  }
</script>

<style lang="scss">
	@import './../../assets/scss/_settings.scss';

	.input.phone {
		width: 100%;

		&:not(.nomargin) {
			margin-bottom: 18px;
		}

		.phone-wrapper {
			$codeWidth: 45px;

			input.code {
				width: $codeWidth;
			}

			input.number {
				width: calc(100% - #{$codeWidth} - 5px);
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