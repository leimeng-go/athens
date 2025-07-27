/// <reference types="vite/client" />

interface ImportMeta {
  readonly env: {
    readonly DEV: boolean;
    readonly PROD: boolean;
    readonly MODE: string;
    readonly SSR: boolean;
    readonly VITE_API_BASE_URL: string;
    [key: string]: string | boolean | undefined;
  };
}