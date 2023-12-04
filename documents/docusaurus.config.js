// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

// const lightCodeTheme = require("prism-react-renderer/themes/github");
// const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "雷池 WAF 社区版",
  tagline: "",
  favicon: "images/favicon.ico",

  // Set the production url of your site here
  url: "https://waf-ce.chaitin.cn/",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/docs",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "chaitin", // Usually your GitHub org/user name.
  projectName: "document", // Usually your repo name.

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "zh-Hans",
    locales: ["zh-Hans"],
  },

  themes: [
    [
      // @ts-ignore
      require.resolve("@easyops-cn/docusaurus-search-local"),
      /** @type {import("@easyops-cn/docusaurus-search-local").PluginOptions} */
      // @ts-ignore
      ({
        docsRouteBasePath: "/",
        // ... Your options.
        // `hashed` is recommended as long-term-cache of index file is possible.
        hashed: true,
        // For Docs using Chinese, The `language` is recommended to set to:
        // ```
        language: ["en", "zh"],
        // ```
      }),
    ],
  ],

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
          // Remove this to remove the "edit this page" links.
          // editUrl: "https://github.com/chaitin/safeline/tree/main/website",
        },
        blog: false,
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: "images/safeline.svg",
      navbar: {
        title: "",
        logo: { alt: "Logo", src: "images/safeline.svg", href: "https://waf-ce.chaitin.cn" },
        items: [
          {
            to: "https://www.bilibili.com/medialist/detail/ml2342694989",
            label: "教学视频",
            position: "right",
          },
          {
            to: "https://demo.waf-ce.chaitin.cn:9443/dashboard",
            label: "演示环境",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: "雷池 SafeLine",
            items: [
               {
                label: "主页",
                to: "https://waf-ce.chaitin.cn",
              },
              {
                label: "开发计划",
                to: "https://waf-ce.chaitin.cn/community",
              },
              {
                label: "付费版本",
                to: "https://waf-ce.chaitin.cn/version",
              },
            ],
          },
          {
            title: "资源",
            items: [
              // {
              //   label: "技术文档",
              //   to: "/",
              // },
              {
                label: "教学视频",
                to: "https://www.bilibili.com/medialist/detail/ml2342694989",
              },
              // {
              //   label: "学习资料",
              //   to: "/",
              // },
              // {
              //   label: "更新日志",
              //   to: "/about/changelog",
              // },
            ],
          },
          {
            title: "关于我们",
            items: [
              {
                label: "长亭科技",
                to: "https://www.chaitin.cn/zh/",
              },
              {
                label: "CT Stack 安全社区",
                to: "https://stack.chaitin.cn/",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} 北京长亭科技有限公司.All rights reserved.`,
      },
      prism: {
        // theme: lightCodeTheme,
        // darkTheme: darkCodeTheme,
      },
      colorMode: {
        defaultMode: 'light',
        disableSwitch: true,
      },
    }),
};

module.exports = config;
