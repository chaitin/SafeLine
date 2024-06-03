/** @type {import('next-sitemap').IConfig} */
const nextSiteMapConfig = {
  siteUrl: 'https://waf-ce.chaitin.cn',
  generateRobotsTxt: true,
  robotsTxtOptions: {
    policies: [{ userAgent: '*', allow: '/', disallow: '' }],
  },
  sitemap: {
    // path: '/sitemap.xml',
    routes: {
      '/community': {
        changefreq: 'always',
      },
    },
  },
  autoLastmod: true,
  priority: 1,
  changefreq: 'daily',
  sitemapSize: 5000,
  transform: async (config, path) => {
    if (!path) {
      return null
    }
    const customFields = config.sitemap.routes[path] || {}

    return {
      loc: path,
      changefreq: customFields.changefreq || config.changefreq,
      priority: config.priority,
      lastmod: config.autoLastmod ? new Date().toISOString() : undefined,
      alternateRefs: config.alternateRefs ?? [],
    }
  },
  additionalPaths: (config) => {
    const paths = ['/docs']
    const result = []
    paths.forEach(async (item) => {
      result.push(await config.transform(config, item))
    })
    return result
  },
}

module.exports = nextSiteMapConfig
