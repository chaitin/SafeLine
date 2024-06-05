import { useEffect } from "react";

import Testimonial from "./testimonial";

export default function Testimonials({
  headline,
  testimonial: testimonials,
  scrollAnchorId,
}) {
  useEffect(() => {
    import("tiny-slider").then(({ tns }) => {
      tns({
        container: ".testimonial-active",
        autoplay: true,
        autoplayTimeout: 3000,
        autoplayButtonOutput: false,
        mouseDrag: true,
        gutter: 0,
        nav: false,
        navPosition: "bottom",
        controls: false,
        // controlsText: [
        //   '<i class="lni lni-chevron-left"></i>',
        //   '<i class="lni lni-chevron-right"></i>',
        // ],
        items: 1,
      });
    });
  });

  return (
    <section id={scrollAnchorId} className="testimonial-section mt-100">
      <div className="container">
        <div className="row justify-content-center">
          <div className="col-xl-7 col-lg-9"></div>
        </div>
      </div>
      <div className="testimonial-active-wrapper">
        <div className="testimonial-active">
          <Testimonial quote="1" name="2" title="3" />
          <Testimonial quote="4" name="5" title="6" />
        </div>
      </div>
    </section>
  );
}
