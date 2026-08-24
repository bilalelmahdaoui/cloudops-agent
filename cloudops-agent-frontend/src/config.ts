const apiBaseUrl = import.meta.env.VITE_API_URL;

if (!apiBaseUrl) {
  throw new Error("VITE_API_URL is required");
}

export const config = {
  apiBaseUrl: apiBaseUrl.replace(/\/$/, ""),
} as const;
