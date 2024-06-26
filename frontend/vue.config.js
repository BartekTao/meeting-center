const { defineConfig } = require('@vue/cli-service')
// module.exports = defineConfig({
//   transpileDependencies: true
// })

const webpack = require('webpack');

// Ref: https://stackoverflow.com/questions/77752897/feature-flag-vue-prod-hydration-mismatch-details-is-not-explicitly-defined
module.exports = defineConfig({
  configureWebpack: {
    plugins: [
      new webpack.DefinePlugin({
        // Vue CLI is in maintenance mode, and probably won't merge my PR to fix this in their tooling
        // https://github.com/vuejs/vue-cli/pull/7443
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'false',
      })
    ],
  },
  devServer: {
    hot: false,          // 禁用热模块替换
    liveReload: false,   // 禁用实时重新加载
    port: 8888
  }
});


// module.exports = {
//   devServer: {
//     port: 8888
//   }
// }
