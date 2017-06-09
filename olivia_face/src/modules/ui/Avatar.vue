<template>
	<div :class="{ avatar: true, big: big }" v-html="(!this.user.image) ? convertNameToInitials(this.user) : ''" :style="{ 'background-image': 'url(' + this.user.image + ')' }"></div>
</template>

<script>
	export default {
		props: {
			user: Object,
			big: {
			  type: Boolean,
				default: false,
			}
		},
		methods: {
		  convertNameToInitials(user) {
        let initials = '&ndash;'

        if (user) {
          if (user.name) {
            initials = user.name.split(' ', 2).map((word) => word[0]).join('').toUpperCase()
          } else if (user.fname && user.sname) {
            initials = user.fname[0] + user.sname[0]
            initials.toUpperCase()
          }
        }

        return initials
			}
		},
		data () {
		  return {}
		}
	}
</script>

<style type="text/scss" lang="scss">
	@import './../../assets/scss/_settings.scss';

	.avatar {
		width: 42px; height: 42px;
		line-height: 42px;
		text-align: center;
		border-radius: 50%;
		border: 1px solid #ffffff;
		box-sizing: border-box;
		background: $blue no-repeat center;
		background-size: cover;
		color: #ffffff;
		font-size: 16px;

		&.big {
			width: 82px; height: 82px;
			line-height: 82px;
			font-size: 32px;
		}
	}
</style>