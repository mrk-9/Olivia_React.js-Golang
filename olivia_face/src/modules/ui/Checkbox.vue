<template>
    <div class="checkbox">
        <input v-bind:id="id" @change="change" type="checkbox">
        <label v-bind:for="id">{{ title }}</label>
    </div>
</template>

<script>
    export default {
        props: {
            id: {
                type: String,
                default: `cb${Math.random()}`,
            },
            title: String,
        },
        methods: {
            change: function (e) {
                this.$emit('input', e.target.checked)
            },
        },
        data () {
            return {}
        }
    }
</script>

<style lang="scss">
    @import './../../assets/scss/_settings.scss';

    .checkbox {
        input[type=checkbox] {
            display: none;
        }

        input[type=checkbox] {
            & + label {
                position: relative;
                font-size: 13px;
                padding-left: 21px;
                cursor: pointer;

                &:before {
                     display: block;
                     position: absolute;
                     width: 14px; height: 14px;
                     margin-top: 5px;
                     border-radius: 1px;
                     border: 1px solid $blue;
                     box-sizing: border-box;
                     color: transparent;

                     content: 'check';
                     font-family: 'Material Icons';
                     -webkit-font-feature-settings: 'liga';
                     -webkit-font-smoothing: antialiased;

                     transition: background-color .1s, color .1s;
                }
            }

            &:checked + label {
                &:before {
                    background-color: $blue;
                    color: #ffffff;
                }
            }
        }
    }
</style>