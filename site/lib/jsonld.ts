/*
export const context = {
  "@vocab": "https://schema.org/",
  rdf: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
  rdfs: "http://www.w3.org/2000/01/rdf-schema#",
  zvava: "https://zvava.org/schema#",
  libdb: "https://0xd14.id#",
};

export function expand(key: `${keyof typeof context}:${string}`) {
  const [prefix, value] = key.split(":");
  return context[prefix as keyof typeof context] + value;
}

expand("libdb:toy"); // type checks
*/

import jsonld from "jsonld";

export type JSONLDDocument = {
  "@context": Record<string, string>;
  "@graph": jsonld.NodeObject[];
};

export type FetchedJSONLD = {
  root: JSONLDDocument;
  self: jsonld.NodeObject & {
    // TypeScript. Shut up.
    [key: string]: any;
  };
};

export const jsonldURL = "/_fs/0xd14.jsonld";

export const jsonldFetch = () =>
  fetch(jsonldURL)
    .then((r) => r.json() as Promise<JSONLDDocument>)
    .then(
      (r) =>
        ({
          root: r,
          self: r["@graph"][0] as jsonld.NodeObject & {
            [key: string]: any;
          },
        }) as FetchedJSONLD,
    )
    .catch((err) => {
      console.error("Failed to fetch self JSON-LD:", err);
      throw err;
    });

export const jsonldDoc = jsonldFetch();

export async function usingContext<T = object>(
  doc: Awaited<typeof jsonldDoc>,
  object: any,
  vocab: string,
): Promise<T> {
  return (await jsonld.compact(
    {
      "@context": doc.root["@context"],
      ...object,
    },
    {
      "@vocab": vocab,
    },
  )) as T;
}
