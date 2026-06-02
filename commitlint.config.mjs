export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [2, 'always', [
      'feat', 'fix', 'chore', 'docs', 'style',
      'refactor', 'test', 'ci', 'perf', 'build', 'revert',
    ]],
    'body-max-line-length': [0],
    'footer-max-line-length': [0],
    'subject-case': [0],
  },
};
