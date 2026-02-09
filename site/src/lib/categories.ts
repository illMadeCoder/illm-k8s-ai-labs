import fs from 'node:fs';
import path from 'node:path';

interface Domain {
  id: string;
  name: string;
  description: string;
  subdomains: string[];
}

interface CategoryData {
  domains: Domain[];
}

const categoriesPath = path.resolve(process.cwd(), 'data', '_categories.json');

let cached: CategoryData | null = null;

function loadCategoryData(): CategoryData {
  if (cached) return cached;
  const raw = fs.readFileSync(categoriesPath, 'utf-8');
  cached = JSON.parse(raw) as CategoryData;
  return cached;
}

export function loadCategories(): Domain[] {
  return loadCategoryData().domains;
}

export function getDomainMeta(domainId: string): Domain | undefined {
  return loadCategories().find((d) => d.id === domainId);
}

export function getAllDomainIds(): string[] {
  return loadCategories().map((d) => d.id);
}
