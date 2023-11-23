import Head from 'next/head'
import '/styles/globals.css'
import NavBar from '@/components/NavBar'
import Footer from '@/components/Footer'
import CNZZScript from '@/components/CNZZ'
import ThemeProvider from "@/components/Theme";

export default function MyApp({ Component, pageProps }: any) {
  return (
    <div className='overflow-x-hidden'>
      <Head>
        <title>雷池 WAF 社区版 | 下一代 Web 应用防火墙 | 免费安装</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="keywords" content="WAF,雷池,长亭,社区版,免费,开源,网站防护"></meta>
        <meta name="description" content="一款足够简单、好用、强大的免费WAF。基于业界领先的语义分析检测技术，保护你的网站不受黑客攻击。"></meta>
      </Head>
      <ThemeProvider>
        <NavBar />
        <Component {...pageProps} />
        <Footer />
      </ThemeProvider>
      <CNZZScript />
    </div>
  )
}