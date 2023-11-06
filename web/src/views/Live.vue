<template>
  <div>
    <!-- 左侧栏 -->
    <div class="left-menus">
      <div
        class="left-menu"
        v-for="(item, index) in list"
        :key="index"
        @click="onChangeTab(index)"
      >
        <AdaptiveWidthScale :maxWidth="90" :isShow="true">
          {{ item.title }}
        </AdaptiveWidthScale>
      </div>
    </div>
    <!-- 右侧主体内容 -->
    <div class="live_page_more">
      <div>1231323123</div>
      <div>1231323123</div>
      <div>1231323123</div>
      <div>1231323123</div>
      <div v-for="(index, item) in 10" :key="index">
        <AdaptiveWidthScale :maxWidth="90" :isShow="true">
          加大渲染计算
        </AdaptiveWidthScale>
      </div>
    </div>
  </div>
</template>

<script>
import AdaptiveWidthScale from "../components/AdaptiveWidthScale.vue";

export default {
  components: {
    AdaptiveWidthScale: AdaptiveWidthScale,
  },
  data() {
    return {
      list: [
        {
          title: "最新战报sdfsdfs sdf sdf ",
        },
        {
          title: "每日战况",
        },
      ],
      activeMenuText: "最新战报",
    };
  },
  mounted() {
    for(let i = 0; i < 5000; i++) {
      console.log(1);
    }
  },
  methods: {
    // 点击切换左侧栏
    onChangeTab(index) {
      console.log(index);
    },
  },
};
</script>

<style lang="less" scoped>
// 左侧栏
.left-menus {
  position: fixed;
  top: 0;
  left: 0;
  width: 120px;
  height: 100%;
  background: rgba(0, 0, 0, 0.2);

  .left-menu {
    display: block;
    height: 43px;
    margin-bottom: 3px;
    cursor: pointer;
    position: relative;
    line-height: 43px;
    font-weight: 700;
    font-size: 13px;
    display: flex;
    align-items: center;
    justify-content: center;

    .i18n_style_outer {
      margin: 0 auto;

      :deep(.i18n_style_inner) {
        line-height: 43px;
        font-weight: 700;
        font-size: 13px;
        color: blue;
        text-transform: uppercase;
      }
    }
  }
}

.live_page_more {
  position: fixed;
  top: 0;
  left: 120px;
  width: 100% - 120px;
  height: 100%;
  overflow: scroll;
}
</style>