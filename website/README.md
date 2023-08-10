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
    rewrite /posts/guide_install /docs/guide/install;
    rewrite /posts/guide_login /docs/guide/login;
    rewrite /posts/guide_config /docs/guide/config;
    rewrite /posts/guide_test /docs/guide/test;
    rewrite /posts/guide_upgrade /docs/guide/upgrade;
    rewrite /posts/faq_install /docs/faq/install;
    rewrite /posts/faq_login /docs/faq/login;
    rewrite /posts/faq_access /docs/faq/access;
    rewrite /posts/faq_config /docs/faq/config;
    rewrite /posts/faq_other /docs/faq/other;
    rewrite /posts/about_syntaxanalysis /docs/about/syntaxanalysis;
    rewrite /posts/about_challenge /docs/about/challenge;
    rewrite /posts/about_changelog /docs/about/changelog;
    rewrite /posts/about_chaitin /docs/about/chaitin;

    proxy_pass http://upstream;
}
```