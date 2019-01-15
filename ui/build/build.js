// set production env, to trigger some further optimizations
process.env.NODE_ENV = "production";

var path = require("path");
// rimraf deals with removing files
var rm = require("rimraf");
var config = require("./config");
var webpack = require("webpack");
var webpackConfig = require("./webpack.prod.conf");

rm(path.join(config.assetsRoot, config.assetsSubDirectory), err => {
  if (err) throw err;
  webpack(webpackConfig, function(err, stats) {
    if (err) throw err;
    process.stdout.write(
      stats.toString({
        colors: true,
        modules: false,
        children: false,
        chunks: false,
        chunkModules: false
      }) + "\n\n\n"
    );

    if(webpackConfig.watch){
      console.info("Watching for source file changes...");
    }
  });
});
