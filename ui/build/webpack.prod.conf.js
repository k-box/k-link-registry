var config = require("./config");
// merge deals with merging two configuration objects
var merge = require("webpack-merge");
var baseWebpackConfig = require("./webpack.base.conf");
var utils = require("./utils");
var webpack = require("webpack");
var HtmlWebpackPlugin = require("html-webpack-plugin");
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var OptimizeCSSPlugin = require("optimize-css-assets-webpack-plugin");
var UglifyJS = require("uglify-js");
var path = require("path");

var env = config.build.env;

// webpackConfig is the production webpack config,
// which is based on the base webpack config
var webpackConfig = merge(baseWebpackConfig, {
  module: {
    // loaders for stylesheets (standalone and embedded in vue bundles)
    rules: utils.styleLoaders({
      sourceMap: config.build.productionSourceMap,
      extract: true
    })
  },
  // should we point the dev tools to the sourcemap?
  devtool: config.build.productionSourceMap ? "#source-map" : false,
  // generate files with hashes, so they dont get cached in a wrong version
  output: {
    path: config.assetsRoot,
    filename: utils.assetsPath("js/[name].[chunkhash].js"),
    chunkFilename: utils.assetsPath("js/[id].[chunkhash].js")
  },
  plugins: [
    // enable further optimizations
    // http://vuejs.github.io/vue-loader/en/workflow/production.html
    new webpack.DefinePlugin({
      "process.env": env
    }),
    // minimize scripts
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      },
      sourceMap: true
    }),
    // extract styles into their own file
    new ExtractTextPlugin({
      filename: utils.assetsPath("css/[name].[contenthash].css")
    }),
    // Compress and dedupe css rules
    new OptimizeCSSPlugin({
      cssProcessorOptions: {
        safe: true
      }
    }),
    // populate our index.html with the correct hashes, minify it further
    new HtmlWebpackPlugin({
      filename: config.index,
      template: "./index.html",
      inject: true,
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true,
        minifyCSS: true
        // more options:
        // https://github.com/kangax/html-minifier#options-quick-reference
      },
      // load chunks in order of dependency
      chunksSortMode: "dependency"
    }),
    // split vendor js into its own file, seperate from application code
    new webpack.optimize.CommonsChunkPlugin({
      name: "vendor",
      minChunks: function (module, count) {
        // any required modules inside node_modules are extracted to vendor
        return (
          module.resource &&
          /\.js$/.test(module.resource) &&
          module.resource.indexOf(
            path.join(__dirname, "../../node_modules")
          ) === 0
        );
      }
    }),
    // extract webpack runtime and module manifest to its own file in order to
    // prevent vendor hash from being updated everytime app bundle is updated
    new webpack.optimize.CommonsChunkPlugin({
      name: "manifest",
      chunks: ["vendor"]
    })
  ]
});

if (config.build.bundleAnalyzerReport) {
  var BundleAnalyzerPlugin = require("webpack-bundle-analyzer")
    .BundleAnalyzerPlugin;
  webpackConfig.plugins.push(new BundleAnalyzerPlugin());
}

module.exports = webpackConfig;