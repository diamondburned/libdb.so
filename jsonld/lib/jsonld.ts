import type { JsonLdDocument } from "jsonld";

export async function fetchDocument<
  T extends object = Record<string, unknown>,
  ExtrasT extends object = JsonLdDocument
>(url: string, extras?: ExtrasT): Promise<T & ExtrasT> {
  const resp = await fetch(url);
  if (!resp.ok) {
    throw new Error(`Failed to fetch document: ${resp.statusText}`);
  }

  let doc = (await resp.json()) as T;
  if (typeof doc !== "object") {
    throw new Error(`Expected document to be an object, got ${typeof doc}`);
  }
  if (extras) {
    doc = { ...doc, ...extras };
  }

  let entries = Object.entries(doc as Record<string, unknown>);
  entries = [
    ...entries.filter(([key]) => key.startsWith("@")),
    ...entries.filter(([key]) => !key.startsWith("@")),
  ];
  return Object.fromEntries(entries) as T & ExtrasT;
}

export function dedent(str: TemplateStringsArray): string {
  const lines = str.join("").split("\n");
  return lines
    .map((line) => line.trim())
    .join(" ")
    .trim();
  // const minIndent = Math.min(
  //   ...lines.map((line) => (line.match(/^ */) ?? [""])[0].length)
  // );
  // return lines
  //   .map((line) => line.slice(minIndent))
  //   .join(" ")
  //   .trim();
}

export function gitFile(uri: `${string}:${string}`): string {
  const url = new URL(uri);
  switch (url.protocol) {
    case "github:": {
      const [user, repo, ...parts] = url.pathname.split("/");
      const path = parts.join("/");
      const ref = url.searchParams.get("ref") ?? "main";
      return `https://raw.githubusercontent.com/${user}/${repo}/${ref}/${path}`;
    }
    default: {
      throw new Error(`Unsupported protocol: ${url.protocol}`);
    }
  }
}

export function assert(condition: boolean, message?: string): asserts condition {
  if (!condition) {
    throw new Error(message);
  }
}
