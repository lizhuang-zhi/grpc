<template>
  <!-- 存在一个隐患！！就是样式的穿透问题 -->
  <div class="i18n_style_outer" :style="{ width: outerWidthVW + 'vw' }">
    <!-- 在同一个组件中设置多个ref相同，只有最后一个ref起效果 -->
    <span class="i18n_style_inner" ref="span_text">
      <slot></slot>
      {{ text }}
    </span>
  </div>
</template>

<script>
export default {
  name: "AdaptiveWidthScale",
  props: {
    // 最大宽度
    maxWidth: Number,
    // 对于一开始为display:none;的使用位置，当v-show为true时进行重新渲染
    isShow: {
      type: Boolean,
      default: true,
    },
    text: String,
  },
  data() {
    return {
      outerWidth: 0, // 外部宽度(px)
      outerWidthVW: 0, // 外部宽度(vw)
      clientWidth: 0,
    };
  },
  watch: {
    // 监听组件显示，重渲染
    isShow: function (val) {
      if (val) {
        this.$nextTick(() => {
          this.init();
        });
      }
    },
    text: function () {
      this.$nextTick(() => {
        this.init();
      });
    },
  },
  mounted() {
    // 获取屏幕宽度，并兼容刘海屏
    this.clientWidth =
      document.documentElement.clientWidth ||
      window.innerWidth ||
      document.body.clientWidth;

    // 判断是否为默认机型（iphone678plus)
    if (this.clientWidth === 736) {
      this.outerWidth = this.maxWidth;
    } else {
      // 非默认机型则对传入的背景宽度进行重新计算
      this.outerWidth = (this.maxWidth / 736) * this.clientWidth;
    }

    this.outerWidthVW = (this.outerWidth / this.clientWidth) * 100;
    // 对于一开始为display:none;的使用位置，当v-show为true时进行重新渲染
    if (this.isShow) {
      this.$nextTick(() => {
        this.init();
      });
    }
  },
  methods: {
    init() {
      // 获取文本宽度
      let currentTextWidth = this.$refs["span_text"].clientWidth;
      if (currentTextWidth == 0) {
        return;
      }
      if (this.outerWidth <= currentTextWidth) {
        // 获取实时伸缩比例
        let proportion = this.outerWidth / currentTextWidth;
        this.$refs["span_text"].style.transform = `scale(${proportion})`;
        this.$refs["span_text"].style.transformOrigin = `center left`;
        this.updateOuterWidth(this.outerWidth);
      } else {
        this.updateOuterWidth(currentTextWidth);
      }
    },
    // 修改外部宽度
    updateOuterWidth(currentTextWidth) {
      this.outerWidth = currentTextWidth;
      this.outerWidthVW = (this.outerWidth / this.clientWidth) * 100;
    },
  },
};
</script>

<style lang="less" scoped>
.i18n_style_outer {
  text-align: center;
  .i18n_style_inner {
    display: inline-block;
    white-space: nowrap;
  }
}
</style>