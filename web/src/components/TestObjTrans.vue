<template>
  <div>
    {{ complexObj }}
  </div>
</template>

<script>
export default {
  name: "TestObjTrans",
  props: {
    complexObj: {
      type: Object,
      default: () => {},
    },
  },
};
</script>

<style scoped lang="less">
</style>
