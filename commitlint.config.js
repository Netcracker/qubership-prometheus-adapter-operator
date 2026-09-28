/* eslint-disable no-undef -- CommonJS Node config; flat eslint has no Node globals */
const { execFileSync } = require('child_process');

const BOT_AUTHORS = new Set(['renovate[bot]', 'dependabot[bot]']);
const HUMAN_HEADER_MAX = 100;
const COMMIT_RECORD_MARKER = '---commitlint-author-lookup---';

function normalizeCommitMessage(msg) {
  return String(msg || '')
    .replace(/\r\n/g, '\n')
    .replace(/\n+$/, '');
}

function authorForCommitMessage(raw) {
  if (process.env.COMMITLINT_COMMIT_AUTHOR) {
    return process.env.COMMITLINT_COMMIT_AUTHOR.trim();
  }

  // commitlint only hands the rule the message text (stdin, --last, or each
  // --from/--to entry). Identify that exact commit by a unique full-message
  // match in the same range super-linter uses. Newest-first subject lookup is
  // unsafe: a later bot commit with the same subject would exempt a human.
  const needle = normalizeCommitMessage(raw);
  if (!needle) {
    return '';
  }

  const from = process.env.GITHUB_BEFORE_SHA;
  const to = process.env.GITHUB_SHA || 'HEAD';
  const gitArgs = from
    ? ['log', `${from}..${to}`, `--format=${COMMIT_RECORD_MARKER}%n%H%n%an%n%B`]
    : ['log', '-n', '50', to, `--format=${COMMIT_RECORD_MARKER}%n%H%n%an%n%B`];

  try {
    const log = execFileSync('git', gitArgs, {
      encoding: 'utf8',
      maxBuffer: 20 * 1024 * 1024,
      stdio: ['ignore', 'pipe', 'ignore'],
    });

    const authors = [];
    for (const chunk of log.split(COMMIT_RECORD_MARKER).slice(1)) {
      const body = chunk.replace(/^\n/, '');
      const nl1 = body.indexOf('\n');
      if (nl1 === -1) {
        continue;
      }
      const nl2 = body.indexOf('\n', nl1 + 1);
      if (nl2 === -1) {
        continue;
      }
      const author = body.slice(nl1 + 1, nl2);
      const message = normalizeCommitMessage(body.slice(nl2 + 1));
      if (message === needle) {
        authors.push(author);
      }
    }

    // Ambiguous or missing match: fail closed (human header-max-length applies).
    if (authors.length !== 1) {
      return '';
    }
    return authors[0];
  } catch {
    return '';
  }
}

module.exports = {
  extends: ['@commitlint/config-conventional'],
  helpUrl: 'https://www.conventionalcommits.org/',
  // We need this until https://github.com/dependabot/dependabot-core/issues/2445
  // is resolved.
  ignores: [(msg) => /Signed-off-by: dependabot\[bot\]/m.test(msg)],
  plugins: [
    {
      rules: {
        // Keep the 100-character header limit for humans. Dependency bots own their
        // subjects and rewrite them on rebase, so only renovate[bot] and
        // dependabot[bot] skip header-max-length.
        'header-max-length-bot-aware': ({ header, raw }) => {
          const author = authorForCommitMessage(raw || header || '');
          if (BOT_AUTHORS.has(author)) {
            return [true];
          }

          const length = (header || '').length;
          if (length <= HUMAN_HEADER_MAX) {
            return [true];
          }

          return [
            false,
            `header must not be longer than ${HUMAN_HEADER_MAX} characters, current length is ${length}`,
          ];
        },
      },
    },
  ],
  rules: {
    'header-max-length': [0],
    'header-max-length-bot-aware': [2, 'always'],
  },
};
