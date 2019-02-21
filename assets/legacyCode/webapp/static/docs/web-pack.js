const path = require('path');

module.exports = {
  entry: './index.css',
  entry: './d3.js',
  entry: './jsldvizapp.js',
  output: {
      filename: 'bundle.js',
      path: path.resolve(__dirname, 'dist')
  }
};
