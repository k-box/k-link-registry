var path = require("path");

module.exports = {
  // location of the index.html file
  index: path.resolve(__dirname, "../dist/index.html"),
  // location of all compiled files
  assetsRoot: path.resolve(__dirname, "../dist"),
  // subdirectory where we will move all static files
  assetsSubDirectory: "./",

  // build settings for production
  build: {
    env: {
      NODE_ENV: '"production"'
    },
    productionSourceMap: true,
    // Run the build command with an extra argument to
    // View the bundle analyzer report after build finishes:
    // `npm run build --report`
    // Set to `true` or `false` to always turn it on or off
    bundleAnalyzerReport: process.env.npm_config_report
  },
  dev: {
    env: {
      NODE_ENV: '"development"'
    },
    produceSourceMap: true
  }
};
