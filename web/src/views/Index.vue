<template>
  <div class="container">
    <div>第一个server接口:</div>
    <div>get /query</div>
    <div>param: user</div>
    <div>手动设置了休眠3s</div>
    <div>返回数据: {{ queryMsg }}</div>
    <div>返回数据: {{ queryData }}</div>

    <br />
    <div>第二个server接口:</div>
    <div>get /list</div>
    <div>param: user</div>
    <div>没有手动设置休眠</div>
    <div>返回数据: {{ list }}</div>

    <br />
    <div>第三个server接口:</div>
    <div>post /form</div>
    <div>
      var formReqParam struct { Name string `json:"name"` Age int `json:"age"`
      Hobby string `json:"hobby"` }
    </div>
    <button @click="formReq">发送请求</button>
    <div>返回数据: {{ formMsg }}</div>
  </div>
</template>  
  
<script>
import { form, query, list } from "@/api/request.js";
export default {
  data() {
    return {
      queryMsg: "", // 第一个接口返回数据
      queryData: "", // 第一个接口返回数据

      list: [], // 第二个接口返回数据

      formMsg: null, // 第三个接口返回数据
    };
  },
  mounted() {
    this.queryInfo();
  },
  methods: {
    queryInfo() {
      query({
        user: "morning",
      })
        .then((res) => {
          console.log(res);
          this.queryMsg = res.message;
          this.queryData = res.data;
        })
        .catch((err) => {});

      list({
        id: 123,
      })
        .then((res) => {
          console.log(res);
          this.list = res.data;
        })
        .catch((err) => {});
    },
    formReq() {
      form({
        name: "leo",
        age: 22,
        hobby: "ball",
      })
        .then((res) => {
          console.log(res);
          this.formMsg = res;
          const data = JSON.stringify({
            title: "成功获取query请求的响应值",
            msg: "数据上报",
            data: res,
          });
          navigator.sendBeacon(
            "http://124.221.157.89:12222/report",
            new Blob([data], { type: "application/json" })
          );
        })
        .catch((err) => {});
    },
  },
};
</script>

<style lang="less" scoped>
.container {
  margin: 40px 50px;
}
</style>