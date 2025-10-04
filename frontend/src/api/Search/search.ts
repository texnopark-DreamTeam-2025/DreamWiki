import { apiFetchPOST, type FetchResponse } from "@/api/api";

export async function search(word: string): Promise<FetchResponse<any>> {
  const res = await apiFetchPOST("/v1/search", { word });
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
