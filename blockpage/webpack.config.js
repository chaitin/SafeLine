var HtmlWebpackPlugin = require("html-webpack-plugin");
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var path = require("path");

module.exports = {
  entry: "./entry.js",
  output: {
    path: path.resolve(__dirname, "./dist"),
    publicPath: "",
    filename: "bundle.js",
  },
  module: {
    rules: [
      {
        test: /\.less$/,
        loader: ExtractTextPlugin.extract({
          fallback: "style-loader",
          use: [
            { loader: "css-loader" },
            { loader: "postcss-loader" },
            { loader: "less-loader" },
          ],
        }),
      },
      { test: /\.jade$/, loader: "pug-loader" },
      { test: /\.png|ico|svg$/, loader: "url-loader?limit=99999" },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      inject: false,
      cache: false,
      template: "index.jade",
      filename: "index.html",
      favicon: "safeline.png",
      title: "网站被拦截，哈哈恍恍惚惚恍恍惚惚恍恍惚惚恍恍惚惚",
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true,
        minifyJS: true,
        minifyCSS: true,
        // more options:
        // https://github.com/kangax/html-minifier#options-quick-reference
      },
    }),
    new ExtractTextPlugin("styles.css"),
  ],
};
