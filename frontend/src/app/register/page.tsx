'use client';

import { FormEvent, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiPost } from '@/lib/api';

type Role = 'EMPLOYER' | 'APPLICANT';

export default function RegisterPage() {
  const router = useRouter();
  const [role, setRole] = useState<Role>('APPLICANT');

  const [email, setEmail] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [password, setPassword] = useState('');

  const [companyName, setCompanyName] = useState('');
  const [fullName, setFullName] = useState('');

  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const body = {
        email,
        password,
        displayName,
        role,
        ...(role === 'EMPLOYER' ? { companyName } : { fullName }),
      };

      await apiPost('/api/auth/register', body);
      router.push('/dashboard');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Ошибка регистрации');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-md px-4 py-10">
      <h1 className="text-2xl font-semibold text-black dark:text-white">Регистрация</h1>

      <form
        onSubmit={onSubmit}
        className="mt-6 space-y-4 rounded-2xl border border-black/10 bg-white/70 p-4 backdrop-blur dark:border-white/10 dark:bg-black/30"
      >
        <label className="block">
          <div className="text-sm font-medium text-black/70 dark:text-white/70">Роль</div>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value as Role)}
            className="mt-1 h-11 w-full rounded-xl border border-black/10 bg-white/60 px-4 text-sm outline-none dark:border-white/15 dark:bg-black/25 dark:text-white"
          >
            <option value="APPLICANT">Соискатель</option>
            <option value="EMPLOYER">Работодатель</option>
          </select>
        </label>

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
          <div className="text-sm font-medium text-black/70 dark:text-white/70">Имя</div>
          <input
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            required
            className="mt-1 h-11 w-full rounded-xl border border-black/10 bg-white/60 px-4 text-sm outline-none dark:border-white/15 dark:bg-black/25 dark:text-white"
          />
        </label>

        {role === 'EMPLOYER' ? (
          <label className="block">
            <div className="text-sm font-medium text-black/70 dark:text-white/70">Название компании</div>
            <input
              value={companyName}
              onChange={(e) => setCompanyName(e.target.value)}
              required
              className="mt-1 h-11 w-full rounded-xl border border-black/10 bg-white/60 px-4 text-sm outline-none dark:border-white/15 dark:bg-black/25 dark:text-white"
            />
          </label>
        ) : (
          <label className="block">
            <div className="text-sm font-medium text-black/70 dark:text-white/70">ФИО</div>
            <input
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              required
              className="mt-1 h-11 w-full rounded-xl border border-black/10 bg-white/60 px-4 text-sm outline-none dark:border-white/15 dark:bg-black/25 dark:text-white"
            />
          </label>
        )}

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
          {loading ? 'Создание…' : 'Зарегистрироваться'}
        </button>
      </form>
    </div>
  );
}

