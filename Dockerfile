FROM node:20.5-alpine

RUN apk update
RUN apk add nginx tini

RUN echo -e "                                                                   \n\
server {                                                                        \n\
    listen 80;                                                                  \n\
                                                                                \n\
    location /api/count {                                                       \n\
        proxy_pass https://rivers-telemetry.chaitin.cn:10086;                   \n\
    }                                                                           \n\
    location /api/exist {                                                       \n\
        proxy_pass https://rivers-telemetry.chaitin.cn:10086;                   \n\
    }                                                                           \n\
    location /blazehttp {                                                       \n\
        root /app/;                                                             \n\
        try_files \$uri =404;                                                   \n\
    }                                                                           \n\
    location /release {                                                         \n\
        root /app/;                                                             \n\
        try_files \$uri =404;                                                   \n\
    }                                                                           \n\
    location / {                                                                \n\
        rewrite /posts/guide_introduction /docs/ permanent;                     \n\
        rewrite /posts/guide_install /docs/guide/install permanent;             \n\
        rewrite /docs/上手指南/guide_install /docs/guide/install permanent;      \n\
        rewrite /posts/guide_login /docs/guide/login permanent;                 \n\
        rewrite /docs/上手指南/guide_login /docs/guide/login permanent;          \n\
        rewrite /posts/guide_config /docs/guide/config permanent;               \n\
        rewrite /docs/上手指南/guide_config /docs/guide/config permanent;        \n\
        rewrite /posts/guide_test /docs/guide/test permanent;                   \n\
        rewrite /docs/上手指南/guide_test /docs/guide/test permanent;            \n\
        rewrite /posts/guide_upgrade /docs/guide/upgrade permanent;             \n\
        rewrite /docs/上手指南/guide_upgrade /docs/guide/upgrade permanent;      \n\
        rewrite /posts/faq_install /docs/faq/install permanent;                 \n\
        rewrite /docs/常见问题排查/faq_install /docs/faq/install permanent;       \n\
        rewrite /posts/faq_login /docs/faq/login permanent;                     \n\
        rewrite /docs/常见问题排查/faq_login /docs/faq/login permanent;           \n\
        rewrite /posts/faq_access /docs/guide/config permanent;                 \n\
        rewrite /docs/常见问题排查/faq_access /docs/guide/config permanent;       \n\
        rewrite /posts/faq_config /docs/faq/config permanent;                   \n\
        rewrite /docs/常见问题排查/faq_config /docs/faq/config permanent;         \n\
        rewrite /posts/faq_other /docs/faq/other permanent;                     \n\
        rewrite /docs/常见问题排查/faq_other /docs/faq/other permanent;           \n\
        rewrite /posts/about_syntaxanalysis /docs/about/syntaxanalysis permanent;         \n\
        rewrite /docs/关于雷池/about_syntaxanalysis /docs/about/syntaxanalysis permanent;  \n\
        rewrite /posts/about_challenge /docs/about/challenge permanent;         \n\
        rewrite /docs/关于雷池/about_challenge /docs/about/challenge permanent;  \n\
        rewrite /posts/about_changelog /docs/about/changelog permanent;         \n\
        rewrite /docs/关于雷池/about_changelog /docs/about/changelog permanent;  \n\
        rewrite /posts/about_chaitin /docs/about/chaitin permanent;             \n\
        rewrite /docs/关于雷池/about_chaitin /docs/about/chaitin permanent;      \n\
        rewrite /docs/faq/access /docs/guide/config permanent;                  \n\
        rewrite /docs/faq/config /docs/guide/config permanent;                  \n\
        proxy_pass http://127.0.0.1:3000;                                       \n\
    }                                                                           \n\
}                                                                               \n\
" > /etc/nginx/http.d/default.conf
RUN sed -i 's/access_log/access_log off; #/' /etc/nginx/nginx.conf
RUN nginx -t

COPY release /app/release
COPY blazehttp /app/blazehttp

COPY website /app
WORKDIR /app
RUN npm ci
RUN npm run build

CMD nginx; tini -- npm run serve
