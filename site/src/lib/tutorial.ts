import fs from 'node:fs';
import path from 'node:path';
import type { TutorialFlow } from '../types';

const dataDir = path.resolve(process.cwd(), 'data');

export function loadTutorial(slug: string): TutorialFlow | null {
  const filePath = path.join(dataDir, `${slug}.tutorial.json`);
  if (!fs.existsSync(filePath)) return null;
  const raw = fs.readFileSync(filePath, 'utf-8');
  return JSON.parse(raw) as TutorialFlow;
}
