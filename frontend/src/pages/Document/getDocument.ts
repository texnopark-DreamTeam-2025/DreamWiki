import { getDiagnosticInfo } from "@/client";

export async function getPageInfo(pageId: string) {
  const res = await getDiagnosticInfo({
    body: {
      page_id: pageId,
    },
  });

  return res;
}
