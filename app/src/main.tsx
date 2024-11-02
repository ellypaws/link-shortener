import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import '@/styles/globals.css';
import LinkShortener from "@/components/link-shortener.tsx";
import {Toaster} from "@/components/ui/toaster.tsx";

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <LinkShortener/>
    </main>
    <Toaster/>
  </React.StrictMode>
);
