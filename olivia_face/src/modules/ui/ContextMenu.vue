<template>
    <div :class="{ cmenu: true, visible: visible }" :style="{ top: top + 'px', left: left + 'px' }">
        <slot></slot>
    </div>
</template>

<script>
    export default {
        props: {
            visible: {
                type: Boolean,
                default: true,
            },
            top: {
                type: String,
                default: 'auto',
            },
            left: {
                type: String,
                default: 'auto',
            },
        },
        data () {
            return {}
        },
    }
</script>

<style lang="scss" rel="stylesheet/scss">
    @import './../../assets/scss/_settings.scss';

    .cmenu {
        position: absolute;
        min-width: 172px;
        background-color: #ffffff;
        opacity: 0;
        pointer-events: none;
        transition: margin-top .3s, opacity .3s;
        z-index: 100;
        border: 1px solid $border;
        @include shadow-big();

        &.visible {
            margin-top: 0;
            opacity: 1;
            pointer-events: auto;
        }

        a {
            display: block;
            position: relative;
            width: 100%; height: 44px;
            padding: 0 19px;
            line-height: 44px;
            box-sizing: border-box;
            color: $text;
            transition: background-color .1s;

            &:before {
                display: block;
                position: absolute;
                width: calc(100% - 19px - 19px); height: 100%;
                content: '';
                margin: -1px 0 0 0;
                border-bottom: 1px solid $border;
            }

            &:hover {
                background-color: $border;
                text-decoration: none;
            }

            &:last-child:before {
                border-bottom: none;
            }
        }
    }
</style>