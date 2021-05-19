const path = require("path");
// const webpack = require('webpack');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { merge } = require("webpack-merge");

const common = {

	resolve: {
		extensions: [".ts", ".tsx", ".scss", ".css", ".js"]
	},

	entry: [`./src/index.tsx`],

	plugins: [
		new MiniCssExtractPlugin({filename: `index.css`}),
		// new webpack.HotModuleReplacementPlugin(),
	],

	module: {
		rules: [
			{
				test: /\.ts(x?)$/,
				exclude: /node_modules/,
				use: "ts-loader"
			},
			// {
			// 	test: /\.s[ac]ss$/i,
			// 	use: ["style-loader", "css-loader", "sass-loader"],
			// 	exclude: /node_modules/
			// },
			{
				test: /\.(s?)css$/,
				use: [MiniCssExtractPlugin.loader, "css-loader", "sass-loader"],
				// exclude: /node_modules/
			},
			// {
			// 	test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
			// 	loader: "url-loader?limit=10000&mimetype=application/font-woff"
			// },
			{
				test: /\.(ttf|eot)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
				use: [{loader: "file-loader"}]
			},
			{
				test: /\.(png|svg|jpg|gif)$/,
				use: [{loader: "file-loader",
				options: {
					publicPath: "img",
					outputPath: "img"
				}}]
			},
			{
				enforce: "pre",
				test: /\.js$/,
				use: [{loader: "source-map-loader"}]
			},
			{
				test: /\.less$/,
				use: [
					{
						loader: "style-loader", // creates style nodes from JS strings
					},
					{
						loader: "css-loader", // translates CSS into CommonJS
					},
					{
						loader: "less-loader", // compiles Less to CSS
					},
				],
			}
		]
	},

	output: {
		filename: "index.js",
		path: path.resolve(__dirname, "dist")
	},

	externals: [
		{
			react: "React",
			"react-dom": "ReactDOM",
		},
	],

	devServer: {
		contentBase: path.resolve(__dirname, './dist'),
		hot: true,
	  },
};

module.exports = (env) => {

	if (env === "dev") {
		return merge(common, {
			mode: "development",
			devtool: "source-map",
		});
	}

	return merge(common, {
		mode: "production",
	});
}