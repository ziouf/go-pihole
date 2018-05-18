<template>
    <!-- small box -->
    <div class="small-box" :class="bgClass">
    <div class="inner">
        <h3 v-bind="value"></h3>

        <p v-bind="txt"></p>
    </div>
    <div class="icon">
        <i class="fa" :class="iconClass" />
    </div>
    <a :href="href" class="small-box-footer">
        More info <i class="fa fa-arrow-circle-right"></i>
    </a>
    </div>
</template>

<script>
export default {
    name: 'widget-smallbox',
    props: [ 'value', 'txt', 'bgClass', 'iconClass', 'href' ]
};
</script>

<style scoped>
</style>