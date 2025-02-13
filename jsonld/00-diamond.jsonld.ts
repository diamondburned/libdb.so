import { fetchDocument, fetchJSON, gitFile, dedent } from "./lib/jsonld.js";
import { type WebringData } from "libwebring/lib/webring.js";
import context from "./_context.json";

const _88x31 = (name: string) => "https://libdb.so/_fs/pics/88x31/" + name;

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
        "libdb:88x31Badge/link": "#",
        "libdb:88x31Badge/image": _88x31("d14.gif"),
      },
      {
        "libdb:88x31Badge/link": "https://0xd14.id",
        "libdb:88x31Badge/image": _88x31("d14-dollcode.png"),
      },
    ],
    "libdb:other88x31": [
      {
        "libdb:88x31Badge/alt": "My Mastodon",
        "libdb:88x31Badge/link": "https://tech.lgbt/@diamond",
        "libdb:88x31Badge/image": _88x31("mastodon_button_3.gif"),
      },
      {
        "libdb:88x31Badge/alt": "My GitHub",
        "libdb:88x31Badge/link": "https://github.com/diamondburned",
        "libdb:88x31Badge/image": _88x31("github.png"),
      },
      {
        "libdb:88x31Badge/alt": "sophie's brain out on a table",
        "libdb:88x31Badge/link": "https://zvava.org",
        "libdb:88x31Badge/image": _88x31("zvava.org.png"),
      },
      {
        "libdb:88x31Badge/alt": 'unit ⎈-657a7269, "ezri"',
        "libdb:88x31Badge/link": "https://ezri.pet",
        "libdb:88x31Badge/image": _88x31("ezri.pet.png"),
      },
      {
        "libdb:88x31Badge/alt": "bunbun.dev",
        "libdb:88x31Badge/link": "https://bunbun.dev",
        "libdb:88x31Badge/image": _88x31("bunbun.dev.gif"),
      },
      {
        "libdb:88x31Badge/alt": "joel's website",
        "libdb:88x31Badge/link": "https://jdr.sh",
        "libdb:88x31Badge/image": _88x31("jdr.sh.png"),
      },
      {
        "libdb:88x31Badge/alt": "Not A Person",
        "libdb:88x31Badge/link": "https://voidgoddess.org/emptyspaces/notaperson/",
        "libdb:88x31Badge/image": _88x31("nap.png"),
      },
      {
        "libdb:88x31Badge/alt": "i'm a hypnoslut! <3 // made by bunbun",
        "libdb:88x31Badge/link": "#",
        "libdb:88x31Badge/image": _88x31("hypnoslut.png"),
      },
      {
        "libdb:88x31Badge/alt": "Mozilla",
        "libdb:88x31Badge/link": "https://www.mozilla.org",
        "libdb:88x31Badge/image": _88x31("www.mozilla.org.png"),
      },
      {
        "libdb:88x31Badge/alt": "hrt.coffee",
        "libdb:88x31Badge/link": "https://hrt.coffee",
        "libdb:88x31Badge/image": _88x31("hrt.coffee.gif"),
      },
      {
        "libdb:88x31Badge/alt": "Xenia Trans Pride",
        "libdb:88x31Badge/link": "https://www.kernel.org",
        "libdb:88x31Badge/image": _88x31("xeniatrans_now.gif"),
      },
      {
        "libdb:88x31Badge/alt": "Transfem Science",
        "libdb:88x31Badge/link": "https://transfemscience.org",
        "libdb:88x31Badge/image": _88x31("trans-your-gender.gif"),
      },
      {
        "libdb:88x31Badge/alt": "Powered by NixOS",
        "libdb:88x31Badge/link": "https://nixos.org/",
        "libdb:88x31Badge/image": _88x31("powered_by_nixos.gif"),
      },
      {
        "libdb:88x31Badge/alt": "Built with Nix",
        "libdb:88x31Badge/link": "https://nixos.org/",
        "libdb:88x31Badge/image": _88x31("built_with_nix.gif"),
      },
      {
        "libdb:88x31Badge/alt": "Let's all love Lain",
        "libdb:88x31Badge/link": "https://fauux.neocities.org/lovelain",
        "libdb:88x31Badge/image": _88x31("lain.gif"),
      },
      {
        "libdb:88x31Badge/alt": "Join Mastodon!",
        "libdb:88x31Badge/link": "https://joinmastodon.org",
        "libdb:88x31Badge/image": _88x31("mastodon-flat.png"),
      },
    ],
    "libdb:dollcode": "▖▖▖▖▌▘▘▌",
    "libdb:webring": [
      await (async () => {
        const webring = await fetchJSON<Record<string, string>[]>(
          gitFile("github:stellophiliac/roboring/websites.json"),
        );
        return {
          "@context": { "@vocab": `${context.libdb}Webring/` },
          "@id": "libdb:webring/roboring",
          "@type": "libdb:Webring",
          version: 1,
          name: "roboring",
          root: "https://stellophiliac.github.io/roboring",
          self: "diamond #0xd14",
          ring: webring.map((site) => {
            return {
              "@id": `https://stellophiliac.github.io/roboring#${site.slug}`,
              name: site.name,
              link: site.url,
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
          },
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
