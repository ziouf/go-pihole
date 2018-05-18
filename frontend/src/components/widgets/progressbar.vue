<template>
    <div>
        <div class="progress">
            <div class="progress-bar" v-bind:style="getWidth()" />
        </div>
        <span class="progress-description" v-bind="description" />
    </div>
</template>

<script>
    export default {
        name: 'widget-progressbar', 
        props: [ 'value', 'description' ],
        methods: {
            getWidth: function() {
                return { 
                    "width": value == "" ? "0%" : value
                }
            }
        }
    }
</script>

<style scoped>

</style>