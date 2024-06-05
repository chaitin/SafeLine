import Image from "next/image";

export default function Feature({ headline, description, icon }) {
  return (
    <div className="col-lg-4 col-md-4">
      <div className="single-feature">
        <div className="feature-icon">
          <svg
            className="icon_svg"
            style={{ width: "20px", height: "20px", marginLeft: "8px" }}
          >
            <use xlinkHref={`#${icon}`} />
          </svg>
        </div>
        <div className="feature-content">
          <h4>{headline}</h4>
          <p>{description}</p>
        </div>
      </div>
    </div>
  );
}
