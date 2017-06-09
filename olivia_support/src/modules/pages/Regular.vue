<template>
	<div class="page regular" @click="contextMenuVisibility = false">
		<div class="header-wrap">
			<div class="header">
				<Logo/>

				<div class="user" v-if="this.user.id != 0">
					<CMenu top="60" :visible="contextMenuVisibility">
						<a href="#" @click="logout">Log out</a>
					</CMenu>

					<!--<div class="avatar" style="background-image: url('http://vignette1.wikia.nocookie.net/erbparodies/images/d/d7/Filepicker-DdnyfEIRyqs1xryV8AQG_tony_montana.jpg/revision/latest?cb=20150111071616');"></div>-->
					<div :class="{ name: true, active: currentRoute == '/account' }" @click.stop="toggleContextMenuVisibility">{{ user.fname + ' ' + user.sname }}</div>
				</div>
				<router-link class="signup" to="/auth/signup" v-if="this.user.id == 0">Sign Up</router-link>
			</div>
		</div>

		<div class="wrapper">
			<slot></slot>
		</div>
	</div>
</template>

<script>
  import Logo from './../Logo.vue'
  import CMenu from './../ui/ContextMenu.vue'

  import AuthResource from './../../resources/Auth'

  export default {
    components: {
      Logo,
      CMenu,
    },
    methods: {
      toggleContextMenuVisibility() {
        this.contextMenuVisibility = !this.contextMenuVisibility;
      },
			upgrade() {
        this.$router.push('/plans')
			},
      logout() {
        AuthResource.logout.bind(AuthResource, this.$router)()
      },
    },
    data () {
      AuthResource.getProfile().then(
        response => {
          this.user = response.body;
				},
      )

      return {
        currentRoute: this.$router.currentRoute.path,
        menuItems: {
          '/home': 'Prospects',
          '/scripts': 'Scripts & Templates',
        },
        contextMenuVisibility: false,
        user: {
          id: 0,
          fname: '...',
					sname: '...',
				},
      }
    }
  }
</script>

<style lang="scss" rel="stylesheet/scss">
	@import './../../assets/scss/_settings.scss';

	$width: 1376px;
	$headerHeight: 74px;

	.regular {
		justify-content: flex-start !important;
		flex-direction: column;


		.header-wrap {
			width: 100%;
			height: $headerHeight;
			margin-bottom: 32px;
			background-color: #ffffff;
			@include shadow();

			.header {
				display: flex;
				align-items: center;
				justify-content: space-between;
				max-width: 1376px;
				min-width: 1024px;
				width: 100%;
				height: 100%;
				margin: 0 auto 0;
				padding: 0 20px;
				box-sizing: border-box;

				.signup {
					height: 34px;
					margin-left: 20px;
					padding: 0 20px;
					line-height: 34px;
					background-color: $blue;
					color: #ffffff;
					border-radius: 17px;
					border: none;
					cursor: pointer;
				}

				.menu {
					height: $headerHeight;

					a.item:hover {
						text-decoration: none;
						color: $blue;
					}

					.item {
						display: inline-block;
						height: $headerHeight;
						padding: 0 60px;
						font-size: 14px;
						font-weight: 600;
						line-height: $headerHeight;
						color: $text;
						border-bottom: 4px solid transparent;
						box-sizing: border-box;
						transition: color .1s, border-bottom-color .1s;

						&.active {
							color: $blue;
							border-bottom-color: $blue;
						}

						.big-title, .small-title {
							height: 24px;
							line-height: 24px;
						}

						.small-title {
							margin-left: 10px;
							color: rgba($text, .53);
						}

						button {
							-webkit-appearance: none;
							-moz-appearance: none;
							appearance: none;
							height: 24px;
							margin-left: 20px;
							padding: 0 10px;
							background-color: $blue;
							color: #ffffff;
							border-radius: 12px;
							border: none;
							cursor: pointer;
						}
					}
				}

				.user {
					display: flex;
					justify-content: flex-end;
					align-items: center;
					position: relative;
					width: 200px;

					.avatar, .name {
						display: inline-block;
					}

					.avatar {
						width: 42px;
						height: 42px;
						margin-right: 10px;
						border-radius: 50%;
						background: no-repeat center;
						background-size: cover;
					}

					.name {
						position: relative;
						padding-right: 20px;
						font-size: 14px;
						font-weight: 700;
						cursor: pointer;

						&.active {
							color: $blue;
						}

						&:after {
							display: block;
							position: absolute;
							top: 2px;
							right: 0;
							height: 19px;
							margin-left: 6px;
							line-height: 19px;
							font-size: 18px;
							color: $placeholder;

							content: 'keyboard_arrow_down';
							font-family: 'Material Icons';
							-webkit-font-feature-settings: 'liga';
							-webkit-font-smoothing: antialiased;
						}
					}

					.cmenu {
						right: 0;
						margin-top: -25px;
					}
				}
			}
		}

		.wrapper {
			max-width: $width;
			min-width: 1024px;
			width: 100%;
			height: calc(100% - 20px - 106px);
			margin: 0 auto 0;
			padding: 0 20px;
			box-sizing: border-box;
		}
	}

	@media (max-width: 1310px) {
		.regular {
			.header-wrap {
				.header {
					.menu {
						.item {
							padding: 0 30px;
						}
					}
				}
			}
		}
	}

	@media (max-width: 1130px) {
		.regular {
			.header-wrap {
				.header {
					.menu {
						.item {
							.small-title {
								display: none;
							}
						}
					}
				}
			}
		}
	}
</style>
