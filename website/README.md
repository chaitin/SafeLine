# 雷池社区官网

使用 [Docusaurus 2](https://docusaurus.io/) 开发。包含两部分内容：页面和文档。

- 页面，在 src 目录下，手写 js
- 文档，在 docs 目录下，手写 markdown

### 开发

```sh
# 开发
npm start
# 支持 搜索功能的预览
npm run serve -- --build --host 0.0.0.0
```

### 部署

```
docker build -t website .
```