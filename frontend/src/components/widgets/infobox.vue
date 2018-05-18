<template>
    <div class="info-box">
    <span class="info-box-icon bg-aqua"><i class="fa" :class="iconClass" /></span>

    <div class="info-box-content">
        <span class="info-box-text" v-bind="txt"></span>
        <span class="info-box-number" v-bind="value"></span>
    </div>
    <!-- /.info-box-content -->
    </div>
    <!-- /.info-box -->
</template>

<script>
export default {
  name: "widget-mail-infobox",
  props: [ 'txt', 'value', 'iconClass' ],
};
</script>

<style scoped>
</style>