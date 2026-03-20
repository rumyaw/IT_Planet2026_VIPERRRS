'use client';

import { useEffect, useMemo, useRef } from 'react';
import type { OpportunityMarkerDTO } from '@/lib/dtos';

type YandexYmPlacemark = {
  events: {
    add: (eventName: string, cb: () => void) => void;
  };
};

type YandexYmMap = {
  geoObjects: {
    removeAll: () => void;
    add: (obj: unknown) => void;
  };
};

type YandexYmaps = {
  ready: (cb: () => void) => void;
  Map: new (container: HTMLElement, opts: Record<string, unknown>) => YandexYmMap;
  Placemark: new (
    coords: [number, number],
    props: Record<string, unknown>,
    opts: Record<string, unknown>
  ) => YandexYmPlacemark;
};

declare global {
  interface Window {
    ymaps?: YandexYmaps;
  }
}

let yandexScriptPromise: Promise<void> | null = null;

function loadYandexScript(apiKey: string) {
  if (typeof window === 'undefined') return Promise.resolve();
  if (window.ymaps) return Promise.resolve();
  if (yandexScriptPromise) return yandexScriptPromise;

  yandexScriptPromise = new Promise<void>((resolve, reject) => {
    const existing = document.querySelector<HTMLScriptElement>('script[data-ymaps="1"]');
    if (existing) {
      existing.addEventListener('load', () => resolve());
      existing.addEventListener('error', () => reject(new Error('Yandex script load failed')));
      return;
    }

    const script = document.createElement('script');
    script.dataset.ymaps = '1';
    script.src = `https://api-maps.yandex.ru/2.1/?apikey=${encodeURIComponent(apiKey)}&lang=ru_RU`;
    script.async = true;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error('Yandex script load failed'));
    document.head.appendChild(script);
  });

  return yandexScriptPromise;
}

type Props = {
  apiKey: string;
  markers: OpportunityMarkerDTO[];
  favoriteIds?: Set<string>;
  onSelect?: (m: OpportunityMarkerDTO) => void; // pinned on click
  onHover?: (m: OpportunityMarkerDTO | null) => void; // ephemeral on hover
  className?: string;
};

export default function YandexMap({
  apiKey,
  markers,
  favoriteIds,
  onSelect,
  onHover,
  className,
}: Props) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<YandexYmMap | null>(null);

  const center = useMemo(() => {
    if (!markers.length) return [55.7558, 37.6173] as [number, number]; // Moscow default
    const avgLat = markers.reduce((s, m) => s + m.lat, 0) / markers.length;
    const avgLng = markers.reduce((s, m) => s + m.lng, 0) / markers.length;
    return [avgLat, avgLng] as [number, number];
  }, [markers]);

  useEffect(() => {
    if (!apiKey) return;
    if (!containerRef.current) return;

    loadYandexScript(apiKey)
      .then(() => {
        const ymaps = window.ymaps;
        if (!ymaps) return;
        ymaps.ready(() => {
          if (!mapRef.current) {
            mapRef.current = new ymaps.Map(containerRef.current!, {
              center,
              zoom: markers.length ? 10 : 7,
              controls: ['zoomControl', 'fullscreenControl'],
            });
          }

          // Replace markers.
          const map = mapRef.current;
          map?.geoObjects.removeAll();

          markers.forEach((m) => {
            const isFav = favoriteIds?.has(m.id) ?? false;
            const preset = isFav ? 'islands#redIcon' : 'islands#blueIcon';

            const placemark = new ymaps.Placemark(
              [m.lat, m.lng],
              {},
              { preset }
            );
            placemark.events.add('click', () => onSelect?.(m));
            placemark.events.add('mouseenter', () => onHover?.(m));
            placemark.events.add('mouseleave', () => onHover?.(null));
            map?.geoObjects.add(placemark);
          });
        });
      })
      .catch(() => {
        // Keep UI usable even if Yandex script fails.
      });
  }, [apiKey, markers, favoriteIds, onSelect, onHover, center]);

  return (
    <div
      ref={containerRef}
      className={className ?? 'w-full h-[420px] rounded-2xl bg-black/5'}
    />
  );
}

