'use client';

import { FormEvent, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiPost } from '@/lib/api';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await apiPost('/api/auth/login', { email, password });
      router.push('/dashboard');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Ошибка входа');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-md px-4 py-10">
      <h1 className="text-2xl font-semibold text-black dark:text-white">Вход</h1>

      <form onSubmit={onSubmit} className="mt-6 space-y-4 rounded-2xl border border-black/10 bg-white/70 p-4 backdrop-blur dark:border-white/10 dark:bg-black/30">
        <label className="block">
          <div className="text-sm font-medium text-black/70 dark:text-white/70">Email</div>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            type="email"
            required
            className="mt-1 h-11 w-full rounded-xl border border-black/10 bg-white/60 px-4 text-sm outline-none dark:border-white/15 dark:bg-black/25 dark:text-white"
          />
        </label>

        <label className="block">
          <div className="text-sm font-medium text-black/70 dark:text-white/70">Пароль</div>
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            type="password"
            required
            className="mt-1 h-11 w-full rounded-xl border border-black/10 bg-white/60 px-4 text-sm outline-none dark:border-white/15 dark:bg-black/25 dark:text-white"
          />
        </label>

        {error ? <div className="text-sm font-medium text-rose-600">{error}</div> : null}

        <button
          disabled={loading}
          type="submit"
          className="h-11 w-full rounded-xl bg-indigo-600 px-4 text-sm font-semibold text-white transition hover:bg-indigo-500 disabled:opacity-60"
        >
          {loading ? 'Вход…' : 'Войти'}
        </button>
      </form>
    </div>
  );
}

