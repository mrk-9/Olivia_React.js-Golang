<template>
	<button :class="{ buttonElement: true, outline: this.outline }" :style="style" @click="clickHandler">
		<slot></slot>
	</button>
</template>

<script>
  export default {
    props: {
      width: {
        type: String,
        default: '119',
      },
      onclick: {
        type: Function,
        default: () => {
        },
      },
      outline: {
        type: Boolean,
        default: false,
      },
    },
    data () {
      return {
        style: {
          width: this.width + 'px',
        },
        clickHandler: this.onclick,
      }
    },
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../../assets/scss/_settings.scss';

	.buttonElement {
		height: 34px;
		color: #ffffff;
		font-size: 15px;
		font-weight: 700;
		border: 2px solid $blue;
		border-radius: 19px;
		box-sizing: border-box;
		cursor: pointer;
		background-color: $blue;
		transition: background-color .2s;

		&:hover {
			background-color: lighten($blue, 5%);
		}

		&.outline {
			background-color: transparent;
			color: $blue;

			&:hover {
				background-color: lighten($blue, 35%);
				/*color: #ffffff;*/
			}
		}
	}
</style>