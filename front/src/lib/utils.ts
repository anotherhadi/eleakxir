import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// take "2025-09-13T21:14:46.13030464+02:00"
// return "13/09/2025 21:14"
export function formatDate(date: string) {
  const d = new Date(date);
  const day = String(d.getDate()).padStart(2, "0");
  const month = String(d.getMonth() + 1).padStart(2, "0");
  const year = d.getFullYear();
  const hours = String(d.getHours()).padStart(2, "0");
  const minutes = String(d.getMinutes()).padStart(2, "0");
  return `${day}/${month}/${year} ${hours}:${minutes}`;
}

export function convertNanoSeconds(nanoseconds: number): string {
  const ONE_MS_IN_NS = 1e6;
  const ONE_S_IN_NS = 1e9;
  const ONE_MIN_IN_NS = 6e10;
  const ONE_HOUR_IN_NS = 3.6e12;

  if (nanoseconds < ONE_MS_IN_NS) {
    return `${nanoseconds} ns`;
  } else if (nanoseconds < ONE_S_IN_NS) {
    const ms = Math.round(nanoseconds / ONE_MS_IN_NS);
    return `${ms} ms`;
  } else if (nanoseconds < ONE_MIN_IN_NS) {
    const s = Math.round(nanoseconds / ONE_S_IN_NS);
    return `${s} s`;
  } else if (nanoseconds < ONE_HOUR_IN_NS) {
    const totalSeconds = Math.round(nanoseconds / ONE_S_IN_NS);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes}m${seconds}s`;
  } else {
    const totalMinutes = Math.round(nanoseconds / ONE_MIN_IN_NS);
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    return `${hours}h${minutes}m`;
  }
}

type FlatObject = { [key: string]: any };

export function FlattenObject(obj: object): FlatObject {
  const flattened: FlatObject = {};

  function recurse(currentObj: any, prefix: string = ""): void {
    for (const key in currentObj) {
      if (Object.prototype.hasOwnProperty.call(currentObj, key)) {
        const value = currentObj[key];
        const newKey = prefix ? `${prefix}.${key}` : key;

        if (
          typeof value === "object" &&
          value !== null &&
          !Array.isArray(value)
        ) {
          // Si la valeur est un objet, on continue la récursion
          recurse(value, newKey);
        } else if (Array.isArray(value)) {
          // Si la valeur est un tableau, on itère sur ses éléments
          value.forEach((item, index) => {
            // On continue la récursion pour les objets dans le tableau
            if (typeof item === "object" && item !== null) {
              recurse(item, `${newKey}.${index}`);
            } else {
              // On ajoute les valeurs primitives
              flattened[`${newKey}.${index}`] = item;
            }
          });
        } else {
          // Si la valeur est une primitive, on l'ajoute à l'objet aplati
          flattened[newKey] = value;
        }
      }
    }
  }

  recurse(obj);
  return flattened;
}
