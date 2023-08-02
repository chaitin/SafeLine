// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

const cnzz = `<script type="text/javascript" src="https://v1.cnzz.com/z_stat.php?id=1281262430&web_id=1281262430"/>`;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "长亭雷池 WAF 社区版",
  tagline: "Dinosaurs are cool",
  favicon: "images/favicon.ico",
  plugins: [require.resolve("@cmfcmf/docusaurus-search-local")],

  // Set the production url of your site here
  url: "https://waf-ce.chaitin.cn/",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "chaitin", // Usually your GitHub org/user name.
  projectName: "safeline-ce-site", // Usually your repo name.

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "zh-Hans",
    locales: ["zh-Hans"],
  },

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/chaitin/safeline/tree/main/website/docs",
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/chaitin/safeline/tree/main/website/blog",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themes: [
    [
      // @ts-ignore
      require.resolve("@easyops-cn/docusaurus-search-local"),
      /** @type {import("@easyops-cn/docusaurus-search-local").PluginOptions} */
      // @ts-ignore
      ({
        // ... Your options.
        // `hashed` is recommended as long-term-cache of index file is possible.
        hashed: true,
        // For Docs using Chinese, The `language` is recommended to set to:
        // ```
        language: ["zh"],
        // ```
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: "images/safeline.png",
      navbar: {
        title: "",
        logo: { alt: "Logo", src: "images/logo.png" },
        items: [
          {
            to: "/docs/上手指南/guide_introduction",
            label: "技术文档",
            position: "right",
          },
          { to: "/detection", label: "效果对比", position: "right" },
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
          {
            href: "https://github.com/chaitin/safeline",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: " ",
            items: [
              {
                label: "北京长亭科技有限公司",
                to: "https://www.chaitin.cn/zh/",
              },
              {
                label: "长亭 B 站主页",
                to: "https://space.bilibili.com/521870525",
              },
            ],
          },
          {
            title: " ",
            items: [
              {
                label: "CT Stack 安全社区",
                href: "https://stack.chaitin.cn/",
              },
              {
                label: "长亭合作伙伴论坛",
                href: "https://bbs.chaitin.cn/",
              },
            ],
          },
          {
            title: " ",
            items: [
              {
                label: "长亭百川云平台",
                to: "https://rivers.chaitin.cn/",
              },
              {
                label: "关于我们",
                href: "https://waf-ce.chaitin.cn/docs/关于雷池/about_chaitin",
              },
            ],
          },
          {
            title: " ",
            items: [
              {
                label: "长亭 GitHub 主页",
                to: "https://github.com/chaitin",
              },
            ],
          },
        ],
        copyright: `${
          process.env.NODE_ENV === "production" ? cnzz : ""
        } <b>Copyright</b> © ${new Date().getFullYear()} 
        北京长亭未来科技有限公司京 ICP 备 19035216 号京公网安备 11010802020947 号`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
