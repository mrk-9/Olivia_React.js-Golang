<template>
    <div class="loader">
        <scale-loader :loading="loading" color="#39c3dd" width="3px" height="20px"></scale-loader>
    </div>
</template>

<script>
    import ScaleLoader from 'vue-spinner/src/ScaleLoader.vue'

    export default {
        components: {
            ScaleLoader
        },
        props: {
            loading: {
                type: Boolean,
                default: false,
            }
        },
        data () {
            return {}
        },
    }
</script>

<style lang="scss">
    .loader {
        margin: 8px 0 0;
    }
</style>