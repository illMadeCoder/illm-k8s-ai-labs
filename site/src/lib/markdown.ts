import type {
  ExperimentSummary,
  BodyBlock,
  CapabilitiesMatrix,
  AnalysisFeedback,
  QueryResult,
  DataPoint,
} from '../types';
import { formatDuration, formatDate, formatCost, formatValue } from './format';
import { getHypothesisClaim, getQuestions, getFocus, getMachineVerdict } from './experiments';

/**
 * Compute summary statistics for a range metric (multiple data points).
 * Groups by series when labels vary across data points.
 */
function summarizeRange(qr: QueryResult): string {
  const data = qr.data!;

  // Group data points by series
  const seriesMap = new Map<string, { label: string; points: DataPoint[] }>();
  const distinguishingKey = findDistinguishingLabelKey(data);

  for (const dp of data) {
    const seriesLabel = distinguishingKey && dp.labels?.[distinguishingKey]
      ? dp.labels[distinguishingKey]
      : '';
    const key = seriesLabel;
    if (!seriesMap.has(key)) {
      seriesMap.set(key, { label: seriesLabel, points: [] });
    }
    seriesMap.get(key)!.points.push(dp);
  }

  const seriesEntries = [...seriesMap.values()];

  if (seriesEntries.length === 1) {
    // Single series — inline stats
    const stats = computeStats(seriesEntries[0].points);
    return formatStats(stats, qr.unit);
  }

  // Multi-series — bullet list
  const bullets = seriesEntries.map((s) => {
    const stats = computeStats(s.points);
    return `- ${s.label}: ${formatStats(stats, qr.unit)}`;
  });
  return '\n' + bullets.join('\n');
}

/** Find a label key whose values differ across data points, for series grouping. */
function findDistinguishingLabelKey(data: DataPoint[]): string | null {
  if (data.length === 0) return null;
  const firstLabels = data[0].labels;
  if (!firstLabels) return null;

  for (const key of Object.keys(firstLabels)) {
    const firstVal = firstLabels[key];
    if (data.some((dp) => dp.labels?.[key] !== firstVal)) {
      return key;
    }
  }
  return null;
}

function computeStats(points: DataPoint[]): { min: number; max: number; avg: number; count: number } {
  let min = Infinity;
  let max = -Infinity;
  let sum = 0;
  for (const p of points) {
    if (p.value < min) min = p.value;
    if (p.value > max) max = p.value;
    sum += p.value;
  }
  return { min, max, avg: sum / points.length, count: points.length };
}

function formatStats(stats: { min: number; max: number; avg: number; count: number }, unit?: string): string {
  return `min ${formatValue(stats.min, unit)}, max ${formatValue(stats.max, unit)}, avg ${formatValue(stats.avg, unit)} (${stats.count} samples)`;
}

interface MarkdownProps {
  experiment: ExperimentSummary;
  displayName: string;
  tags: string[];
}

