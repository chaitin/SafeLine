import React from "react";
export default function FooterCopyright({ copyright }) {
  return (
    <>
      <div
        className="footer__copyright"
        // Developer provided the HTML, so assume it's safe.
        // eslint-disable-next-line react/no-danger
        dangerouslySetInnerHTML={{ __html: copyright }}
      />
      {/* 自定义组件的目的是为了插入 cnzz 统计代码 */}
      {process.env.NODE_ENV == "production" && (
        <script
          type="text/javascript"
          src="https://v1.cnzz.com/z_stat.php?id=1281262430&web_id=1281262430"
        />
      )}
    </>
  );
}
