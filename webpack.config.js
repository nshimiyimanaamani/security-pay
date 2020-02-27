module.exports = {
  resolve: {
    alias: {
      // Alias for using source of BootstrapVue
      "bootstrap-vue$": "bootstrap-vue/src/index.js"
    }
  },
  optimization: { minimize: true, mangleWasmImports: true },
  module: {
    rules: [
      {
        test: /\.js$/,
        // Exclude transpiling `node_modules`, except `bootstrap-vue/src`
        exclude: /node_modules\/(?!bootstrap-vue\/src\/)/,
        use: {
          loader: "babel-loader",
          options: {
            presets: ["env"]
          }
        }
      }
    ]
  }
};
if (process.env.NODE_ENV === "production") {
  module.exports.plugins = (module.exports.plugins || []).concat([
    new webpack.DefinePlugin({
      "process.env": {
        NODE_ENV: '"production"'
      }
    }),
    new webpack.optimize.UglifyJsPlugin()
  ]);
}
new webpack.optimize.CommonsChunkPlugin({
  name: "vendor",
  minChunks: function(module) {
    return module.context && module.context.indexOf("node_modules") !== -1;
  }
});
