module.exports = {
    devServer: {
        historyApiFallback: true,
        open: true,
        host: '0.0.0.0',
        port: 8838,
        // 代理(非拦截请求发送代理)
        proxy: {
            '/': {
                target: process.env.VUE_APP_PROXY_URL,
                changeOrigin: true,
            }
        }
    },
}