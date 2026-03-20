'use client';

import { useEffect } from 'react';

export default function ThemeBootstrap() {
  useEffect(() => {
    const saved = window.localStorage.getItem('theme');
    const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false;
    const theme = saved === 'light' || saved === 'dark' ? saved : prefersDark ? 'dark' : 'light';
    document.documentElement.classList.toggle('dark', theme === 'dark');
  }, []);

  return null;
}