export function experimentToMarkdown({ experiment, displayName, tags }: MarkdownProps): string {
  const lines: string[] = [];
  const push = (...s: string[]) => lines.push(...s);
  const blank = () => lines.push('');

  const analysis = experiment.analysis;
  const verdict = analysis?.hypothesisVerdict;
  const verdictLabel = verdict
    ? verdict.charAt(0).toUpperCase() + verdict.slice(1)
    : null;
  const machineVerdict = getMachineVerdict(experiment);
  const hypothesisClaim = getHypothesisClaim(experiment);
  const questions = getQuestions(experiment);
  const focusAreas = getFocus(experiment);
  const hasRichAnalysis = !!analysis?.abstract;

  // --- Title & description ---
  push(`# ${displayName}`);
  blank();
  push(experiment.description);
  blank();

  if (tags.length > 0) {
    push(`**Tags:** ${tags.join(', ')}`);
    blank();
  }

  // --- Stats row ---
  const stats: string[] = [];
  const displayVerdict = verdictLabel ?? (machineVerdict
    ? machineVerdict.charAt(0).toUpperCase() + machineVerdict.slice(1)
    : null);
  if (displayVerdict) stats.push(`**Hypothesis:** ${displayVerdict}`);
  stats.push(`**Created:** ${formatDate(experiment.createdAt)}`);
  stats.push(`**Duration:** ${formatDuration(experiment.durationSeconds)}`);
  stats.push(`**Targets:** ${experiment.targets.length}`);
  if (experiment.costEstimate) stats.push(`**Cost:** ${formatCost(experiment.costEstimate.totalUSD)}`);
  stats.push(`**Workflow:** ${experiment.workflow.phase}`);
  push(stats.join(' | '));
  blank();

  // --- Architecture diagram ---
  if (analysis?.architectureDiagram) {
    push('## Architecture');
    blank();
    const lang = analysis.architectureDiagramFormat === 'mermaid' ? 'mermaid' : '';
    push(`\`\`\`${lang}`);
    push(analysis.architectureDiagram);
    push('```');
    blank();
  }

  // --- Targets table ---
  push('## Targets');
  blank();
  push('| Name | Cluster Type | Machine Type | Nodes |');
  push('|------|-------------|-------------|-------|');
  for (const t of experiment.targets) {
    push(`| ${t.name} | ${t.clusterType} | ${t.machineType ?? '—'} | ${t.nodeCount ?? '—'} |`);
  }
  blank();

  // --- Cost ---
  if (experiment.costEstimate) {
    push('## Cost');
    blank();
    const cost = experiment.costEstimate;
    if (cost.perTarget && Object.keys(cost.perTarget).length > 0) {
      push('| Target | Cost |');
      push('|--------|------|');
      for (const [name, usd] of Object.entries(cost.perTarget)) {
        push(`| ${name} | ${formatCost(usd)} |`);
      }
      push(`| **Total** | **${formatCost(cost.totalUSD)}** |`);
    } else {
      push(`**Total:** ${formatCost(cost.totalUSD)}`);
    }
    if (cost.note) {
      blank();
      push(`*${cost.note}*`);
    }
    blank();
  }

  // --- Overview: Hypothesis ---
  if (hypothesisClaim) {
    push('## Hypothesis');
    blank();
    push(`> ${hypothesisClaim}`);
    blank();
  }

  // --- Abstract ---
  if (hasRichAnalysis && analysis) {
    push('## Abstract');
    blank();
    push(`*${analysis.model} · ${formatDate(analysis.generatedAt)}*`);
    blank();
    push(analysis.abstract!);
    blank();
  }

  // --- Questions ---
  if (questions.length > 0) {
    push('## Questions');
    blank();
    questions.forEach((q, i) => push(`${i + 1}. ${q}`));
    blank();
  }

  // --- Focus areas ---
  if (focusAreas.length > 0) {
    push('## Focus Areas');
    blank();
    focusAreas.forEach((f) => push(`- ${f}`));
    blank();
  }

  // --- Legacy summary (when no rich analysis) ---
  if (!hasRichAnalysis && analysis?.summary) {
    push('## Summary');
    blank();
    push(analysis.summary);
    blank();
  }

  // --- Capabilities matrix ---
  if (hasRichAnalysis && analysis?.capabilitiesMatrix) {
    renderCapabilitiesMatrix(lines, analysis.capabilitiesMatrix);
  }

  // --- Analysis body ---
  if (hasRichAnalysis && analysis?.body?.blocks?.length) {
    push('## Analysis');
    blank();

    const textIndices = analysis.body.blocks.reduce<number[]>(
      (acc, b, i) => { if (b.type === 'text') acc.push(i); return acc; }, [],
    );
    const firstTextIdx = textIndices.length > 0 ? textIndices[0] : -1;
    const lastTextIdx = textIndices.length > 1 ? textIndices[textIndices.length - 1] : -1;

    analysis.body.blocks.forEach((block, blockIdx) => {
      if (block.type === 'text') {
        const header = blockIdx === firstTextIdx ? 'Analysis'
          : blockIdx === lastTextIdx ? 'Conclusion'
          : undefined;
        if (header) {
          push(`### ${header}`);
          blank();
        }
        push(block.content);
        blank();
      } else {
        renderBlock(lines, block, experiment, 3);
      }
    });
  }

  // --- Metrics ---
  const allMetricEntries = experiment.metrics
    ? Object.entries(experiment.metrics.queries).filter(([, qr]) => !qr.error)
    : [];
  const metricEntries = allMetricEntries.filter(([, qr]) => qr.data && qr.data.length > 0);

  if (metricEntries.length > 0) {
    push(`## Metrics (${metricEntries.length})`);
    blank();
    for (const [name, qr] of metricEntries) {
      push(`### ${name}${qr.unit ? ` (${qr.unit})` : ''}`);
      blank();
      if (qr.description) {
        push(qr.description);
        blank();
      }
      if (qr.data && qr.data.length === 1) {
        push(`**${formatValue(qr.data[0].value, qr.unit)}**`);
      } else {
        push(summarizeRange(qr));
      }
      blank();
      if (analysis?.metricInsights?.[name]) {
        push(`> ${analysis.metricInsights[name]}`);
        blank();
      }
    }
  }

  // --- Feedback ---
  if (hasRichAnalysis && analysis?.feedback) {
    renderFeedback(lines, analysis.feedback);
  }

  return lines.join('\n');
}

