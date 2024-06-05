import Image from "next/image";

export default function TwoColumnWithImage({
  headline,
  subheadline,
  image,
  imagePosition,
  scrollAnchorId,
}) {
  return (
    <section id={scrollAnchorId} className="cta-section">
      <div className="container">
        <div className="row">
          {image && imagePosition === "left" && (
            <div className="col-lg-6 order-last order-lg-first">
              <div className="left-image cta-image ">
                <Image
                  src={`/images/${image}`}
                  height={400}
                  width={600}
                  alt=""
                  sizes="100vw"
                  style={{
                    width: "100%",
                    height: "auto",
                  }}
                />
              </div>
            </div>
          )}
          <div className="col-lg-6">
            <div className="cta-content-wrapper">
              <div className="section-title">
                <h2 className="mb-20">{headline}</h2>
                <div
                  style={{ color: "rgba(0,0,0,0.7)" }}
                  dangerouslySetInnerHTML={{ __html: subheadline }}
                />
              </div>
            </div>
          </div>
          {image && imagePosition === "right" && (
            <div className="col-lg-6">
              <div className="right-image cta-image text-lg-end">
                <Image
                  src={`/images/${image}`}
                  height={400}
                  width={600}
                  alt=""
                  sizes="100vw"
                  style={{
                    width: "100%",
                    height: "auto",
                  }}
                />
              </div>
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
