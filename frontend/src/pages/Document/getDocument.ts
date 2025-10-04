import { apiFetchPOST, type FetchResponse } from "@/api/api";

export async function getPageInfo(pageId: string): Promise<FetchResponse<any>> {
  const res = await apiFetchPOST("/v1/diagnostic-info/get", {
    page_id: pageId,
  });

  if (res.ok && res.data) {
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
