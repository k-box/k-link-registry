var path = require("path");
var config = require("./config");
var utils = require("./utils");
var vueLoaderConfig = require("./vue-loader.conf");

// resolve is a helper function to make file paths relative
function resolve(dir) {
  return path.join(__dirname, "..", dir);
}

var isWatchRequested = process.env.WATCH ? true : false;

module.exports = {
  watch: isWatchRequested,
  // entrypoint for our app
  entry: {
    app: "./src/main.js"
  },
  // default output settings
  // (used for index.html)
  output: {
    path: config.assetsRoot,
    publicPath: path.join("{{ .BaseURL }}", "static/"),
    filename: "[name].js"
  },
  resolve: {
    // these file types will be importable directly
    extensions: [".js", ".vue", ".json"],
    alias: {
      // *vue files will be processed by the vue compiler
      vue$: "vue/dist/vue.esm.js",
      // @-prefixed paths will be relative to src/
      "@": resolve("src")
    }
  },
  module: {
    rules: [{
        // yaml files will be converted into a json
        test: /\.(yml|yaml)$/,
        loader: 'yml-loader'
      },
      {
        // js and vue files will pass eslint
        test: /\.(js|vue)$/,
        loader: "eslint-loader",
        enforce: "pre",
        include: [resolve("src"), resolve("test")],
        options: {
          formatter: require("eslint-friendly-formatter")
        }
      },
      {
        // vue files will be processed by the vue loader
        test: /\.vue$/,
        loader: "vue-loader",
        options: vueLoaderConfig
      },
      {
        // js files will be processed with babel, to transpile ES6 down
        // to regular JS.
        test: /\.js$/,
        loader: "babel-loader",
        include: [resolve("src"), resolve("test")]
      },
      {
        // images will be inlined (base64), if they are below 10k in size.
        // Otherwise they will be stored with the rest of the assets.
        test: /\.(png|jpe?g|gif|svg)(\?.*)?$/,
        loader: "url-loader",
        options: {
          limit: 10000,
          name: utils.assetsPath("img/[name].[hash:6].[ext]")
        }
      },
      {
        // web fonts will always be inlined (in the CSS file)
        test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
        loader: "url-loader",
        options: {
          // limit: 10000,
          name: utils.assetsPath("fonts/[name].[hash:7].[ext]")
        }
      }
    ]
  }
};