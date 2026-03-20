'use client';

export default function ThemeToggle() {
  const onToggle = () => {
    const isDark = document.documentElement.classList.contains('dark');
    const next = isDark ? 'light' : 'dark';
    document.documentElement.classList.toggle('dark', next === 'dark');
    window.localStorage.setItem('theme', next);
  };

  return (
    <button
      type="button"
      onClick={onToggle}
      className="rounded-full border border-black/10 bg-white/70 px-3 py-1.5 text-sm font-medium text-black/80 backdrop-blur transition hover:bg-white dark:border-white/15 dark:bg-black/40 dark:text-white/85 dark:hover:bg-black/55"
      aria-label="Toggle theme"
    >
      Тема
    </button>
  );
}