// ---- Block renderers ----

function renderBlock(
  lines: string[],
  block: BodyBlock,
  experiment: ExperimentSummary,
  headingLevel: number,
): void {
  const push = (...s: string[]) => lines.push(...s);
  const blank = () => lines.push('');
  const h = '#'.repeat(headingLevel);

  switch (block.type) {
    case 'text':
      push(block.content);
      blank();
      break;

    case 'topic':
      push(`${h} ${block.title}`);
      blank();
      for (const inner of block.blocks) {
        renderBlock(lines, inner, experiment, headingLevel + 1);
      }
      break;

    case 'metric': {
      const qr = experiment.metrics?.queries?.[block.key];
      if (!qr || qr.error || !qr.data?.length) {
        push(`*Metric \`${block.key}\` not available.*`);
        blank();
        break;
      }
      if (qr.data.length === 1) {
        push(`**${block.key}**${qr.unit ? ` (${qr.unit})` : ''}: ${formatValue(qr.data[0].value, qr.unit)}`);
      } else {
        push(`**${block.key}**${qr.unit ? ` (${qr.unit})` : ''}: ${summarizeRange(qr)}`);
      }
      blank();
      if (block.insight) {
        push(`> ${block.insight}`);
        blank();
      }
      break;
    }

    case 'comparison':
      push('| | Value | |');
      push('|---|---|---|');
      for (const item of block.items) {
        push(`| **${item.label}** | ${item.value} | ${item.description ?? ''} |`);
      }
      blank();
      break;

    case 'capabilityRow': {
      push(`**${block.capability}**`);
      blank();
      for (const [tech, assessment] of Object.entries(block.values)) {
        push(`- **${tech}:** ${assessment}`);
      }
      blank();
      break;
    }

    case 'table':
      push(`| ${block.headers.join(' | ')} |`);
      push(`|${block.headers.map(() => '---').join('|')}|`);
      for (const row of block.rows) {
        push(`| ${row.join(' | ')} |`);
      }
      if (block.caption) {
        blank();
        push(`*${block.caption}*`);
      }
      blank();
      break;

    case 'architecture': {
      const lang = block.format === 'mermaid' ? 'mermaid' : '';
      push(`\`\`\`${lang}`);
      push(block.diagram);
      push('```');
      if (block.caption) {
        blank();
        push(`*${block.caption}*`);
      }
      blank();
      break;
    }

    case 'callout':
      push(`> **[${block.variant.toUpperCase()}] ${block.title}**`);
      push(`> ${block.content}`);
      blank();
      break;

    case 'recommendation':
      push(`- **[${block.priority.toUpperCase()}] ${block.title}**${block.effort ? ` _(${block.effort} effort)_` : ''} — ${block.description}`);
      blank();
      break;

    case 'row':
      for (const child of block.blocks) {
        renderBlock(lines, child, experiment, headingLevel);
      }
      break;
  }
}

function renderCapabilitiesMatrix(lines: string[], matrix: CapabilitiesMatrix): void {
  const push = (...s: string[]) => lines.push(...s);
  const blank = () => lines.push('');
  const techs = matrix.technologies;

  push('## Capabilities');
  blank();

  if (matrix.summary) {
    push(`**Summary:** ${matrix.summary}`);
    blank();
  }

  for (const category of matrix.categories) {
    push(`### ${category.name}`);
    blank();
    push(`| Capability | ${techs.join(' | ')} |`);
    push(`|---|${techs.map(() => '---').join('|')}|`);
    for (const cap of category.capabilities) {
      const values = techs.map((t) => cap.values[t] ?? '—');
      push(`| ${cap.name} | ${values.join(' | ')} |`);
    }
    blank();
  }
}

function renderFeedback(lines: string[], feedback: AnalysisFeedback): void {
  const push = (...s: string[]) => lines.push(...s);
  const blank = () => lines.push('');

  const hasRecs = feedback.recommendations && feedback.recommendations.length > 0;
  const hasDesign = feedback.experimentDesign && feedback.experimentDesign.length > 0;

  if (hasRecs) {
    push('## Next Steps');
    blank();
    for (const rec of feedback.recommendations!) {
      push(`- ${rec}`);
    }
    blank();
  }

  if (hasDesign) {
    push('## Design Improvements');
    blank();
    for (const item of feedback.experimentDesign!) {
      push(`- ${item}`);
    }
    blank();
  }
}
