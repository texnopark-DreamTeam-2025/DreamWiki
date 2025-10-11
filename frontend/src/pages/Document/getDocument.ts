import { getDiagnosticInfo } from "@/client";

export async function getPageInfo(pageId: string) {
  const res = await getDiagnosticInfo({
    body: {
      page_id: pageId,
    },
  });

  if (!res.error && res.data) {
    return {
      ok: true,
      status: res.status,
      data: res.data,
      error: undefined,
    };
  }

  return {
    ok: false,
    status: res.status,
    data: undefined,
    error: res.error,
  };
}
