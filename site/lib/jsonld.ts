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
