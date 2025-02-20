import nextra from "nextra";

const withNextra = nextra({
	latex: true,
	search: {
		codeblocks: true,
	},
});

export default withNextra({
	reactStrictMode: true,
	images: {
		remotePatterns: [
			{
				protocol: "https",
				hostname: "assets.seobotai.com",
			},
		],
	},
});
