import { CourseOpen, CourseStructure } from "@/types/courseMenu";
import { LessonsStructure } from "@/types/lesson";
import { UserProfile } from "@/types/users";
import { QuestionsStructure } from "./types/question";

export const IP = "https://skill-force.ru";
export const PORT = "80";

export interface Course {
  id: number;
  price: number;
  purchases_amount: number;
  creator_id: number;
  time_to_pass: number;
  title: string;
  description: string;
  rating: number;
  src_image: string;
  tags: string[];
  is_favorite: boolean;
}

export async function apiFetch(url: string, options = {}) {
  try {
    const response = await fetch(`${IP}/api${url}`, {
      credentials: "include",
      ...options,
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || `Ошибка: ${response.status}`);
    }

    return data;
  } catch (error) {
    console.error(`Ошибка запроса к ${url}:`, error);
    return null;
  }
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

export async function searchForm(keyword: string) {
  const data = await apiFetch(`/v1/search`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ keyword }),
  });
  return data?.bucket_courses || [];
}

export async function ListOfSearching() {
  const data = await apiFetch("/", {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  });
  return data
    ? data
    : "Ошибка"
}

