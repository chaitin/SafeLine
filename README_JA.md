<p align="center">
  <img src="/images/banner.png" width="400" />
</p>

<h4 align="center">
  SafeLine - Web アプリケーションを安全に
</h4>

<p align="center">
  <a target="_blank" href="https://ly.safepoint.cloud/laA8asp">🏠 Website</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://ly.safepoint.cloud/w2AeHhb">📖 Docs</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://ly.safepoint.cloud/hSMd4SH">🔍 Live Demo</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://discord.gg/SVnZGzHFvn">🙋‍♂️ Discord</a> &nbsp; | &nbsp;
  <a target="_blank" href="/README.md">English</a> &nbsp; | &nbsp;
  <a target="_blank" href="/README_CN.md">中文版</a>
</p>

## 👋 はじめに

SafeLine は、Web アプリケーションを攻撃やエクスプロイトから保護するためのセルフホスト型 **`WAF（Web Application Firewall：Web アプリケーションファイアウォール）`** です。

Web アプリケーションファイアウォールは、Web アプリケーションとインターネットの間を流れる HTTP トラフィックをフィルタリングおよび監視することで、Web アプリケーションの保護を支援します。一般的に、`SQL インジェクション`、`XSS`、`コードインジェクション`、`OS コマンドインジェクション`、`CRLF インジェクション`、`LDAP インジェクション`、`XPath インジェクション`、`RCE`、`XXE`、`SSRF`、`パストラバーサル`、`バックドア`、`ブルートフォース`、`HTTP フラッド`、`ボットの悪用`などの攻撃から Web アプリケーションを保護します。

#### 💡 動作の仕組み

<img src="/images/how-it-works.png" width="800" />

Web アプリケーションの前段に WAF を配置することで、Web アプリケーションとインターネットの間にシールドが設けられます。プロキシサーバーが仲介役を用いてクライアントマシンのアイデンティティを保護するのに対し、WAF はリバースプロキシの一種であり、クライアントがサーバーに到達する前に WAF を経由させることで、サーバーが直接さらされることを防ぎます。

WAF は、Web アプリケーションへ流れる悪意のある HTTP/S トラフィックをフィルタリング・監視・ブロックし、許可されていないデータがアプリケーションから流出するのを防ぐことで、Web アプリケーションを保護します。これは、どのトラフィックが悪意のあるもので、どのトラフィックが安全かを判断するための一連のポリシーに従うことで実現されます。プロキシサーバーがクライアントのアイデンティティを保護する仲介役として機能するのと同様に、WAF も同じように動作しますが、潜在的に悪意のあるクライアントから Web アプリケーションサーバーを保護するリバースプロキシの仲介役として機能します。

主な機能は次のとおりです。

- Web 攻撃に対する防御
- 能動的なボット悪用防御
- HTML と JS コードの暗号化
- IP ベースのレート制限
- Web アクセス制御リスト

#### ⚡️ スクリーンショット

| <img src="./images/screenshot-1.png" width=370 /> | <img src="./images/screenshot-2.png" width=370 /> |
| ------------------------------------------------- | ------------------------------------------------- | 
| <img src="./images/screenshot-3.png" width=370 /> | <img src="./images/screenshot-4.png" width=370 /> | 

