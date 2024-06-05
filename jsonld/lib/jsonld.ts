export async function fetchDocument(
  url: string,
  extras: Record<string, unknown> = {}
): Promise<Record<string, unknown>> {
  const resp = await fetch(url);
  if (!resp.ok) {
    throw new Error(`Failed to fetch document: ${resp.statusText}`);
  }

  const doc = {
    ...(await resp.json()),
    ...extras,
  };

  let entries = Object.entries(doc);
  entries = [
    ...entries.filter(([key]) => key.startsWith("@")),
    ...entries.filter(([key]) => !key.startsWith("@")),
  ];
  return Object.fromEntries(entries);
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
