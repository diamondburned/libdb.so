import { fetchDocument, gitFile, dedent } from "./lib/jsonld.js";
import context from "./_context.json";

const graph = [
  {
    "@type": "Person",
    "@id": "https://0xd14.id",
    additionalType: ["Bot", "zvava:Bot"],
    identifier: "d14",
    name: "Diamond",
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
    "zvava:bot": true,
    "zvava:pronouns": ["it/its", "she/her"],
    "libdb:webring": [
      {
        "@context": { "@vocab": `${context.libdb}Webring/` },
        "@type": "libdb:Webring",
        version: 1,
        name: "robot girls",
        ring: [
          {
            name: "diamond",
            link: "https://0xd14.id",
          },
          {
            name: "SZOFIÁ",
            link: "https://737a6f6669e1.id",
          },
        ],
      },
      await fetchDocument(
        gitFile("github:diamondburned/acmfriends-webring/webring.json?ref=<3-spring-2023"),
        {
          "@context": { "@vocab": `${context.libdb}Webring/` },
          "@type": "libdb:Webring",
        }
      ),
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

function githubFile(ownerRepo: string, path: string, branch = "main") {
  return `https://raw.githubusercontent.com/${ownerRepo}/${branch}/${path}`;
}
