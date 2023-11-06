<template>
  <div id="app" class="move_scroll_style">
    <router-view />
  </div>
</template>

<style lang="less">
body {
  padding: 0;
  margin: 0;
  #app {
    width: 100vw;
    height: 100vh;
    background-image: url(./assets/index_bg_icon.png);
    background-size: 30% 30%;
    background-repeat: no-repeat;
    background-position: 100% 100%;
  }
}

// 去除滚动条样式类（全局-垂直方向）
.move_scroll_style {
  // 去除滚动样式条
  scrollbar-width: none;
  /* firefox */
  -ms-overflow-style: none;
  /* IE 10+ */
  overflow-x: hidden;
  overflow-y: auto;
  &::-webkit-scrollbar {
    display: none;
    /* Chrome Safari */
  }
}
</style>