[ライブデモ](https://demo.waf.chaitin.com:9443/) を見る

## 🔥 機能

主な機能の一覧は次のとおりです。

- **`Web 攻撃のブロック`**
  - `SQL インジェクション`、`XSS`、`コードインジェクション`、`OS コマンドインジェクション`、`CRLF インジェクション`、`XXE`、`SSRF`、`パストラバーサル`など、あらゆる Web 攻撃に対して防御します。
- **`レート制限`**
  - 定義した上限を超えるトラフィックをスロットリングすることで、`DoS 攻撃`、`ブルートフォース攻撃`、`トラフィックの急増`、その他の悪用から Web アプリケーションを保護します。
- **`アンチボットチャレンジ`**
  - アンチボットチャレンジによって `ボット攻撃` から Web サイトを保護します。人間のユーザーは許可され、クローラーやボットはブロックされます。
- **`認証チャレンジ`**
  - 認証チャレンジを有効にすると、訪問者はパスワードを入力する必要があり、入力しない場合はブロックされます。
- **`動的保護`**
  - 動的保護を有効にすると、Web サーバー上の HTML および JS コードがアクセスのたびに動的に暗号化されます。

#### 🧩 ショーケース

|                               | 正当なユーザー                                       | 悪意のあるユーザー                                               |
| ----------------------------- | --------------------------------------------------- | ---------------------------------------------------------------- | 
| **`Web 攻撃のブロック`**       | <img src="./images/skeleton.png" width=270 />       | <img src="./images/blocked-for-attack-detected.png" width=270 /> |
| **`レート制限`**               | <img src="./images/skeleton.png" width=270 />       | <img src="./images/blocked-for-access-too-fast.png" width=270 /> |
| **`アンチボットチャレンジ`**   | <img src="./images/captcha-1.gif" width=270 />      | <img src="./images/captcha-2.gif" width=270 />                     |
| **`認証チャレンジ`**           | <img src="./images/auth-1.gif" width=270 />         | <img src="./images/auth-2.gif" width=270 />                        |
| **`HTML 動的保護`**            | <img src="./images/dynamic-html-1.png" width=270 /> | <img src="./images/dynamic-html-2.png" width=270 />              |
| **`JS 動的保護`**              | <img src="./images/dynamic-js-1.png" width=270 />   | <img src="./images/dynamic-js-2.png" width=270 />                | 

## 🚀 クイックスタート

> [!WARNING]
> 中国本土のユーザーが国際版をインストールするとクラウドサービスに接続できなくなる場合があります。[中文版インストールドキュメント](https://docs.waf-ce.chaitin.cn/zh/%E4%B8%8A%E6%89%8B%E6%8C%87%E5%8D%97/%E5%AE%89%E8%A3%85%E9%9B%B7%E6%B1%A0) をご覧ください。

#### 📦 インストール

SafeLine のインストール方法については、[インストールガイド](https://docs.waf.chaitin.com/en/GetStarted/Deploy) を参照してください。

#### ⚙️ Web アプリケーションの保護

[設定](https://docs.waf.chaitin.com/en/GetStarted/AddApplication) を参照してください。

## 📋 詳細情報

#### 効果の評価

| 指標              | ModSecurity, Level 1 | CloudFlare, Free     | SafeLine, Balance      | SafeLine, Strict      |
| ----------------- | -------------------- | -------------------- | ---------------------- | --------------------- |
| サンプル総数      | 33669                | 33669                | 33669                  | 33669                 |
| **検知率**        | 69.74%               | 10.70%               | 71.65%                 | **76.17%**            |
| **誤検知率**      | 17.58%               | 0.07%                | **0.07%**              | 0.22%                 |
| **正確度**        | 82.20%               | 98.40%               | **99.45%**             | 99.38%                |


#### SafeLine は本番環境で使えますか？

はい、SafeLine は本番環境で利用可能です。

- 世界中で 180,000 件を超えるインストール実績
- 1,000,000 を超える Web サイトを保護
- 毎日 30,000,000,000 を超える HTTP リクエストを処理

#### 🙋‍♂️ コミュニティ

[Discord](https://discord.gg/SVnZGzHFvn) に参加してコミュニティのサポートを受けてください。コアチームのメンバーは Discord で STAFF ロールによって識別できます。

- チャンネル [#feedback](https://discord.com/channels/1243085666485534830/1243120292822253598)：新機能に関する議論用。
- チャンネル [#FAQ](https://discord.com/channels/1243085666485534830/1263761679619981413)：よくある質問用。
- チャンネル [#general](https://discord.com/channels/1243085666485534830/1243115843919806486)：その他の質問用。

コミュニティ向けにはいくつかの連絡手段がありますが、主なものは Discord です。これらに加えて、新しい Issue の作成には GitHub Issues もご利用いただけます。

<p align="left">
  <a target="_blank" href="https://discord.gg/SVnZGzHFvn"><img src="https://img.shields.io/badge/Discord-5865F2?style=flat&logo=discord&logoColor=white"></a> &nbsp;
  <a target="_blank" href="https://x.com/safeline_waf"><img src="https://img.shields.io/badge/X.com-000000?style=flat&logo=x&logoColor=white"></a> &nbsp;
  <a target="_blank" href="/images/wechat.png"><img src="https://img.shields.io/badge/WeChat-07C160?style=flat&logo=wechat&logoColor=white"></a>
</p>

#### 💪 PRO 版

近日公開予定！

#### 📝 ライセンス

詳細は [LICENSE](/LICENSE.md) を参照してください。
