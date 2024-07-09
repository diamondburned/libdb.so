import { fetchDocument, fetchJSON, gitFile, dedent } from "./lib/jsonld.js";
import { type WebringData } from "libwebring/lib/webring.js";
import context from "./_context.json";

const _88x31 = (name: string) => "https://libdb.so/_fs/88x31/" + name;
const _88x31s = [];

const graph = [
  {
    "@id": "https://0xd14.id",
    "@type": "Person",
    name: ["Diamond", "d14"],
    identifier: "0xd14",
    url: [
      "https://0xd14.id", //
      "https://libdb.so",
      "https://diamondx.pet",
    ],
    email: [
      "mailto:diamond@libdb.so", //
    ],
    sameAs: [
      "https://tech.lgbt/@diamond",
      "https://girlcock.club/@diamond",
      "https://github.com/diamondburned",
    ],
    height: "169 cm",
    gender: [
      "girlthing", //
      "bot",
      "female",
      "transgender",
    ],
    birthDate: "2002-02-27",
    knowsLanguage: ["en", "vi"],
    description: dedent`
      Hi, I'm Diamond! I'm a 4th-year Computer Science major 👩🎓 and past
	    Software Engineer Intern 👩‍💻 🖥️.
	  `,
    additionalType: ["Bot", "zvava:Bot"],
    "zvava:bot": true,
    "zvava:pronouns": ["it/its", "she/her"],
    "libdb:88x31": [
      {
        "libdb:Other88x31/link": "#",
        "libdb:Other88x31/image": _88x31("d14.gif"),
      },
      {
        "libdb:Other88x31/link": "https://noe.sh/dollcode/",
        "libdb:Other88x31/image": _88x31("d14-dollcode.png"),
      },
    ],
    "libdb:other88x31": [
      {
        "libdb:Other88x31/alt": "zvava.org",
        "libdb:Other88x31/link": "https://zvava.org",
        "libdb:Other88x31/image": _88x31("zvava.org.png"),
      },
      {
        "libdb:Other88x31/alt": "ezri.pet",
        "libdb:Other88x31/link": "https://ezri.pet",
        "libdb:Other88x31/image": _88x31("ezri.pet.png"),
      },
      {
        "libdb:Other88x31/alt": "Not A Person",
        "libdb:Other88x31/link": "https://voidgoddess.org/emptyspaces/notaperson/",
        "libdb:Other88x31/image": _88x31("nap.png"),
      },
      {
        "libdb:Other88x31/alt": "Mozilla",
        "libdb:Other88x31/link": "https://www.mozilla.org",
        "libdb:Other88x31/image": _88x31("www.mozilla.org.png"),
      },
      {
        "libdb:Other88x31/alt": "hrt.coffee",
        "libdb:Other88x31/link": "https://hrt.coffee",
        "libdb:Other88x31/image": _88x31("hrt.coffee.gif"),
      },
    ],
    "libdb:dollcode": "▖▖▖▖▌▘▘",
    "libdb:webring": [
      await (async () => {
        const webring = await fetchJSON<Record<string, string>[]>(
          gitFile("github:diamondburned/roboring/websites.json")
        );
        return {
          "@context": { "@vocab": `${context.libdb}Webring/` },
          "@id": "libdb:webring/roboring",
          "@type": "libdb:Webring",
          version: 1,
          name: "roboring",
          link: "https://stellophiliac.github.io/roboring",
          self: "diamond #0xd14",
          ring: webring.map((site, i) => {
            const isNext = webring[(i + 1) % webring.length]?.slug == "0xd14";
            const isPrev = webring[(webring.length - i - 1) % webring.length]?.slug == "0xd14";
            return {
              "@id": `https://stellophiliac.github.io/roboring#${site.slug}`,
              name: site.name,
              link:
                isNext || isPrev
                  ? `https://stellophiliac.github.io/roboring/0xd14/${isPrev ? "previous" : "next"}`
                  : site.url,
            };
          }),
        };
      })(),
      await (async () => {
        const webring = await fetchDocument<WebringData>(
          gitFile("github:diamondburned/acmfriends-webring/webring.json?ref=<3-spring-2023"),
          {
            "@context": { "@vocab": `${context.libdb}Webring/` },
            "@id": "libdb:webring/acmfriends",
            "@type": "libdb:Webring",
            link: "https://github.com/diamondburned/acmfriends-webring",
          }
        );
        webring.ring = webring.ring.map((site) => ({
          "@id": "libdb:webring/acmfriends/" + site.name,
          name: site.name,
          link: site.link.includes("://") ? site.link : `https://${site.link}`,
        }));
        return webring;
      })(),
    ],
    // resume: await fetchDocument(
    //   githubFile("diamondburned", "resume", "resume.json"), //
    //   {
    //     "@context": "libdb:Resume",
    //     "@type": "libdb:Resume",
    //   }
    // ),
  },
];

export default graph;
