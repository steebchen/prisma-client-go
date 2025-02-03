const { BlogClient } = require("seobot");

/** @type {import("next-sitemap").IConfig} */
module.exports = {
	siteUrl: "https://goprisma.org",
	generateRobotsTxt: true,
	exclude: ["/docs/README", "*/_meta"],
	additionalPaths: async () => {
		const key =
			process.env.SEOBOT_API_KEY || "a8c58738-7b98-4597-b20a-0bb1c2fe5772";

		const client = new BlogClient(key);

		const { articles } = await client.getArticles(0, 10);

		return articles.map((article) => ({
			loc: `/blog/${article.slug}`,
		}));
	},
};
