import Script from 'next/script'

export default function CNZZScript() {
  return (
    <>
      <Script
        id="cnzz-js"
        src="https://v1.cnzz.com/z_stat.php?id=1281262430&web_id=1281262430"
        strategy="afterInteractive"
      />
    </>
  )
};