"use client";

import type { User } from "./types";

const TOKEN_KEY = "tofex_admin_token";
const USER_KEY = "tofex_admin_user";

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(TOKEN_KEY);
}

export function getStoredUser(): User | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as User;
  } catch {
    return null;
  }
}

export function setSession(token: string, user: User) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

export function hasPermission(user: User | null, perm: string): boolean {
  if (!user?.roles) return false;
  for (const role of user.roles) {
    for (const p of role.permissions ?? []) {
      if (p.permission === "*" || p.permission === perm) return true;
    }
  }
  return false;
}

export function isSuperAdmin(user: User | null): boolean {
  return user?.roles?.some((r) => r.name === "super_admin") ?? false;
}
