import { apiFetchGET, FetchResponse } from "@/api/api";

export async function search(): Promise<FetchResponse<any>> {
  const res = await apiFetchGET("/search");
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
