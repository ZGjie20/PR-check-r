const GITHUB_PR_URL_PATTERN =
  /^https?:\/\/github\.com\/[^/]+\/[^/]+\/pull\/\d+\/?$/i;

export function isValidGitHubPrUrl(url: string): boolean {
  return GITHUB_PR_URL_PATTERN.test(url.trim());
}

export function parsePrNumber(url: string): number | null {
  const match = url.trim().match(/\/pull\/(\d+)/i);
  if (!match) return null;
  return Number.parseInt(match[1], 10);
}
