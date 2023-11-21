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
        <title>长亭雷池 WAF 社区版</title>
        <link rel="icon" href="/favicon.ico" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="keywords" content="WAF,雷池,长亭,社区版,免费,开源,网站防护"></meta>
        <meta name="description" content="长亭雷池 WAF 社区版"></meta>
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