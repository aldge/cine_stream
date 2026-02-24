const path = require('path');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = (env, argv) => {
  const isProduction = argv.mode === 'production';
  return {
    entry: './src/cine-player.js',
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: isProduction ? 'cine-player.min.js' : 'cine-player.js',
      library: 'CinePlayer',
      libraryTarget: 'umd',
      globalObject: 'this'
    },
    mode: argv.mode || 'development',
    optimization: {
      minimize: isProduction,
      minimizer: isProduction
        ? [
            new TerserPlugin({
              terserOptions: {
                compress: {
                  drop_console: true,
                  drop_debugger: true
                },
                mangle: true,
                output: {
                  comments: false
                }
              }
            })
          ]
        : []
    },
    devServer: {
      static: {
        directory: path.join(__dirname, 'public')
      },
      port: 8080,
      hot: true
    }
  };
};