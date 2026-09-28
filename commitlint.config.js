const { execFileSync } = require('child_process');

const BOT_AUTHORS = new Set(['renovate[bot]', 'dependabot[bot]']);
const HUMAN_HEADER_MAX = 100;

function authorForCommitMessage(raw) {
  if (process.env.COMMITLINT_COMMIT_AUTHOR) {
    return process.env.COMMITLINT_COMMIT_AUTHOR.trim();
  }

  const header = String(raw || '')
    .split('\n')[0]
    .trim();
  if (!header) {
    return '';
  }

  // Prefer the range super-linter passes to commitlint (--from/--to). Fall back to
  // recent HEAD ancestry so a bot commit on its own branch still resolves. Never
  // search --all: that would exempt a human who pastes a bot subject on stdin.
  const from = process.env.GITHUB_BEFORE_SHA;
  const to = process.env.GITHUB_SHA || 'HEAD';
  const gitArgs = from
    ? ['log', `${from}..${to}`, '--format=%an%x09%s']
    : ['log', '-n', '50', to, '--format=%an%x09%s'];

  try {
    const log = execFileSync('git', gitArgs, {
      encoding: 'utf8',
      maxBuffer: 20 * 1024 * 1024,
      stdio: ['ignore', 'pipe', 'ignore'],
    });
    for (const line of log.split('\n')) {
      if (!line) {
        continue;
      }
      const tab = line.indexOf('\t');
      if (tab === -1) {
        continue;
      }
      const author = line.slice(0, tab);
      const subject = line.slice(tab + 1);
      if (subject === header) {
        return author;
      }
    }
  } catch {
    return '';
  }

  return '';
}

module.exports = {
  extends: ['@commitlint/config-conventional'],
  helpUrl: 'https://www.conventionalcommits.org/',
  // We need this until https://github.com/dependabot/dependabot-core/issues/2445
  // is resolved.
  ignores: [(msg) => /Signed-off-by: dependabot\[bot]/m.test(msg)],
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
