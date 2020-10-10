
const TerserPlugin = require("terser-webpack-plugin");
const webpack = require('webpack')
module.exports = {
    publicPath: "/admin", // 基本路径
    outputDir: "www", // 输出文件目录
    lintOnSave: true, // eslint-loader 是否在保存的时候检查
    assetsDir: "static", //放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录。
    pages: {
        index: {
            entry: 'src/main.js',
            template: 'public/index.html',
        },
        meetRoom: {
            entry: 'src/views/MeetRoom/main.js',
            template: 'public/meetRoom.html',
        },
    }, // 以多页模式构建应用程序。
    configureWebpack: {
        plugins: [
            new webpack.ProvidePlugin({
                $: "jquery",
                jQuery: "jquery",
                "windows.jQuery": "jquery"
            })
        ],
        optimization: {
            minimizer: [
                new TerserPlugin({
                    terserOptions: {
                        compress: {
                            pure_funcs: ["console.log"]
                        }
                    }
                })
            ]
        }
    },
    devServer: {
        // host: 'localhost',
        host: "0.0.0.0",
        port: 8080, // 端口号
        https: false, // https:{type:Boolean}
        open: true, //配置自动启动浏览器  http://172.11.11.22:8888/rest/XX/
        hotOnly: true, // 热更新
        // proxy: 'http://localhost:8000'   // 配置跨域处理,只有一个代理
        proxy: {
            //配置自动启动浏览器
            // "/*": {
            //     target: "http://192.168.99.150:8004",
            //     changeOrigin: true,
            //     // ws: true,//websocket支持
            //     secure: false
            // },

            "/*": {
                target: "https://sfu.easyrtc.cn",
                changeOrigin: true,
                // ws: true,//websocket支持
                secure: false
            },

            // "/*": {
            //     target: "http://192.168.99.142:8004",
            //     changeOrigin: true,
            //     // ws: true,//websocket支持
            //     secure: false
            // },
        }
    },
}