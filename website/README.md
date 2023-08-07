# Website

使用 [Docusaurus 2](https://docusaurus.io/), 作为基础框架。

### 开发

```sh
# 代码变动后可以自动更新，但是不能
npm start
# 支持搜索功能，但是无法自动更新
npm run preview
```

### 部署

手动本地构建 docker 镜像，然后运行

```sh
docker build -t website:latest .
docker run --name site -p 3000:80 -d website:latest
```

### 链接替换

使用 nginx rewrite 把更改地址的链接记录下，运行旧链接访问到新地址

```nginx

location / {

    rewrite /posts/guide_introduction /docs/;
    rewrite /posts/guide_install /docs/上手指南/guide_install;
    rewrite /posts/guide_login /docs/上手指南/guide_login;
    rewrite /posts/guide_config /docs/上手指南/guide_config;
    rewrite /posts/guide_test /docs/上手指南/guide_test;
    rewrite /posts/guide_upgrade /docs/上手指南/guide_upgrade;
    rewrite /posts/faq_install /docs/常见问题排查/faq_install;
    rewrite /posts/faq_login /docs/常见问题排查/faq_login;
    rewrite /posts/faq_access /docs/常见问题排查/faq_access;
    rewrite /posts/faq_config /docs/常见问题排查/faq_config;
    rewrite /posts/faq_other /docs/常见问题排查/faq_other;
    rewrite /posts/about_syntaxanalysis /docs/关于雷池/about_syntaxanalysis;
    rewrite /posts/about_challenge /docs/关于雷池/about_challenge;
    rewrite /posts/about_changelog /docs/关于雷池/about_changelog;
    rewrite /posts/about_chaitin /docs/关于雷池/about_chaitin;

    proxy_pass http://upstream;
}
```