export const IP = "https://skill-force.ru";

export type FetchResponse<T> = {
  data?: T;
  error?: string;
  status?: number;
  ok: boolean;
};

export async function apiFetch<T = any>(
  url: string,
  options = {}
): Promise<FetchResponse<T>> {
  try {
    const response = await fetch(`${IP}/api${url}`, {
      credentials: "include",
      ...options,
    });

    const data = await response.json();
    if (!response.ok) {
      return {
        ok: false,
        error: data.error || `Ошибка: ${response.status}`,
        status: response.status,
      };
    }

    return {
      ok: true,
      data,
      status: response.status,
    };
  } catch (error: any) {
    console.error(`Ошибка запроса к ${url}:`, error);
    return {
      ok: false,
      error: error.message || "Неизвестная ошибка",
    };
  }
}

export async function apiFetchGET<T = any>(
  url: string
): Promise<FetchResponse<T>> {
  return apiFetch<T>(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });
}

export async function apiFetchPOST<T = any>(
  url: string,
  body: any
): Promise<FetchResponse<T>> {
  const csrfToken = await fetchCSRFToken();
  return apiFetch<T>(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-CSRF-Token": csrfToken || "",
    },
    body: JSON.stringify(body),
  });
}

export async function fetchCSRFToken() {
  const response = await fetch(`${IP}/api/updateProfile`, {
    method: "GET",
    credentials: "include",
  });

  const csrfToken = response.headers.get("X-Csrf-Token");
  if (!csrfToken) {
    console.error("CSRF token not received");
    return null;
  }

  return csrfToken;
}
