import { assert } from "./lib/jsonld.js";
import context from "./context.json";
import jsonld, { type ContextDefinition } from "jsonld";

const jsonldFiles = import.meta.glob("./*.jsonld{.json,.js,.ts}", { eager: true }) as Record<
  string,
  { default: jsonld.NodeObject[] }
>;

export async function renderDocument(): Promise<jsonld.JsonLdDocument> {
  const allNodes: jsonld.NodeObject[] = [];

  for (let [_, module] of Object.entries(jsonldFiles)) {
    const node = module.default;
    assertNodeObjects(node);
    allNodes.push(...node);
  }

  let document: jsonld.JsonLdDocument = {
    "@context": context,
    "@graph": expandIDsRecursively(context, allNodes),
  };

  document = await jsonld.compact(document, context, {
    compactArrays: true,
    compactToRelative: true,
  });

  // jsonld.compact messes up the key order of objects, so we need to reorder them.
  document["@graph"] = (document["@graph"] as jsonld.NodeObject[]).map((node, i) =>
    reorderSameObject(node, allNodes[i]),
  );

  return document;
}

function expandIDsRecursively<T extends jsonld.NodeObject | jsonld.NodeObject[]>(
  context: ContextDefinition,
  obj: T,
): T {
  if (Array.isArray(obj)) {
    return obj.map((item) =>
      typeof item === "object" //
        ? expandIDsRecursively(context, item as jsonld.NodeObject)
        : item,
    ) as T;
  }

  let id = obj["@id"];
  if (!id) {
    return obj;
  }

  if (typeof id !== "string") {
    throw new Error(`Expected @id to be a string, got ${typeof id}`);
  }

  if (id.includes(":") && !id.includes("://")) {
    const [prefix, suffix] = id.split(":");
    const expanded = context[prefix];
    if (expanded) {
      id = expanded + suffix;
    }
  }

  try {
    const url = new URL(id);
    id = url.href;
  } catch (err) {
    throw new Error(`Failed to parse ID as URL "${id}" in object ${JSON.stringify(obj)}`, {
      cause: err,
    });
  }

  obj["@id"] = id;

  for (const [key, value] of Object.entries(obj)) {
    if (typeof value == "object") {
      obj[key] = expandIDsRecursively(context, value as jsonld.NodeObject | jsonld.NodeObject[]);
    }
  }

  return obj;
}

function reorderSameObject<T extends Record<string, unknown>>(unordered: T, ordered: T): T {
  const result: T = {} as T;
  for (const key of Object.keys(ordered)) {
    result[key as keyof T] = unordered[key] as T[keyof T];
  }
  return result;
}

function assertNodeObjects(doc: unknown): asserts doc is jsonld.NodeObject[] {
  assert(Array.isArray(doc), `Expected document to be an array, got ${typeof doc}`);
  assert(
    doc.every((item) => typeof item === "object"),
    "Expected document to be an array of objects",
  );
}
