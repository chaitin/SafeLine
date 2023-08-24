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

    rewrite /posts/guide_introduction /docs/ permanent;
    rewrite /posts/guide_install /docs/guide/install permanent;
    rewrite /docs/上手指南/guide_install /docs/guide/install permanent;

    rewrite /posts/guide_login /docs/guide/login permanent;
    rewrite /docs/上手指南/guide_login /docs/guide/login permanent;

    rewrite /posts/guide_config /docs/guide/config permanent;
    rewrite /docs/上手指南/guide_config /docs/guide/config permanent;

    rewrite /posts/guide_test /docs/guide/test permanent;
    rewrite /docs/上手指南/guide_test /docs/guide/test permanent;

    rewrite /posts/guide_upgrade /docs/guide/upgrade permanent;
    rewrite /docs/上手指南/guide_upgrade /docs/guide/upgrade permanent;

    rewrite /posts/faq_install /docs/faq/install permanent;
    rewrite /docs/常见问题排查/faq_install /docs/faq/install permanent;

    rewrite /posts/faq_login /docs/faq/login permanent;
    rewrite /docs/常见问题排查/faq_login /docs/faq/login permanent;

    rewrite /posts/faq_access /docs/guide/config permanent;
    rewrite /docs/常见问题排查/faq_access /docs/guide/config permanent;

    rewrite /posts/faq_config /docs/faq/config permanent;
    rewrite /docs/常见问题排查/faq_config /docs/faq/config permanent;

    rewrite /posts/faq_other /docs/faq/other permanent;
    rewrite /docs/常见问题排查/faq_other /docs/faq/other permanent;

    rewrite /posts/about_syntaxanalysis /docs/about/syntaxanalysis permanent;
    rewrite /docs/关于雷池/about_syntaxanalysis /docs/about/syntaxanalysis permanent;

    rewrite /posts/about_challenge /docs/about/challenge permanent;
    rewrite /docs/关于雷池/about_challenge /docs/about/challenge permanent;

    rewrite /posts/about_changelog /docs/about/changelog permanent;
    rewrite /docs/关于雷池/about_changelog /docs/about/changelog permanent;

    rewrite /posts/about_chaitin /docs/about/chaitin permanent;
    rewrite /docs/关于雷池/about_chaitin /docs/about/chaitin permanent;

    rewrite /docs/faq/access /docs/guide/config permanent;
    rewrite /docs/faq/config /docs/guide/config permanent;

    proxy_pass http://upstream;
}
```