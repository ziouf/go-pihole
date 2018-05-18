<template>
    <div class="box box-default" :class="isCollapsed()">
        <div class="box-header with-border">
            <h3 class="box-title" v-bind="title" />

            <div class="box-tools pull-right">
            <button type="button" class="btn btn-box-tool" data-widget="collapse">
                <i class="fa fa-plus"></i>
            </button>
            </div>
            <!-- /.box-tools -->
        </div>
        <!-- /.box-header -->
        <div class="box-body" v-bind="txt" />
        <!-- /.box-body -->
    </div>
    <!-- /.box -->
</template>

<script>
    export default {
        name: 'widget-box',
        props: [ 'title', 'txt', 'collapsed' ],
        data: function() {
            return {
                
            }
        },
        methods: {
            isCollapsed: function() {
                return { 'collapsed-box': collapsed == "" ? true : false }
            },
        },
    }
</script>

<style scoped>

</style>