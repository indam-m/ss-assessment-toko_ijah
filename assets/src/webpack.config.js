const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const pug = {
    test: /\.pug$/,
    use: ['html-loader?attrs=false', 'pug-html-loader']
};

const config = {
    entry: './src/js/app.js',
    output: {
        path: path.resolve(__dirname, '..'),
        filename: '[name].js'
    },
    module: {
        rules: [pug]
    },
    plugins: [
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: 'src/Index.pug',
            inject: false
        }),
        new HtmlWebpackPlugin({
            filename: 'item-amount.html',
            template: 'src/ItemAmount.pug',
            inject: false
        }),
        new HtmlWebpackPlugin({
            filename: 'item-in.html',
            template: 'src/ItemIn.pug',
            inject: false
        }),
        new HtmlWebpackPlugin({
            filename: 'item-out.html',
            template: 'src/ItemOut.pug',
            inject: false
        })
    ]
};
module.exports = config;