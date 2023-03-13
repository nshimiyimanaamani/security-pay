module.exports = {
  productionSourceMap: false,
  configureWebpack: {
    optimization: {
      splitChunks: {
        cacheGroups: {
          vendor: {
            chunks: 'all',
            maxSize: 250000,
            maxAsyncRequests: 250000,
            maxInitialRequests: 250000,
            reuseExistingChunk: true,
            enforce: true,
            priority: 1
          }
        }
      }
    }
  },
  chainWebpack: config => {
    config.module.rules.delete('eslint')
  }
}
