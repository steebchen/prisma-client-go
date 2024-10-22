/** @type {import('next-sitemap').IConfig} */
module.exports = {
  siteUrl: 'https://goprisma.org',
  generateRobotsTxt: true,
  exclude: [
    '/docs/README',
    '*/_meta',
  ],
}
