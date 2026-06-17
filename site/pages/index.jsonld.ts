import { renderDocument } from "$/jsonld";

export async function GET() {
  const document = JSON.stringify(await renderDocument(), null, 2);
  return new Response(document, {
    headers: {
      "Content-Type": "application/json",
      "Content-Disposition": "inline; filename=document.jsonld",
    },
  });
}
