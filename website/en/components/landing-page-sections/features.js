import Feature from "./feature";

export default function Features({
  headline,
  subheadline,
  features,
  scrollAnchorId,
}) {
  return (
    <section id={scrollAnchorId} className="feature-section">
      <div className="container">
        <div className="row">
          <div className="col-lg-12" style={{ paddingLeft: "12px" }}>
            <div className="row">
              {[
                {
                  headline: "Defenses For OWASP Attacks",
                  description:
                    "SafeLine use as an important tool to defense against OWASP Top 10 Attack, such as SQL injection, XSS, Insecure deserialization etc.",
                  icon: "icon-OWASP",
                },
                {
                  headline: "Defenses For 0-Day Attacks",
                  description:
                    "SafeLine use intelligent rule-free detection algorithm to against 0-Day attacks with unknown attack signatures.",
                  icon: "icon-a-0day",
                },
                {
                  headline: "Proactive Bot defense",
                  description:
                    "SafeLine uses advanced algorithms to send capthcha challenge for suspicious users to against automated robot attacks.",
                  icon: "icon-zhudongfangyu",
                },
                {
                  headline: "In-Browser Code Encryption",
                  description:
                    "SafeLine can dynamically encrypt and obfuscate static code in the browser (such as HTML, JavaScript) to against reverse engineering.",
                  icon: "icon-liulanqidaimajiami",
                },
                {
                  headline: "Web Authentication",
                  description:
                    "SafeLine prompting the user for authentication to web apps that lacks valid authentication credentials, Illegal users will be blocked.",
                  icon: "icon-shenfenyanzheng",
                },
                {
                  headline: "Web Access Control List",
                  description:
                    "SafeLine offering fine-grained control over traffic allows you to define a set of rules that determine which requests are allowed or denied.",
                  icon: "icon-a-Webkongzhifangwen",
                },
              ].map(({ headline, description, icon }) => (
                <Feature
                  key={headline}
                  headline={headline}
                  description={description}
                  icon={icon}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
